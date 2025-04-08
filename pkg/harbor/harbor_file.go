package harbor

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

type FileInfo struct {
	absolutePath string
	fileSize     int64
	digest       string
}

type FileJob struct {
	repo    string
	project string
	info    *FileInfo
}

func (j *FileJob) getRepoName() string {
	if j.project != "" {
		return j.project + "/" + j.repo
	}
	return j.repo
}

func (j *FileJob) Run(c Config) error {
	exist, err := j.checkLayerExist(c)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	url := fmt.Sprintf("%s/v2/%s/blobs/uploads/", c.address, j.getRepoName())
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}
	resp, err := c.client.Do(req)
	location := resp.Header.Get("Location")
	if resp.StatusCode != http.StatusAccepted || location == "" {
		return fmt.Errorf("post %s failed, statusCode:%d", url, resp.StatusCode)
	}
	return j.upload(c, location)
}

func (j *FileJob) checkLayerExist(c Config) (bool, error) {
	url := fmt.Sprintf("%s/v2/%s/blobs/%s", c.address, j.getRepoName(), j.info.digest)
	req, err := http.NewRequest(http.MethodHead, url, nil)
	if err != nil {
		return false, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	if resp.StatusCode == http.StatusOK {
		return true, nil
	}
	return false, nil
}

func (j *FileJob) upload(c Config, url string) error {
	f, err := os.Open(j.info.absolutePath)
	if err != nil {
		return err
	}

	defer f.Close()
	contentSize := j.info.fileSize
	// Monolithic Upload
	if c.chunkSize == 0 {
		url = fmt.Sprintf("%s&digest=%s", url, j.info.digest)
		contentLength := strconv.Itoa(int(j.info.fileSize))
		resp, err := j.doUpload(c, http.MethodPut, url, contentLength, "", f)
		if err != nil {
			return err
		}
		if resp.StatusCode != http.StatusCreated {
			return fmt.Errorf("monolithic upload failed, response code: %d", resp.StatusCode)
		}
	} else {
		// Chunked Upload
		index, offset := 0, 0
		buf := make([]byte, c.chunkSize)
		for {
			n, err := f.Read(buf)
			if err == io.EOF {
				break
			}
			offset = index + n
			index = offset
			chunk := buf[0:n]

			contentLength := strconv.Itoa(n)
			contentRange := fmt.Sprintf("%d-%d", index, offset)

			if int64(offset) == contentSize {
				url = fmt.Sprintf("%s&digest=%s", url, j.info.digest)
				resp, err := j.doUpload(c, http.MethodPut, url, contentLength, contentRange, bytes.NewBuffer(chunk))
				if err != nil {
					return err
				}
				if resp.StatusCode != http.StatusCreated {
					return fmt.Errorf("chunked upload faild,response code: %d", resp.StatusCode)
				}
				break
			} else {
				resp, err := j.doUpload(c, http.MethodPatch, url, contentLength, contentRange, bytes.NewBuffer(chunk))
				if err != nil {
					return err
				}
				location := resp.Header.Get("Location")
				if resp.StatusCode == http.StatusAccepted && location != "" {
					url = location
				} else {
					return fmt.Errorf("chunked upload faild,response code: %d", resp.StatusCode)
				}
			}
		}
	}

	return nil
}

func (j *FileJob) doUpload(c Config, method, url, contentLength, contentRange string, reader io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("Content-Length", contentLength)
	if contentRange != "" {
		req.Header.Set("Content-Range", contentRange)
	}
	log.Printf("%s %s", method, url)
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
