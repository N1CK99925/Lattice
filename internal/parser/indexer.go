package parser

import (
	"Lattice/internal/repository"
)

func BuildRepositoryIndex(path string) (*RepositoryIndex, error) {
	files, err := repository.Walk(path)
	if err != nil {
		return nil, err

	}
	index := &RepositoryIndex{}
	for _, file := range files {
		parsedFile, err := ParseFile(file)
		if err != nil {
			return nil, err
		}

		index.Files = append(index.Files, parsedFile)
	}
	return index, nil

}
