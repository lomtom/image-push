package harbor

// ParseImageAndTag parse image repo
import (
	"errors"
	"github.com/docker/distribution/uuid"
	"os"
	"path/filepath"
	"strings"
)

func parseImageAndTag(repo string) (string, string, error) {
	index := strings.Index(repo, "/")
	var repoAndTag string
	if index == -1 {
		repoAndTag = repo
	} else {
		repoAndTag = repo[index+1:]
	}
	arr := strings.Split(repoAndTag, ":")
	if len(arr) != 2 {
		return "", "", errors.New("invalid repo tag")
	}
	if strings.Contains(arr[0], "/") {
		repos := strings.Split(arr[0], "/")
		if len(repos) != 2 {
			return "", "", errors.New("invalid repo tag")
		}
		return repos[1], arr[1], nil
	}
	return arr[0], arr[1], nil
}

func generateTmpFolder() (string, error) {
	path := filepath.Join("/tmp/docker_push", uuid.Generate().String())
	err := os.MkdirAll(path, 0755)
	if err != nil {
		return "", err
	}
	return path, nil
}
