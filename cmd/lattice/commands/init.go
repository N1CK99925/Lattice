package commands

import (
	"context"
	"crypto/sha256"
	"fmt"
	"log"
	"log/slog"
	"path/filepath"

	"Lattice/internal/logger"
	"Lattice/internal/models"
	"Lattice/internal/parser"
	"Lattice/internal/repository"
	"Lattice/internal/resolver"
	"Lattice/internal/storage"

	"github.com/spf13/cobra"
)

var indexCmd = &cobra.Command{
	Use:   "init [path]",
	Short: "Index a Go repository into the local database",
	Long:  `Walks the repository at the given path, parses all Go files, resolves symbol references, and persists the index to SQLite.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runIndex,
}

func init() {
	rootCmd.AddCommand(indexCmd)
}

func runIndex(cmd *cobra.Command, args []string) error {
	logger.Init(logger.Config{
		Level: slog.LevelInfo,
		Json:  false,
	})

	absRoot, err := filepath.Abs(args[0])
	if err != nil {
		return fmt.Errorf("resolve path: %w", err)
	}

	store, err := storage.New("./Lattice.db")
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer store.Close()

	paths, err := repository.Walk(absRoot)
	if err != nil {
		return fmt.Errorf("walk repository: %w", err)
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

		src, err := readFileBytes(abs)
		if err != nil {
			return fmt.Errorf("read %s: %w", rel, err)
		}
		h := sha256.Sum256(src)
		hash := fmt.Sprintf("%x", h)

		cur, err := store.IsCurrent(ctx, abs, hash)
		if err != nil {
			return fmt.Errorf("check %s: %w", rel, err)
		}
		if cur {
			log.Printf("skip (unchanged): %s", rel)
			continue
		}

		pf, err := parser.ParseFile(abs)
		if err != nil {
			return fmt.Errorf("parse %s: %w", rel, err)
		}
		index.Files = append(index.Files, pf)
		toPersist = append(toPersist, fileHash{path: abs, hash: hash})
	}

	resolve := resolver.Resolve(&index)

	for _, fh := range toPersist {
		for _, pf := range index.Files {
			if pf.Path == fh.path {
				log.Printf("persist: %s", filepath.Base(fh.path))
				if err := store.PersistFile(ctx, pf, fh.hash, resolve); err != nil {
					return fmt.Errorf("persist %s: %w", fh.path, err)
				}
				break
			}
		}
	}

	log.Printf("done: %d files parsed, %d persisted", len(index.Files), len(toPersist))
	return nil
}

func readFileBytes(path string) ([]byte, error) {
	return parser.ReadFiles(path)
}
