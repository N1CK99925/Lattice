package parser

import (
	"Lattice/internal/repository"
	"path/filepath"
)

func BuildRepositoryIndex(path string) (*RepositoryIndex, error) {
	absRoot, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	files, err := repository.Walk(absRoot)
	if err != nil {
		return nil, err

	}
	index := &RepositoryIndex{}
	for _, file := range files {
		abs := filepath.Join(absRoot, file)
		parsedFile, err := ParseFile(abs)
		if err != nil {
			return nil, err
		}

		index.Files = append(index.Files, parsedFile)
	}
	return index, nil

}
