package harbor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/docker/distribution"
	"github.com/docker/distribution/manifest/schema2"
	"github.com/opencontainers/go-digest"
	"net/http"
)

type Manifests = []Manifest

type Manifest struct {
	Config   string   `json:"Config"`
	RepoTags []string `json:"RepoTags"`
	Layers   []string `json:"Layers"`
}

type ManifestJob struct {
	layerInfos         []*FileInfo
	manifestConfigInfo *FileInfo
	project            string
	repo               string
	tag                string
}

func (j *ManifestJob) getRepoName() string {
	if j.project != "" {
		return j.project + "/" + j.repo
	}
	return j.repo
}

func (j *ManifestJob) Run(c Config) error {
	manifest := &schema2.Manifest{}
	manifest.SchemaVersion = schema2.SchemaVersion.SchemaVersion
	manifest.MediaType = schema2.MediaTypeManifest
	manifest.Config.MediaType = schema2.MediaTypeImageConfig
	manifest.Config.Size = j.manifestConfigInfo.fileSize
	manifest.Config.Digest = digest.Digest(j.manifestConfigInfo.digest)
	for _, v := range j.layerInfos {
		item := distribution.Descriptor{
			MediaType: schema2.MediaTypeUncompressedLayer,
			Size:      v.fileSize,
			Digest:    digest.Digest(v.digest),
		}
		manifest.Layers = append(manifest.Layers, item)
	}
	data, err := json.Marshal(manifest)
	if err != nil {
		return err
	}
	url := fmt.Sprintf("%s/v2/%s/manifests/%s", c.address, j.getRepoName(), j.tag)
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", schema2.MediaTypeManifest)
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("put manifest failed, code is %d", resp.StatusCode)
	}
	return nil
}
