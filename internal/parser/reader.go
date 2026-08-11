package parser

import (
	"Lattice/internal/models"
	cached "Lattice/internal/queries"
	_ "embed"
	"os"
	"strconv"
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

func ParseFile(path string) (models.ParsedFile, error) {

	parsedFile := models.ParsedFile{
		Path: path,
	}

	source, err := ReadFiles(path)
	if err != nil {
		return parsedFile, err
	}
	querySource := cached.QuerySource

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
		ctx := models.ParseContext{
			Path: path,
		}

		if symbol := extractSymbol(match, query, source, ctx); symbol != nil {
			parsedFile.Symbols = append(parsedFile.Symbols, *symbol)
		}

		if imp := extractImport(match, query, source, ctx); imp != nil {
			parsedFile.Imports = append(parsedFile.Imports, *imp)
		}
		if call := extractCall(match, query, source, ctx); call != nil {
			parsedFile.Calls = append(parsedFile.Calls, *call)
		}
		if typeRef := extractTypeRef(match, query, source, ctx); typeRef != nil {
			parsedFile.TypeRefs = append(parsedFile.TypeRefs, *typeRef)
		}
	}
	return parsedFile, nil
}
func symbolScope(defNode *tree_sitter.Node) string {
	for n := defNode.Parent(); n != nil; n = n.Parent() {
		switch n.Kind() {
		case "function_declaration", "method_declaration":
			return "local"
		case "source_file":
			return "package"
		}
	}
	return "package"
}

func extractSymbol(
	match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte, ctx models.ParseContext) *models.Symbol {
	var defNode *tree_sitter.Node
	var nameNode *tree_sitter.Node
	var receiverNode *tree_sitter.Node

	for _, capture := range match.Captures {
		captureName := query.CaptureNames()[capture.Index]
		switch captureName {
		case "symbol.def":
			defNode = &capture.Node
		case "symbol.name":
			nameNode = &capture.Node
		case "symbol.receiver":
			receiverNode = &capture.Node
		}
	}
	if defNode == nil || nameNode == nil {
		return nil
	}

	receiver := ""
	if receiverNode != nil {
		receiver = receiverNode.Utf8Text(source)
	}
	kind, ok := models.GoSymbolKinds[defNode.Kind()]
	if !ok {
		return nil
	}
	return &models.Symbol{
		Name:      nameNode.Utf8Text(source),
		File:      ctx.Path,
		Kind:      kind,
		Receiver:  receiver,
		ID:        ctx.Path + "::" + string(kind) + "::" + receiver + "::" + nameNode.Utf8Text(source) + "::" + strconv.Itoa(int(defNode.StartPosition().Row+1)),
		StartLine: uint32(defNode.StartPosition().Row + 1),
		EndLine:   uint32(defNode.EndPosition().Row + 1),
		Scope:     symbolScope(defNode),
	}

}

func extractImport(match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte, ctx models.ParseContext) *models.Import {
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
	return &models.Import{
		Path: strings.Trim(moduleNode.Utf8Text(source), `"`),
		File: ctx.Path,
	}
}

// This symbol is used to get the parent of the call , so u can map caller -> calee
func enclosingSymbolID(node *tree_sitter.Node, source []byte, ctx models.ParseContext) string {
	for node != nil {
		switch node.Kind() {
		case "function_declaration", "method_declaration":
			name := node.ChildByFieldName("name")
			if name == nil {
				break
			}
			receiver := ""
			if node.Kind() == "method_declaration" {
				if r := node.ChildByFieldName("receiver"); r != nil {
					receiver = r.Utf8Text(source)
				}
			}
			kind := models.GoSymbolKinds[node.Kind()]
			return ctx.Path + "::" + string(kind) + "::" + receiver + "::" + name.Utf8Text(source) + "::" + strconv.Itoa(int(node.StartPosition().Row+1))
		}

		node = node.Parent()
	}

	return ""
}
func extractCall(match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte, ctx models.ParseContext) *models.Call {
	var siteNode *tree_sitter.Node
	var targetNode *tree_sitter.Node
	var receiverNode *tree_sitter.Node
	// var argumentsNode *tree_sitter.Node
	for _, capture := range match.Captures {
		switch query.CaptureNames()[capture.Index] {
		case "call.site":
			siteNode = &capture.Node

		case "call.target":
			targetNode = &capture.Node

		case "call.receiver":
			receiverNode = &capture.Node

			// case "call.arguments":
			// argumentsNode = &capture.Node
		}
	}
	if siteNode == nil || targetNode == nil {
		return nil
	}
	receiver := ""

	if receiverNode != nil {
		receiver = receiverNode.Utf8Text(source)
	}

	return &models.Call{
		ParentSymbolID: enclosingSymbolID(siteNode, source, ctx),
		Target:         targetNode.Utf8Text(source),
		File:           ctx.Path,
		Receiver:       receiver,
		StartLine:      uint32(siteNode.StartPosition().Row + 1),
		EndLine:        uint32(siteNode.EndPosition().Row + 1),
	}
}

func extractTypeRef(match *tree_sitter.QueryMatch, query *tree_sitter.Query, source []byte, ctx models.ParseContext) *models.TypeRef {
	var typeNode *tree_sitter.Node
	for _, capture := range match.Captures {
		captureName := query.CaptureNames()[capture.Index]
		switch captureName {
		case "param.type":
			typeNode = &capture.Node
		}
	}
	if typeNode == nil {
		return nil
	}
	return &models.TypeRef{
		Text:      typeNode.Utf8Text(source),
		File:      ctx.Path,
		StartLine: uint32(typeNode.StartPosition().Row + 1),
		EndLine:   uint32(typeNode.EndPosition().Row + 1),
	}
}

//TODO: this rn does consider var_spec as top level and not a local like it does consider var that are local or more important make sure to add the ParentId column in v2
