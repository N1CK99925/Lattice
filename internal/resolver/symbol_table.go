package resolver

import (
	"Lattice/internal/parser"
	"fmt"
)

func BuildSymbolTable(index *parser.RepositoryIndex) map[string]*parser.Symbol {
	symbolTable := make(map[string]*parser.Symbol)

	for _, file := range index.Files {
		for i := range file.Symbols {
			symbol := &file.Symbols[i]
			symbolTable[symbol.Name] = symbol
		}
	}
	fmt.Printf("%+v\n", *symbolTable["ParseFile"])
	return symbolTable

}
