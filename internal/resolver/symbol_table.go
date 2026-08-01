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
			symbolTable[symbol.ID] = symbol
		}
	}
	symbol, ok := symbolTable["internal/parser/reader.go::ParseFile"]
	if !ok {
		fmt.Println("ParseFile not found")
	} else {
		fmt.Printf("%+v\n", *symbol)
	}
	return symbolTable

}
