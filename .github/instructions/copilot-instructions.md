# Copilot Instructions for Gogue Roguelike Game (Go)

## Project Overview
- **Console roguelike game** in Go 1.22+, inspired by Rogue (1980), using the [goncurses](https://github.com/rthornton128/goncurses) library for terminal UI.
- **Strict multi-layered architecture**:
  - `model/entity/` — domain/game logic (Player, Enemy, Level, Room, Item, etc.)
  - `presentation/cli/` — UI logic, state management, and goncurses rendering
  - `model/datalayer/` — data persistence (save/load, leaderboard)
  - `view/cli/` — low-level rendering helpers for goncurses
  - `cmd/gogue-cli/` — main entrypoint

## Key Architectural Patterns
- **Domain logic is framework-agnostic**: No goncurses or I/O in `model/entity/`.
- **Presentation layer** (`presentation/cli/`) manages game state, user input, and delegates rendering to `view/cli/`.
- **Data layer** (`model/datalayer/`) handles all persistent storage (JSON for saves/leaderboard).
- **Dependency inversion**: Domain does not import presentation or data layers.
- **State pattern**: Game and menu states are managed as a stack (`States []state.State`).

## Developer Workflows
- **Build**: `make build` (outputs to `bin/gogue-cli.exe`)
- **Test**: `make test` (runs all Go tests)
- **Coverage**: `make test-coverage` (outputs HTML to `bin/coverage.html`)
- **Lint**: `make lint` (requires `golangci-lint`)
- **Run**: `go run ./cmd/gogue-cli/main.go` (ensure terminal supports goncurses)

## Project Conventions
- **All code lives in `src/`**. Do not place Go code outside this directory.
- **Strict separation of concerns**: UI, domain, and data logic must not mix.
- **Game state is turn-based**: All world updates are triggered by player actions.
- **Entities**: See `model/entity/` for canonical struct definitions (Player, Enemy, Backpack, etc.).
- **Menu and state transitions**: See `presentation/cli/state/` and `model/menu.go`.
- **Rendering**: All goncurses drawing is in `view/cli/`.
- **Controls**: WASD for movement, h/j/k/e for item use (see README for full mapping).

## Integration Points
- **goncurses**: Only used in presentation/view layers.
- **Persistent data**: JSON files, handled in `model/datalayer/`.
- **No direct cross-layer imports**: Use interfaces and dependency injection where needed.

## Examples
- **Add a new enemy type**: Implement in `model/entity/enemy.go`, update logic in domain, rendering in view, and spawning in level generation.
- **Add a new menu option**: Update `model/menu.go` and `presentation/cli/state/main_menu.go`.

## References
- See `README.md` for full gameplay and architecture requirements.
- Example C code for algorithms in `code-samples/rogue_sample/`.

---
If any section is unclear or incomplete, please provide feedback for further iteration.
