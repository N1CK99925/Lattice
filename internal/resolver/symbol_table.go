package resolver

import (
	"Lattice/internal/models"
)

type SymbolTable struct {
	ByID   map[string]*models.Symbol
	ByName map[string][]*models.Symbol
}

func BuildSymbolTable(index *models.RepositoryIndex) SymbolTable {

	table := SymbolTable{
		ByID:   make(map[string]*models.Symbol),
		ByName: make(map[string][]*models.Symbol),
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
