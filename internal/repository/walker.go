package repository

import (
	"fmt"
	"io/fs"
	"log"
	"os"
)

func Run_Walker() error {
	fmt.Println("Do you want to run the current repository or another file path?")
	fmt.Println("1. Run the current repository")
	fmt.Println("2. Run another file path")
	var choice int
	fmt.Scanf("%d", &choice)

	switch {
	case choice == 1:
		Walk_dir(".")
		fmt.Println("Walking current repository")
		return nil

	case choice == 2:
		fmt.Println("Enter the file path")
		var path string
		fmt.Scanf("%s", &path)
		Walk_dir(path)
		fmt.Println("Walking another file path")
		return nil

	default:
		fmt.Println("Invalid choice")
		return nil
	}
}

func Walk_dir(root_path string) error {
	root := root_path
	filesystem := os.DirFS(root)

	fs.WalkDir(filesystem, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Printf("err %v", err)
		}
		fmt.Println(path)
		return nil
	})
	return nil
}
