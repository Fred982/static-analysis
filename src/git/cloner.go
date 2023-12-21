package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

func CloneRepository(repoURL, localPath string) error {
	parts := strings.Split(repoURL, "/")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")
	repoPath := filepath.Join(localPath, repoName)
	_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
		URL:      repoURL,
		Progress: os.Stdout,
	})
	if err != nil {
		return fmt.Errorf("error cloning repository: %v", err)
	}
	return nil
}
