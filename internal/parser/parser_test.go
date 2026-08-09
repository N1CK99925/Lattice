package parser_test

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"Lattice/internal/parser"
	"Lattice/internal/models"
	"Lattice/internal/repository"
)

// -update rewrites the golden files instead of comparing.
var update = flag.Bool("update", false, "rewrite golden files")

// snapshot is a deterministic, path-relative copy of one parsed file.
// It exists so the golden output is byte-stable across machines and runs.
type snapshot struct {
	Path     string          `json:"path"`
	Symbols  []models.Symbol `json:"symbols"`
	Imports  []models.Import `json:"imports"`
	Calls    []models.Call   `json:"calls"`
	TypeRefs []models.TypeRef `json:"type_refs"`
}

func TestIndexGolden(t *testing.T) {
	root := filepath.Join("testdata", "fixture")

	// Walk + parse using the real pipeline.
	paths, err := repository.Walk(root)
	if err != nil {
		t.Fatalf("walk: %v", err)
	}

	abs, err := filepath.Abs(root)
	if err != nil {
		t.Fatalf("abs: %v", err)
	}

	var snap []snapshot
	for _, p := range paths {
		pf, err := parser.ParseFile(filepath.Join(abs, p))
		if err != nil {
			t.Fatalf("parse %s: %v", p, err)
		}

		// Normalize every path/ID from machine-absolute to fixture-relative,
		// then sort so output order is canonical.
		rel := filepath.ToSlash(p)
		s := snapshot{
			Path:     rel,
			Symbols:  sortSymbols(normalizeSymbols(pf.Symbols, abs)),
			Imports:  sortImports(normalizeImports(pf.Imports, abs)),
			Calls:    sortCalls(normalizeCalls(pf.Calls, abs)),
			TypeRefs: sortTypeRefs(normalizeTypeRefs(pf.TypeRefs, abs)),
		}
		// normalizeCalls must also fix ParentSymbolID strings.
		s.Calls = fixParentIDs(s.Calls, abs)
		snap = append(snap, s)
	}

	got, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	got = append(got, '\n')

	goldenPath := filepath.Join("testdata", "golden", "index.golden")
	if *update {
		if err := os.MkdirAll(filepath.Dir(goldenPath), 0o755); err != nil {
			t.Fatalf("mkdir golden: %v", err)
		}
		if err := os.WriteFile(goldenPath, got, 0o644); err != nil {
			t.Fatalf("write golden: %v", err)
		}
		t.Logf("wrote %s", goldenPath)
		return
	}

	want, err := os.ReadFile(goldenPath)
	if err != nil {
		t.Fatalf("read golden (run with -update first): %v", err)
	}
	if string(got) != string(want) {
		t.Errorf("index diverged from golden.\n--- want\n+++ got\n%s",
			diffLines(string(want), string(got)))
	}
}

// normalizeSymbols rewrites the machine-absolute root in File/ID to a
// fixture-relative prefix, so the snapshot is the same on any machine.
func normalizeSymbols(in []models.Symbol, abs string) []models.Symbol {
	out := make([]models.Symbol, len(in))
	for i, s := range in {
		s.File = strings.TrimPrefix(s.File, abs+"/")
		s.ID = strings.TrimPrefix(s.ID, abs+"/")
		out[i] = s
	}
	return out
}

func normalizeImports(in []models.Import, abs string) []models.Import {
	out := make([]models.Import, len(in))
	for i, im := range in {
		im.File = strings.TrimPrefix(im.File, abs+"/")
		out[i] = im
	}
	return out
}

func normalizeCalls(in []models.Call, abs string) []models.Call {
	out := make([]models.Call, len(in))
	for i, c := range in {
		c.File = strings.TrimPrefix(c.File, abs+"/")
		out[i] = c
	}
	return out
}

func fixParentIDs(in []models.Call, abs string) []models.Call {
	for i := range in {
		in[i].ParentSymbolID = strings.TrimPrefix(in[i].ParentSymbolID, abs+"/")
	}
	return in
}

func sortSymbols(in []models.Symbol) []models.Symbol {
	sort.Slice(in, func(i, j int) bool {
		a, b := in[i], in[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Kind != b.Kind {
			return a.Kind < b.Kind
		}
		if a.Name != b.Name {
			return a.Name < b.Name
		}
		return a.StartLine < b.StartLine
	})
	return in
}

func sortImports(in []models.Import) []models.Import {
	sort.Slice(in, func(i, j int) bool {
		a, b := in[i], in[j]
		if a.File != b.File {
			return a.File < b.File
		}
		return a.Path < b.Path
	})
	return in
}

func sortCalls(in []models.Call) []models.Call {
	sort.Slice(in, func(i, j int) bool {
		a, b := in[i], in[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.StartLine != b.StartLine {
			return a.StartLine < b.StartLine
		}
		return a.Target < b.Target
	})
	return in
}

func normalizeTypeRefs(in []models.TypeRef, abs string) []models.TypeRef {
	out := make([]models.TypeRef, len(in))
	for i, tr := range in {
		tr.File = strings.TrimPrefix(tr.File, abs+"/")
		out[i] = tr
	}
	return out
}

func sortTypeRefs(in []models.TypeRef) []models.TypeRef {
	sort.Slice(in, func(i, j int) bool {
		a, b := in[i], in[j]
		if a.StartLine != b.StartLine {
			return a.StartLine < b.StartLine
		}
		return a.Text < b.Text
	})
	return in
}

// diffLines prints a minimal +/- diff of two multi-line strings.
func diffLines(want, got string) string {
	ws := strings.Split(want, "\n")
	gs := strings.Split(got, "\n")
	var b strings.Builder
	max := len(ws)
	if len(gs) > max {
		max = len(gs)
	}
	for i := 0; i < max; i++ {
		var w, g string
		if i < len(ws) {
			w = ws[i]
		}
		if i < len(gs) {
			g = gs[i]
		}
		switch {
		case i >= len(ws):
			b.WriteString("+ " + g + "\n")
		case i >= len(gs):
			b.WriteString("- " + w + "\n")
		case w != g:
			b.WriteString("- " + w + "\n")
			b.WriteString("+ " + g + "\n")
		}
	}
	return b.String()
}