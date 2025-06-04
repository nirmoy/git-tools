package gittools

// Commit represents a Git commit
type Commit struct {
	Hash    string
	Subject string
	Author  string
	Date    string
}

// Use ASCII unit separator (\x1f) as a safe delimiter for git log output
const LogDelimiter = "\x1f"

// ANSI color codes
const (
	ColorReset  = "\033[0m"
	ColorYellow = "\033[33m"
	ColorGreen  = "\033[32m"
	ColorCyan   = "\033[36m"
) 