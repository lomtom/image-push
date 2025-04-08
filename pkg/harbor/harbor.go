package harbor

import (
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/mholt/archiver/v3"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	manifestFilename = "manifest.json"
)

type Config struct {
	address   string
	project   string
	chunkSize int
	tmpFolder string
	client    *http.Client
}

func NewConfig(address, username, password, project, archiveName string, skipTLS bool, ChunkSize int, file io.Reader) (*Config, error) {
	tmpFolder, err := generateTmpFolder()
	if err != nil {
		return nil, err
	}

	archivePath := filepath.Join(tmpFolder, archiveName)
	localFile, err := os.Create(archivePath)
	if err != nil {
		return nil, err
	}
	defer localFile.Close()

	_, err = io.Copy(localFile, file)
	if err != nil {
		return nil, err
	}

	// unarchive
	err = archiver.Unarchive(archivePath, tmpFolder)
	if err != nil {
		return nil, fmt.Errorf("unarchive failed, %+v", err)
	}

	address = strings.TrimSuffix(address, "/")
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipTLS},
		Proxy: func(req *http.Request) (*url.URL, error) {
			req.SetBasicAuth(username, password)
			return nil, nil
		},
	}
	return &Config{
		address:   address,
		project:   project,
		chunkSize: ChunkSize,
		tmpFolder: tmpFolder,
		client: &http.Client{
			Transport: tr,
		},
	}, nil
}

type Job interface {
	Run(Config) error
}

func (c *Config) Push() error {
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(c.tmpFolder)

	manifests, err := c.parseManifestFile()
	if err != nil {
		return err
	}

	jobs := make([]Job, 0)
	for _, manifest := range manifests {
		for _, repoTag := range manifest.RepoTags {
			repo, tag, err := parseImageAndTag(repoTag)
			if err != nil {
				log.Printf("paseRepoTag faild: %s", err)
				continue
			}

			// layers
			layerInfos := make([]*FileInfo, len(manifest.Layers))
			for k, layerPath := range manifest.Layers {
				fileInfo, err := c.getFileInfo(layerPath)
				if err != nil {
					return err
				}
				layerInfos[k] = fileInfo
				layerJob := &FileJob{
					repo:    repo,
					project: c.project,
					info:    fileInfo,
				}
				jobs = append(jobs, layerJob)
			}

			// contianer config
			fileInfo, err := c.getFileInfo(manifest.Config)
			if err != nil {
				return err
			}
			containerConfigJob := &FileJob{
				repo:    repo,
				project: c.project,
				info:    fileInfo,
			}
			jobs = append(jobs, containerConfigJob)

			// manifest
			manifestJob := &ManifestJob{
				layerInfos:         layerInfos,
				manifestConfigInfo: fileInfo,
				repo:               repo,
				project:            c.project,
				tag:                tag,
			}
			jobs = append(jobs, manifestJob)
		}
	}

	for _, v := range jobs {
		if err = v.Run(*c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) getFileInfo(path string) (*FileInfo, error) {
	absolutePath := filepath.Join(c.tmpFolder, path)
	f, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, err
	}
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return nil, err
	}
	hash := hex.EncodeToString(h.Sum(nil))
	return &FileInfo{
		absolutePath: absolutePath,
		fileSize:     fileInfo.Size(),
		digest:       fmt.Sprint("sha256:", hash),
	}, nil
}

func (c *Config) parseManifestFile() (Manifests, error) {
	manifestFilePath := filepath.Join(c.tmpFolder, manifestFilename)
	data, err := os.ReadFile(manifestFilePath)
	if err != nil {
		return nil, err
	}
	manifests := make(Manifests, 0)
	if err = json.Unmarshal(data, &manifests); err != nil {
		return nil, err
	}
	return manifests, nil
}
