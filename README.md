# Lattice

Lattice builds a graph of your codebase so agents (LLMs) can get real context about it — what symbols exist, who calls what, and how files depend on each other. It parses Go source with tree-sitter and turns it into a queryable call graph.

## How it works

The pipeline is: walk → parse → resolve → persist.

1. **Walk** — `internal/repository` walks the directory and collects every `.go` file, skipping `.git`.
2. **Parse** — `internal/parser` parses each file with tree-sitter, using the SCM queries embedded in `internal/queries` (via `go:embed`). For every file it extracts:
   - **Symbols** — functions, methods, types, constants, and package-level variables, each with a stable `ID` (`path::kind::receiver::name::line`) and source line range.
   - **Imports** — every imported module.
   - **Calls** — call sites covering simple calls, method calls, package-qualified calls, chained calls, and functions passed as first-class values. Each call is tagged with its enclosing symbol's `ParentSymbolID`, so caller → callee edges are map-able.
   - **Type references** — parameter/receiver types, struct field types, return types, and composite-literal types, which drive file-level `type_use` edges.
3. **Resolve** — `internal/resolver` builds a symbol table (lookup by `ID` and by `Name`) and maps every call target to its definition, producing the resolved call edges.
4. **Persist** — `internal/storage` writes everything into a SQLite database. Files are hashed (sha256), and a file is skipped entirely if its stored hash already matches (i.e. it's unchanged since the last run). Each parsed file is stored in a single transaction: the file row, its symbols (with `defines` edges from a synthetic per-file symbol), plus `imports`, `type_use`, and `calls` edges — calls resolve to the symbol they hit, or fall back to an external target.
5. **Graph** — `internal/graph` is where the final graph representation will live (currently a stub).

## Project structure

```
cmd/lattice/            CLI entry point — reads a repo path, builds the index, resolves and persists it
internal/
  config/               app configuration (stub)
  graph/                graph representation (stub)
  logger/               structured slog logger (JSON or text, level-configurable)
  models/               shared types — symbols, imports, calls, type refs, parsed files
  parser/               tree-sitter parsing — per-file reader, repository indexer
  queries/              embedded tree-sitter SCM query source (go.scm)
  repository/           filesystem walker
  resolver/             symbol table + call resolution
  storage/              SQLite store — migrations, persistence, change detection
```

## Usage

```sh
go build -o lattice ./cmd/lattice
./lattice
# type a path to a Go repository and press Enter
```

The CLI reads the root path from stdin, indexes the repo, and stores everything in `./Lattice.db`. Files whose content hasn't changed since the last run are skipped, so repeat runs are fast.

## Status

Done so far:

- [x] Directory walker (Go files only, `.git` skipped)
- [x] Tree-sitter parsing for Go: symbols, imports, calls, type refs
- [x] Call sites tagged with their enclosing symbol (`ParentSymbolID`) for caller → callee mapping
- [x] Symbol table + call-to-definition resolution
- [x] Structured logging via `slog`
- [x] SQLite persistence with WAL, embedded migrations, and hash-based change detection
- [ ] Graph representation (`internal/graph`)
- [ ] Additional languages beyond Go
- [ ] Proper CLI flags / subcommands
