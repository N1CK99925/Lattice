package resolver

import (
	"Lattice/internal/models"
	"fmt"
)

func Resolve(index *models.RepositoryIndex) {
	symbolTable := BuildSymbolTable(index)
	// for k := range symbolTable.ByID {
	// fmt.Println(k)
	// }
	for _, file := range index.Files {
		for _, call := range file.Calls {
			symbols, ok := symbolTable.ByName[call.Target]
			if !ok {
				continue
			}
			fmt.Printf("%s -> %s\n", call.Target, symbols[0].Name)
		}
	}
}

//what this does is it takes the RepositoryIndex and then maps it to their respective symbols
