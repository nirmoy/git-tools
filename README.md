# Git Tools

A collection of Git utilities for branch and commit management.

## Installation

```bash
go build -o git-tools ./src/
```

## Available Commands

### find-missing
Find commits in one branch that are missing from another branch.

```bash
./git-tools find-missing <branch1> <branch2>
```

**Examples:**
```bash
./git-tools find-missing feature-branch main
./git-tools find-missing main develop
```

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