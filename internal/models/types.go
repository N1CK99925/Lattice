package models

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

type ParseContext struct {
	Path    string
	Package string
}

type Symbol struct {
	ID        string
	Name      string
	Package   string
	File      string
	Receiver  string
	Kind      SymbolKind
	StartLine uint32
	EndLine   uint32
	Scope     string
}

type Import struct {
	Path string
	File string
}
type Call struct {
	ParentSymbolID string
	Target         string
	File           string
	Receiver       string
	StartLine      uint32
	EndLine        uint32
}
type TypeRef struct {
	Text      string
	File      string
	StartLine uint32
	EndLine   uint32
}
type ParsedFile struct {
	Path     string
	Package  string
	Symbols  []Symbol
	Imports  []Import
	Calls    []Call
	TypeRefs []TypeRef
}

type RepositoryIndex struct {
	Files []ParsedFile
}
