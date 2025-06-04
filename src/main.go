package gittools

import (
	"fmt"
	"os"
)

func RunCLI() {
	if len(os.Args) < 2 {
		PrintUsage()
		os.Exit(1)
	}

	subcmd := os.Args[1]

	switch subcmd {
	case "find-missing":
		if len(os.Args) != 4 {
			fmt.Fprintf(os.Stderr, "Usage: %s find-missing <branch1> <branch2>\n", os.Args[0])
			os.Exit(1)
		}
		FindMissing(os.Args[2], os.Args[3])
	case "grep-branch":
		if len(os.Args) == 3 {
			GrepBranch(false, os.Args[2])
		} else if len(os.Args) == 4 && os.Args[2] == "--all" {
			GrepBranch(true, os.Args[3])
		} else {
			fmt.Fprintf(os.Stderr, "Usage: %s grep-branch [--all] \"text\"\n", os.Args[0])
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcmd)
		PrintUsage()
		os.Exit(1)
	}
}

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  git-tools find-missing <branch1> <branch2>   # Find commits in branch1 missing from branch2")
	fmt.Println("  git-tools grep-branch [--all] \"text\"         # List branches and commits where text exists in commit message")
	fmt.Println("    --all: search all refs (branches, remotes, tags)")
} 