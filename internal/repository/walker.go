package repository

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func Walk(root string) ([]string, error) {

	files, err := Walk_dir(root)
	if err != nil {
		return nil, err
	}
	return files, nil

}

func Walk_dir(root_path string) ([]string, error) {
	root := root_path
	filesystem := os.DirFS(root)
	var Files []string
	err := fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		//  TODO: create a func to add more files ignore types
		if d.IsDir() && strings.HasPrefix(d.Name(), ".git") {
			return fs.SkipDir
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			Files = append(Files, path)
		}
		return nil
	})
	return Files, err
}
