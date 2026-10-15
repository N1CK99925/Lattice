package commands

import (
	"bytes"
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"

	"log/slog"

	"Lattice/internal/logger"
	"Lattice/internal/storage"
	"Lattice/internal/storage/generated"
)

const fixtureDir = "../../../go_code_test"

// TestMain initializes the package logger so storage.New (which logs on open)
// works in tests, the same way the CLI commands do before touching a store.
func TestMain(m *testing.M) {
	logger.Init(logger.Config{Level: slog.LevelInfo})
	os.Exit(m.Run())
}

// fixtureAbs is the absolute path of the go_code_test directory, resolved
// from the package directory (go test runs with cwd set to the package).
var fixtureAbs = func() string {
	abs, err := filepath.Abs(fixtureDir)
	if err != nil {
		panic(err)
	}
	return abs
}()

type fixtureSym struct {
	file string
	name string
	line int
}

// fixtureSymbols describes the symbols the parser produces for the files in
// go_code_test. The IDs use the same `<path>::<kind>::<receiver>::<name>::<line>`
// format the parser emits.
var fixtureSymbols = []fixtureSym{
	{file: "main.go", name: "A", line: 3},
	{file: "main.go", name: "B", line: 7},
	{file: "main.go", name: "C", line: 11},
	{file: "main.go", name: "D", line: 15},
	{file: "main.go", name: "E", line: 19},
	{file: "main.go", name: "F", line: 23},
	{file: "main.go", name: "G", line: 27},
	{file: "main.go", name: "H", line: 31},
	{file: "main.go", name: "I", line: 35},
	{file: "main.go", name: "J", line: 39},
	{file: "main.go", name: "K", line: 43},
	{file: "extra.go", name: "L", line: 3},
	{file: "extra.go", name: "M", line: 7},
	{file: "extra.go", name: "N", line: 11},
	{file: "extra.go", name: "O", line: 15},
	{file: "extra.go", name: "P", line: 19},
}

// callEdges is the resolved call graph inside go_code_test.
var callEdges = [][2]string{
	{"A", "B"},
	{"B", "C"},
	{"C", "D"},
	{"D", "E"},
	{"E", "F"},
	{"F", "G"},
	{"G", "H"},
	{"H", "I"},
	{"I", "J"},
	{"J", "K"},
	{"L", "C"},
	{"M", "E"},
	{"N", "G"},
	{"O", "I"},
	{"P", "C"},
}

func symbolID(t *testing.T, name string) string {
	t.Helper()
	for _, s := range fixtureSymbols {
		if s.name == name {
			return filepath.Join(fixtureAbs, s.file) + "::function::::" + s.name + "::" + strconv.Itoa(s.line)
		}
	}
	t.Fatalf("unknown fixture symbol %q", name)
	return ""
}

func ids(t *testing.T, names ...string) []string {
	t.Helper()
	out := make([]string, 0, len(names))
	for _, n := range names {
		out = append(out, symbolID(t, n))
	}
	return out
}

func chdirTemp(t *testing.T) {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(t.TempDir()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(wd) })
}

// captureStdout redirects os.Stdout for the duration of fn and returns
// everything the process wrote to stdout.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	done := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = buf.ReadFrom(r)
		done <- buf.String()
	}()

	fn()
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	out := <-done
	os.Stdout = old
	return out
}

// symbolLines keeps only the lines that look like symbol IDs (logger and
// "depth: N" lines are filtered out).
func symbolLines(out string) []string {
	var got []string
	for _, ln := range strings.Split(out, "\n") {
		if strings.Contains(ln, "::") {
			got = append(got, ln)
		}
	}
	return got
}

// runDeps runs the deps CLI handler for symbol at the given depth.
func runDeps(t *testing.T, symbol string, depth int) []string {
	t.Helper()
	depsCmd.SetContext(context.Background())
	if err := depsCmd.Flags().Set("depth", strconv.Itoa(depth)); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() {
		if err := getDeps(depsCmd, []string{symbol}); err != nil {
			t.Errorf("getDeps: %v", err)
		}
	})
	return symbolLines(out)
}

// runDependents runs the dependents CLI handler for symbol at the given depth.
func runDependents(t *testing.T, symbol string, depth int) []string {
	t.Helper()
	dependentsCmd.SetContext(context.Background())
	if err := dependentsCmd.Flags().Set("depth", strconv.Itoa(depth)); err != nil {
		t.Fatal(err)
	}
	out := captureStdout(t, func() {
		if err := getDependents(dependentsCmd, []string{symbol}); err != nil {
			t.Errorf("getDependents: %v", err)
		}
	})
	return symbolLines(out)
}

