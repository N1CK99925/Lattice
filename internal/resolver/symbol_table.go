package resolver

import (
	"Lattice/internal/parser"
)

type SymbolTable struct {
	ByID   map[string]*parser.Symbol
	ByName map[string][]*parser.Symbol
}

func BuildSymbolTable(index *parser.RepositoryIndex) SymbolTable {

	table := SymbolTable{
		ByID:   make(map[string]*parser.Symbol),
		ByName: make(map[string][]*parser.Symbol),
	}
	for _, file := range index.Files {
		for i := range file.Symbols {
			symbol := &file.Symbols[i]
			table.ByID[symbol.ID] = symbol
			table.ByName[symbol.Name] = append(table.ByName[symbol.Name], symbol)
		}
	}
	// symbol, ok := table.ByID["internal/parser/reader.go::ParseFile"]
	// if !ok {
	// 	fmt.Println("ParseFile not found")
	// } else {
	// 	// fmt.Printf("%+v\n", *symbol)
	// }
	return table

}
