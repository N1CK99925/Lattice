package parser

import (
	"fmt"
	"os"

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

func ParseFile(path string) error {
	source, err := ReadFiles(path)
	if err != nil {
		return err
	}
	parser := tree_sitter.NewParser()
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())
	defer parser.Close()
	parser.SetLanguage(language)

	tree := parser.Parse(source, nil)
	defer tree.Close()

	root := tree.RootNode()
	fmt.Println(root.ToSexp())

	return nil
}
