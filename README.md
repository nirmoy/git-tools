# Git Tools

A collection of Git utilities for branch and commit management with interactive terminal interfaces.

## Installation

```bash
go build -o git-tools ./src/
```

## Available Commands

### find-missing
Find commits in one branch that are missing from another branch.

```bash
./git-tools find-missing [--browse|-i] [--tui|-t] <branch1> <branch2>
```

**Display Modes:**

**Normal output** (default):
```bash
./git-tools find-missing feature-branch main
./git-tools find-missing main develop
```

**Interactive browse mode** (with pager):
```bash
./git-tools find-missing --browse origin/feature main
./git-tools find-missing -i feature main
```

**Terminal User Interface** (dual-pane with full git show):
```bash
./git-tools find-missing --tui origin/feature main
./git-tools find-missing -t feature main
```

**TUI Features:**
- **Dual-pane layout**: Commit list (left) and patch viewer (right)
- **Full git show output**: Complete commit details with colored diffs
- **Navigation**: ↑↓/jk (navigate), ←→/hl (horizontal scroll), PgUp/PgDn/Space (scroll patches)
- **Responsive design**: Header wraps in narrow terminals, horizontal scrolling for long commits
- **Keyboard shortcuts**: Enter (focus patch), Escape (back to list), q (quit)

### grep-branch
Search for text in commit messages across branches.

```bash
./git-tools grep-branch [--all] "search text"
```

**Options:**
- `--all`: Search all refs (branches, remotes, tags) instead of just local branches

**Examples:**
```bash
./git-tools grep-branch "fix bug"
./git-tools grep-branch --all "authentication"
```

## Requirements

- Git must be installed and available in PATH
- Must be run from within a Git repository
- Terminal with color support recommended for best TUI experience 