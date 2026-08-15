package resolver

import (
	"Lattice/internal/models"
)

type Resolver struct {
	table SymbolTable
}

func ResolveCall(table SymbolTable, call models.Call) (*models.Symbol, bool) {
	caller, ok := table.ByID[call.ParentSymbolID]
	if !ok {
		return nil, false
	}
	candidates, ok := table.ByName[call.Target]
	if !ok {
		return nil, false
	}
	for _, symbol := range candidates {
		if symbol.File == caller.File {
			return symbol, true
		}
	}
	if len(candidates) == 1 {
		return candidates[0], true
	}
	return nil, false
}

func Resolve(index *models.RepositoryIndex) func(models.Call) (string, bool) {
	table := BuildSymbolTable(index)

	return func(call models.Call) (string, bool) {
		symbol, ok := ResolveCall(table, call)

		if !ok {
			return "", false
		}

		return symbol.ID, true
	}
}

//what this does is it takes the RepositoryIndex and then maps it to their respective symbols
