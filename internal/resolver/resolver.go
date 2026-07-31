package resolver

import (
	"Lattice/internal/parser"
	"fmt"
)

func Resolve(index *parser.RepositoryIndex) {
	symbolTable := BuildSymbolTable(index)
	for _, file := range index.Files {
		for _, call := range file.Calls {
			symbol, ok := symbolTable[call.Target]
			if !ok {
				continue
			}
			fmt.Printf("%s -> %s\n", call.Target, symbol.Name)
		}
	}
}

//what this does is it takes the RepositoryIndex and then maps it to their respective symbols
