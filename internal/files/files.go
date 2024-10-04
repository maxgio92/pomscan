package files

import (
	"os"
	"path/filepath"
)

func FindFiles(dir, name string) ([]string, error) {
	found := []string{}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		path := filepath.Join(dir, file.Name())
		if file.IsDir() {
			f, _ := FindFiles(path, name)
			found = append(found, f...)
		} else if file.Name() == name {
			found = append(found, path)
		}
	}

	return found, nil
}
