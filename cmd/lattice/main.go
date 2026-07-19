package main

import (
	"Lattice/internal/repository"
	"log"
)

func main() {
	err := repository.Run_Walker()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Lattice is running")
}
