package parser

import (
	"Lattice/internal/models"
	"Lattice/internal/repository"
	"path/filepath"
)

func BuildRepositoryIndex(path string) (*models.RepositoryIndex, error) {
	absRoot, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	files, err := repository.Walk(absRoot)
	if err != nil {
		return nil, err

	}
	index := &models.RepositoryIndex{}
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
