# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Setup (installs go, goreleaser, golangci-lint, gofumpt, vhs via mise)
mise install
make setup          # go mod tidy

# Test
make test           # go test -v ./... -count=1
go test -v ./app/... -count=1   # single package
make coverage       # generate coverage.txt (excludes uitests)
make coverage-html  # open coverage report in browser

# Lint & Format
make lint           # golangci-lint run
make format         # gofumpt -l -w .

# Build
make build          # goreleaser build --single-target --snapshot --output ./jqcompletion

# Manual UI tests (requires a terminal; sets DEBUGLOG=1)
make testuijsonview
make testuiqueryview
```

Enable debug logging by setting `DEBUGLOG=1`; logs are written to `debug.log`.

## Architecture

This is a [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI app (Elm architecture) with a split-pane layout:

```
┌─────────────────┬─────────────────┐
│   queryview     │    jsonview     │
│  (left panel)   │  (right panel)  │
│  text input +   │  scrollable     │
│  fuzzy list     │  JSON output    │
└─────────────────┴─────────────────┘
```

### Packages

- **`app/` (package `main`)** — Top-level bubbletea model and entry point
  - `main.go`: CLI arg parsing, optional editor launch (when no file given), bubbletea `Program` setup. On exit, prints the final jq query.
  - `app.go`: Root `model` struct composing `queryview.Model` and `jsonview.Model`. Handles `WindowSizeMsg` (splits width 50/50), `SelectionChangedMsg` (updates JSON output on list nav), `KeyEnter` (evaluates current query), `KeyCtrlS` (saves query to history and opens result read-only in `$EDITOR`).
  - `json.go`: JSON key extraction via `gojq` (`GetUnnestedKeys` uses `paths` to enumerate all dot-paths), jq query execution (`QueryJsonData`, `RobustQueryJsonData`), and `KeySearchEngine` (implements `queryview.Engine` using `sahilm/fuzzy`). `RobustQueryJsonData` falls back to progressively shorter prefix queries if the exact query fails.
  - `editor.go`: Editor integration using `$EDITOR` (default: vim). Handles temp file creation, file copy/load/save.

- **`queryview/`** — Left panel: text input with debounced fuzzy completion list
  - Implements `Engine` interface (`Query(string) []string`) for pluggable search backends.
  - Debounces typing by 100ms before issuing `queryRequestMsg` → `queryResponseMsg` cycle.
  - Emits `SelectionChangedMsg` when the list selection changes (on nav or new results).

- **`jsonview/`** — Right panel: scrollable viewport rendering pretty-printed JSON.

- **`uitests/`** — Standalone manual UI test programs for `jsonview` and `queryview` in isolation (not unit tests; run interactively).

### Data flow

1. User types in `queryview` input → debounced `queryRequestMsg` → `KeySearchEngine.Query()` (fuzzy match over all JSON dot-paths) → list updated.
2. User navigates list → `SelectionChangedMsg` → `app` calls `RobustQueryJsonData` → `jsonview` displays result.
3. User presses Enter → uses raw text input value as jq query.
4. User presses Ctrl+S → appends current query to history, saves filtered JSON to temp file, opens in `$EDITOR` (read-only for vim/nvim).
