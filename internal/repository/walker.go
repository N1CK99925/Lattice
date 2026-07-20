package repository

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func Run_Walker() ([]string, error) {
	fmt.Println("Do you want to run the current repository or another file path?")
	fmt.Println("1. Run the current repository")
	fmt.Println("2. Run another file path")
	var choice int
	fmt.Scanf("%d", &choice)

	switch {
	case choice == 1:
		files, err := Walk_dir(".")
		fmt.Println("Walking current repository")
		if err != nil {
			return nil, err
		}
		return files, nil

	case choice == 2:
		fmt.Println("Enter the file path")
		var path string
		fmt.Scanf("%s", &path)
		fmt.Println("Walking another file path")
		files, err := Walk_dir(path)
		if err != nil {
			return nil, err
		}
		return files, nil

	default:
		fmt.Println("Invalid choice")
		return nil, nil
	}

}

func Walk_dir(root_path string) ([]string, error) {
	root := root_path
	filesystem := os.DirFS(root)
	var Files []string
	fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("err %v", err)
		}
		if d.IsDir() && strings.HasPrefix(d.Name(), ".git") {
			return fs.SkipDir
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			Files = append(Files, path)
		}
		fmt.Println(Files)
		return nil
	})
	return Files, nil
}
