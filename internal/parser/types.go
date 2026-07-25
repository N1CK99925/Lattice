package parser

//this contains all the types that are used in Query of tree sitter

type SymbolKind string

const (
	FunctionSymbol SymbolKind = "function"
	MethodSymbol   SymbolKind = "method"
	TypeSymbol     SymbolKind = "type"
	ConstantSymbol SymbolKind = "constant"
	VariableSymbol SymbolKind = "variable"
)

var GoSymbolKinds = map[string]SymbolKind{
	"function_declaration": FunctionSymbol,
	"method_declaration":   MethodSymbol,
	"type_spec":            TypeSymbol,
	"const_spec":           ConstantSymbol,
	"var_spec":             VariableSymbol,
}

type Symbols struct {
	Name      string
	Kind      SymbolKind
	StartLine uint32
	EndLine   uint32
}

type ParsedFile struct {
	Path    string
	Symbols []Symbols
}

type RepositoryIndex struct {
	Files []ParsedFile
}
