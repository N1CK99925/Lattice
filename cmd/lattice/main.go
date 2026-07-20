package main

import (
	"Lattice/internal/parser"
	"Lattice/internal/repository"
	"log"
)

func main() {
	log.Println("Lattice is running")
	log.Println("...")
	files, err := repository.Run_Walker()
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		err := parser.ParseFile(file)
		if err != nil {
			log.Fatal(err)
		}
	}
}
