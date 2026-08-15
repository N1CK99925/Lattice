package main

import (
	"Lattice/internal/logger"
	"Lattice/internal/models"
	"Lattice/internal/parser"
	"Lattice/internal/repository"
	"Lattice/internal/resolver"
	"Lattice/internal/storage"
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

func main() {
	logger.Init(logger.Config{
		Level: slog.LevelInfo,
		Json:  true,
	})
	log.Println("Lattice is running")

	store, err := storage.New("./Lattice.db")
	if err != nil {
		log.Fatal(err)
	}

	var root string
	fmt.Scanf("%s", &root)

	absRoot, err := filepath.Abs(root)
	if err != nil {
		log.Fatal(err)
	}

	paths, err := repository.Walk(absRoot)
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	var index models.RepositoryIndex
	type fileHash struct {
		path string
		hash string
	}
	var toPersist []fileHash

	for _, rel := range paths {
		abs := filepath.Join(absRoot, rel)

		src, err := os.ReadFile(abs)
		if err != nil {
			log.Fatal(err)
		}
		h := sha256.Sum256(src)
		hash := fmt.Sprintf("%x", h)

		cur, err := store.IsCurrent(ctx, abs, hash)
		if err != nil {
			log.Fatal(err)
		}
		if cur {
			log.Printf("skip (unchanged): %s", rel)
			continue
		}

		pf, err := parser.ParseFile(abs)
		if err != nil {
			log.Fatal(err)
		}
		index.Files = append(index.Files, pf)
		toPersist = append(toPersist, fileHash{path: abs, hash: hash})
	}

	// table := resolver.BuildSymbolTable(&index)
	// resolve := func(target string) (string, bool) {
	// 	syms, ok := table.ByName[target]
	// 	if !ok || len(syms) == 0 {
	// 		return "", false
	// 	}
	// 	return syms[0].ID, true
	// }
	resolve := resolver.Resolve(&index)
	for _, fh := range toPersist {
		for _, pf := range index.Files {
			if pf.Path == fh.path {
				log.Printf("persist: %s", filepath.Base(fh.path))
				if err := store.PersistFile(ctx, pf, fh.hash, resolve); err != nil {
					log.Fatalf("persist %s: %v", fh.path, err)
				}
				break
			}
		}
	}

	log.Printf("done: %d files parsed, %d persisted", len(index.Files), len(toPersist))
}
