package parser

import (
	"fmt"
	"os"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func ReadFiles(file string) ([]byte, error) {

	source, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return source, nil
}

func ParseFile(path string) (ParsedFile, error) {
	parsedFile := ParsedFile{
		Path: path,
	}

	source, err := ReadFiles(path)
	if err != nil {
		return parsedFile, err
	}

	querySource, err := os.ReadFile("internal/queries/go.scm")
	if err != nil {
		fmt.Printf("%v", err)
		return parsedFile, err
	}

	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())

	defer parser.Close()

	parser.SetLanguage(language)

	tree := parser.Parse(source, nil)
	defer tree.Close()

	root := tree.RootNode()

	query, queryErr := tree_sitter.NewQuery(language, string(querySource))
	if queryErr != nil {
		return parsedFile, queryErr
	}
	defer query.Close()

	queryCursor := tree_sitter.NewQueryCursor()
	defer queryCursor.Close()

	matches := queryCursor.Matches(query, root, source)

	for {
		match := matches.Next()
		if match == nil {
			break
		}

		if symbol := extractSymbol(match, query, source); symbol != nil {
			parsedFile.Symbols = append(parsedFile.Symbols, *symbol)
		}

		if imp := extractImport(match, query, source); imp != nil {
			parsedFile.Imports = append(parsedFile.Imports, *imp)
		}
	}
	return parsedFile, nil
}
func extractSymbol(
	match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte) *Symbol {
	var defNode *tree_sitter.Node
	var nameNode *tree_sitter.Node

	for _, capture := range match.Captures {
		captureName := query.CaptureNames()[capture.Index]
		switch captureName {
		case "symbol.def":
			defNode = &capture.Node
		case "symbol.name":
			nameNode = &capture.Node
		}
	}
	if defNode == nil || nameNode == nil {
		return nil
	}
	kind, ok := GoSymbolKinds[defNode.Kind()]
	if !ok {
		return nil
	}
	return &Symbol{
		Name:      nameNode.Utf8Text(source),
		Kind:      kind,
		StartLine: uint32(defNode.StartPosition().Row + 1),
		EndLine:   uint32(defNode.EndPosition().Row + 1),
	}
}

func extractImport(match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte) *Import {
	var moduleNode *tree_sitter.Node
	for _, capture := range match.Captures {
		captureName := query.CaptureNames()[capture.Index]
		switch captureName {
		case "import.module":
			moduleNode = &capture.Node
		}
	}
	if moduleNode == nil {
		return nil
	}
	return &Import{
		Path: strings.Trim(moduleNode.Utf8Text(source), `"`),
	}
}

//TODO: CACHE QUERY
//TODO: this rn does consider var_spec as top level and not a local like it does consider var that are local or more important make sure to add the ParentId column in v2