// seedFixtureDB creates ./Lattice.db (in the current working directory) and
// loads it with the fixture symbols and the call edges between them, so the
// CLI's graph loading path is exercised against a known graph.
func seedFixtureDB(t *testing.T) {
	t.Helper()
	store, err := storage.New("Lattice.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	ctx := context.Background()

	fileIDs := map[string]int64{}
	for _, file := range []string{"main.go", "extra.go"} {
		id, err := store.Queries.UpsertFile(ctx, generated.UpsertFileParams{
			Path:     filepath.Join(fixtureAbs, file),
			Language: "go",
			Hash:     "fixture",
			ParsedAt: time.Now(),
		})
		if err != nil {
			t.Fatal(err)
		}
		fileIDs[file] = id
	}

	for _, s := range fixtureSymbols {
		if err := store.Queries.InsertSymbol(ctx, generated.InsertSymbolParams{
			ID:        symbolID(t, s.name),
			FileID:    fileIDs[s.file],
			Name:      s.name,
			Kind:      "function",
			StartLine: int64(s.line),
			EndLine:   int64(s.line),
			Exported:  true,
			Scope:     "package",
		}); err != nil {
			t.Fatal(err)
		}
	}

	for _, e := range callEdges {
		if err := store.Queries.InsertEdge(ctx, generated.InsertEdgeParams{
			SourceSymbol: symbolID(t, e[0]),
			TargetSymbol: sql.NullString{String: symbolID(t, e[1]), Valid: true},
			Kind:         "calls",
		}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestDeps(t *testing.T) {
	chdirTemp(t)
	seedFixtureDB(t)

	a := symbolID(t, "A")
	p := symbolID(t, "P")
	k := symbolID(t, "K")

	tests := []struct {
		name   string
		symbol string
		depth  int
		want   []string
	}{
		{
			name:   "full chain",
			symbol: a,
			depth:  defaultDepth(),
			want:   ids(t, "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"),
		},
		{
			name:   "depth one stops at direct deps",
			symbol: a,
			depth:  1,
			want:   ids(t, "A", "B"),
		},
		{
			name:   "depth zero returns only the symbol itself",
			symbol: a,
			depth:  0,
			want:   ids(t, "A"),
		},
		{
			name:   "cross file dependency",
			symbol: p,
			depth:  defaultDepth(),
			want:   ids(t, "P", "C", "D", "E", "F", "G", "H", "I", "J", "K"),
		},
		{
			name:   "leaf has no dependencies",
			symbol: k,
			depth:  defaultDepth(),
			want:   ids(t, "K"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runDeps(t, tt.symbol, tt.depth)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("deps %s (depth %d) = %v, want %v", tt.symbol, tt.depth, got, tt.want)
			}
		})
	}
}

func TestDependents(t *testing.T) {
	chdirTemp(t)
	seedFixtureDB(t)

	c := symbolID(t, "C")
	k := symbolID(t, "K")
	a := symbolID(t, "A")

	tests := []struct {
		name   string
		symbol string
		depth  int
		want   []string
	}{
		{
			name:   "everything that transitively calls",
			symbol: c,
			depth:  defaultDepth(),
			want:   ids(t, "C", "B", "L", "P", "A"),
		},
		{
			name:   "depth one stops at direct callers",
			symbol: c,
			depth:  1,
			want:   ids(t, "C", "B", "L", "P"),
		},
		{
			name:   "leaf gathers the whole graph",
			symbol: k,
			depth:  defaultDepth(),
			want:   ids(t, "K", "J", "I", "H", "O", "G", "F", "N", "E", "D", "M", "C", "B", "L", "P", "A"),
		},
		{
			name:   "depth stops the reverse walk",
			symbol: k,
			depth:  2,
			want:   ids(t, "K", "J", "I"),
		},
		{
			name:   "uncalled symbol has no dependents",
			symbol: a,
			depth:  defaultDepth(),
			want:   ids(t, "A"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := runDependents(t, tt.symbol, tt.depth)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("dependents %s (depth %d) = %v, want %v", tt.symbol, tt.depth, got, tt.want)
			}
		})
	}
}

// defaultDepth returns the CLI's default depth flag value.
func defaultDepth() int {
	d, err := depsCmd.Flags().GetInt("depth")
	if err != nil {
		return 5
	}
	return d
}

// functionName extracts the function name from a parser-style symbol ID
// (<path>::function::::<name>::<line>), falling back to "" for non-function
// symbols (e.g. the per-file "defines" symbol).
func functionName(id string) string {
	segments := strings.Split(id, "::")
	if len(segments) < 5 || segments[len(segments)-4] != "function" {
		return ""
	}
	return segments[len(segments)-2]
}

// functionNames returns the sorted, de-duplicated function names found in
// the printed symbol lines.
func functionNames(t *testing.T, lines []string) []string {
	t.Helper()
	seen := map[string]bool{}
	var names []string
	for _, ln := range lines {
		if name := functionName(ln); name != "" && !seen[name] {
			seen[name] = true
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names
}

// TestDepsIntegration drives the whole pipeline against the real go_code_test
// sources: index (parse + resolve + persist) and then ask for the deps of one
// of the symbols through the CLI handler. This is the path a real user takes
// after running `lattice init`.
func TestDepsIntegration(t *testing.T) {
	abs := filepath.Join(fixtureAbs)
	chdirTemp(t)

	indexCmd.SetContext(context.Background())
	_ = captureStdout(t, func() {
		if err := runIndex(indexCmd, []string{abs}); err != nil {
			t.Errorf("runIndex: %v", err)
		}
	})

	got := runDeps(t, symbolID(t, "A"), defaultDepth())
	if !reflect.DeepEqual(got, ids(t, "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K")) {
		t.Fatalf("integration deps = %v, want %v", got, ids(t, "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"))
	}
}

// TestDependentsIntegration drives the whole pipeline as well, but asserts on
// the set of function callers because the real index also wires in per-file
// "defines" edges that the seeded tests do not contain.
func TestDependentsIntegration(t *testing.T) {
	abs := filepath.Join(fixtureAbs)
	chdirTemp(t)

	indexCmd.SetContext(context.Background())
	_ = captureStdout(t, func() {
		if err := runIndex(indexCmd, []string{abs}); err != nil {
			t.Errorf("runIndex: %v", err)
		}
	})

	got := functionNames(t, runDependents(t, symbolID(t, "C"), defaultDepth()))
	want := []string{"A", "B", "C", "L", "P"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("integration dependents = %v, want %v", got, want)
	}
}
