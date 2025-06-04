package gittools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func GrepBranch(searchAll bool, text string) {
	if !IsGitRepo() {
		fmt.Fprintf(os.Stderr, "Error: Not in a Git repository\n")
		os.Exit(1)
	}

	var logArgs []string
	if searchAll {
		logArgs = []string{"log", "--all", "--grep", text, "--pretty=format:%H" + LogDelimiter + "%s" + LogDelimiter + "%D"}
	} else {
		logArgs = []string{"log", "--branches", "--grep", text, "--pretty=format:%H" + LogDelimiter + "%s" + LogDelimiter + "%D"}
	}
	cmd := exec.Command("git", logArgs...)
	output, err := cmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running git log: %v\n", err)
		os.Exit(1)
	}

	seen := make(map[string]struct{}) // To avoid duplicate commit/branch pairs
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, LogDelimiter, 3)
		if len(parts) < 3 {
			continue
		}
		hash := parts[0]
		subject := parts[1]
		refs := parts[2]
		for _, ref := range strings.Split(refs, ",") {
			ref = strings.TrimSpace(ref)
			if ref == "" {
				continue
			}
			if strings.HasPrefix(ref, "tag: ") {
				continue // skip tags
			}
			if idx := strings.Index(ref, "->"); idx != -1 {
				ref = strings.TrimSpace(ref[idx+2:])
			}
			key := ref + ":" + hash
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			fmt.Printf("%s%s%s %s %s%s%s\n", ColorYellow, hash[:8], ColorReset, ref, ColorGreen, subject, ColorReset)
		}
	}
} 