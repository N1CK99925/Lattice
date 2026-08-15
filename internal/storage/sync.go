package storage

import (
	"Lattice/internal/models"
	"Lattice/internal/storage/generated"
	"context"
	"database/sql"
	"path/filepath"
	"time"
)

func (s *Store) IsCurrent(ctx context.Context, path, hash string) (bool, error) {
	f, err := s.Queries.GetFileByPath(ctx, path)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return f.Hash == hash, nil
}

type ResolveFunc func(models.Call) (string, bool)

func (s *Store) PersistFile(ctx context.Context, pf models.ParsedFile, hash string, resolve ResolveFunc) error {
	tx, err := s.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	q := s.Queries.WithTx(tx)

	fileID, err := q.UpsertFile(ctx, generated.UpsertFileParams{
		Path:     pf.Path,
		Language: "go",
		Hash:     hash,
		ParsedAt: time.Now(),
	})
	if err != nil {
		return err
	}

	if err := q.DeleteEdgesByFile(ctx, generated.DeleteEdgesByFileParams{
		FileID:   fileID,
		FileID_2: fileID,
	}); err != nil {
		return err
	}

	if err := q.DeleteSymbolsByFile(ctx, fileID); err != nil {
		return err
	}

	fileSymID := pf.Path + "::file"
	if err := q.InsertSymbol(ctx, generated.InsertSymbolParams{
		ID:        fileSymID,
		FileID:    fileID,
		Name:      filepath.Base(pf.Path),
		Receiver:  sql.NullString{},
		Kind:      "file",
		StartLine: 1,
		EndLine:   1,
		Exported:  false,
		Scope:     "package",
	}); err != nil {
		return err
	}

	for _, sym := range pf.Symbols {
		if err := q.InsertSymbol(ctx, generated.InsertSymbolParams{
			ID:        sym.ID,
			FileID:    fileID,
			Name:      sym.Name,
			Receiver:  nullStr(sym.Receiver),
			Kind:      string(sym.Kind),
			StartLine: int64(sym.StartLine),
			EndLine:   int64(sym.EndLine),
			Exported:  len(sym.Name) > 0 && sym.Name[0] >= 'A' && sym.Name[0] <= 'Z',
			Scope:     sym.Scope,
		}); err != nil {
			return err
		}
		if err := q.InsertEdge(ctx, generated.InsertEdgeParams{
			SourceSymbol: fileSymID,
			TargetSymbol: sql.NullString{String: sym.ID, Valid: true},
			Kind:         "defines",
		}); err != nil {
			return err
		}
	}
	for _, imp := range pf.Imports {
		if err := q.InsertEdge(ctx, generated.InsertEdgeParams{
			SourceSymbol:   fileSymID,
			TargetExternal: sql.NullString{String: imp.Path, Valid: true},
			Kind:           "imports",
		}); err != nil {
			return err
		}
	}

	for _, tr := range pf.TypeRefs {
		if err := q.InsertEdge(ctx, generated.InsertEdgeParams{
			SourceSymbol:   fileSymID,
			TargetExternal: sql.NullString{String: tr.Text, Valid: true},
			Kind:           "type_use",
		}); err != nil {
			return err
		}
	}

	for _, call := range pf.Calls {
		caller := call.ParentSymbolID
		if caller == "" {
			continue
		}

		if resolve != nil {
			if targetID, ok := resolve(call); ok {
				if err := q.InsertEdge(ctx, generated.InsertEdgeParams{
					SourceSymbol: caller,
					TargetSymbol: sql.NullString{String: targetID, Valid: true},
					Kind:         "calls",
				}); err != nil {
					return err
				}
				continue
			}
		}
		external := call.Target
		if call.Receiver != "" {
			external = call.Receiver + "." + call.Target
		}
		if err := q.InsertEdge(ctx, generated.InsertEdgeParams{
			SourceSymbol:   caller,
			TargetExternal: sql.NullString{String: external, Valid: true},
			Kind:           "calls",
		}); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func nullStr(s string) sql.NullString {
	if s == "" {
		return sql.NullString{}
	}
	return sql.NullString{String: s, Valid: true}
}
