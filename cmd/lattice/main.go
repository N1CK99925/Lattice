package main

import (
	"Lattice/internal/parser"
	"Lattice/internal/repository"
	"fmt"
	"log"
)

func main() {
	log.Println("Lattice is running")
	log.Println("...")
	var root string
	fmt.Scanf("%s", &root)
	files, err := repository.Walk(root)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		log.Println("Parsing:", file)

		parsedFile, err := parser.ParseFile(file)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("Symbols:")
		for _, symbol := range parsedFile.Symbols {
			log.Println(symbol.Name, symbol.Kind)
		}

		log.Println("Imports:")
		for _, imp := range parsedFile.Imports {
			log.Println(imp.Path)
		}

		log.Println("Calls:")
		for _, call := range parsedFile.Calls {
			log.Printf("%s.%s", call.Receiver, call.Target)
		}
	}
}
