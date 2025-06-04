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
		handleFindMissing()
	case "grep-branch":
		handleGrepBranch()
	default:
		fmt.Fprintf(os.Stderr, "Unknown subcommand: %s\n", subcmd)
		PrintUsage()
		os.Exit(1)
	}
}

func handleFindMissing() {
	args := os.Args[2:]
	interactive := false
	tui := false
	
	// Check for interactive flags
	var branches []string
	for _, arg := range args {
		if arg == "--browse" || arg == "-i" || arg == "--interactive" {
			interactive = true
		} else if arg == "--tui" || arg == "-t" {
			tui = true
		} else {
			branches = append(branches, arg)
		}
	}
	
	if len(branches) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s find-missing [--browse|-i] [--tui|-t] <branch1> <branch2>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  --browse, -i, --interactive: Browse commits interactively with detailed view\n")
		fmt.Fprintf(os.Stderr, "  --tui, -t: Launch Terminal User Interface (dual-pane view)\n")
		os.Exit(1)
	}
	
	if tui {
		FindMissingTUI(branches[0], branches[1])
	} else if interactive {
		FindMissingInteractive(branches[0], branches[1])
	} else {
		FindMissing(branches[0], branches[1])
	}
}

func handleGrepBranch() {
	args := os.Args[2:]
	if len(args) == 1 {
		GrepBranch(false, args[0])
	} else if len(args) == 2 && args[0] == "--all" {
		GrepBranch(true, args[1])
	} else {
		fmt.Fprintf(os.Stderr, "Usage: %s grep-branch [--all] \"text\"\n", os.Args[0])
		os.Exit(1)
	}
}

func PrintUsage() {
	fmt.Println("Usage:")
	fmt.Println("  git-tools find-missing [--browse|-i] [--tui|-t] <branch1> <branch2>")
	fmt.Println("                         # Find commits in branch1 missing from branch2")
	fmt.Println("                         # --browse/-i: Interactive detailed view")
	fmt.Println("                         # --tui/-t: Terminal User Interface (dual-pane)")
	fmt.Println("  git-tools grep-branch [--all] \"text\"")
	fmt.Println("                         # List branches and commits where text exists in commit message")
	fmt.Println("                         # --all: search all refs (branches, remotes, tags)")
} 