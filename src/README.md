# Git Tools - Source Code Organization

This directory contains the reorganized source code for the git-tools application, with each subcommand separated into its own file for better maintainability.

## File Structure

```
src/
├── main.go           # Main entry point and command routing
├── types.go          # Shared data structures and constants
├── utils.go          # Common utility functions
├── find_missing.go   # Implementation of the 'find-missing' subcommand
├── grep_branch.go    # Implementation of the 'grep-branch' subcommand
└── README.md         # This file
```

## File Descriptions

### `main.go`
- Contains the main function and command-line argument parsing
- Routes subcommands to their respective handler functions
- Contains the `printUsage()` function for displaying help information

### `types.go`
- Defines shared data structures like the `Commit` struct
- Contains constants used across multiple files (color codes, delimiters)

### `utils.go` 
- Contains utility functions used by multiple subcommands:
  - `isGitRepo()` - checks if current directory is a Git repository
  - `branchExists()` - checks if a Git branch exists
  - `normalizeSubject()` - normalizes commit subject strings

### `find_missing.go`
- Implements the `find-missing` subcommand functionality
- Contains functions specific to finding missing commits between branches:
  - `findMissing()` - main handler function
  - `getMissingCommits()` - retrieves commits missing from target branch
  - `getAllSubjects()` - gets all commit subjects from a branch

### `grep_branch.go`
- Implements the `grep-branch` subcommand functionality
- Contains the `grepBranch()` function for searching commit messages across branches

## Benefits of This Organization

1. **Separation of Concerns**: Each subcommand has its own file, making the code easier to navigate and maintain
2. **Shared Resources**: Common types, constants, and utilities are centralized in dedicated files
3. **Scalability**: Adding new subcommands is as simple as creating a new file and updating the router in `main.go`
4. **Maintainability**: Bugs and features for specific subcommands can be addressed in isolation
5. **Clean Main**: The main function is now focused solely on argument parsing and routing

## Building

To build the project:
```bash
go build -o ../git-tools .
``` 