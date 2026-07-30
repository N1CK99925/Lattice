package main

import (
	"Lattice/internal/parser"
	"fmt"
	"log"
)

func main() {
	log.Println("Lattice is running")

	var root string
	fmt.Scanf("%s", &root)

	index, err := parser.BuildRepositoryIndex(root)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range index.Files {
		log.Println("Parsing:", file.Path)

		log.Println("Symbols:")
		for _, symbol := range file.Symbols {
			log.Println(symbol.Name, symbol.Kind)
		}

		log.Println("Imports:")
		for _, imp := range file.Imports {
			log.Println(imp.Path)
		}

		log.Println("Calls:")
		for _, call := range file.Calls {
			log.Printf("%s.%s", call.Receiver, call.Target)
		}
	}
}
