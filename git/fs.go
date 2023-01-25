package git

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// FindRepositoryRoot looks up the tree from the specified absolute path
// to determine the path of the parent git repository
func FindRepositoryRoot(dir string) (string, error) {
	for {
		if dirHasChild(dir, ".git") {
			return dir, nil
		}

		dir = filepath.Dir(dir)
		if dir == "." || dir == "" || dir == filepath.Dir(dir) {
			break
		}
	}

	return "", fmt.Errorf("could not find the repository root")
}

// dirHasChild determines if the specified absolute path to a directory contains
// a child with the desired name.
func dirHasChild(dir string, childName string) bool {
	children, err := os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return false
	}
	for _, child := range children {
		if child.Name() == childName {
			return true
		}
	}
	return false
}
