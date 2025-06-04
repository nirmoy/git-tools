package gittools

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"
)

func FindMissing(branch1, branch2 string) {
	// Check if we're in a Git repository
	if !IsGitRepo() {
		fmt.Fprintf(os.Stderr, "Error: Not in a Git repository\n")
		os.Exit(1)
	}

	// Check if branches exist
	if !BranchExists(branch1) {
		fmt.Fprintf(os.Stderr, "Error: Branch '%s' does not exist\n", branch1)
		os.Exit(1)
	}

	if !BranchExists(branch2) {
		fmt.Fprintf(os.Stderr, "Error: Branch '%s' does not exist\n", branch2)
		os.Exit(1)
	}

	fmt.Printf("Finding commits in '%s' that are missing from '%s'...\n\n", branch1, branch2)

	// Get commits that are in branch1 but not in branch2 (by hash)
	missingCommits, err := getMissingCommits(branch1, branch2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting missing commits: %v\n", err)
		os.Exit(1)
	}

	// Get all commit subjects from branch2 for subject-based comparison (normalized)
	branch2Subjects, err := getAllSubjects(branch2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting subjects from branch2: %v\n", err)
		os.Exit(1)
	}

	// Filter missingCommits: only keep those whose normalized subject does NOT exist in branch2
	filteredCommits := make([]Commit, 0, len(missingCommits))
	for _, commit := range missingCommits {
		normSubj := NormalizeSubject(commit.Subject)
		if _, ok := branch2Subjects[normSubj]; !ok {
			filteredCommits = append(filteredCommits, commit)
		}
	}

	if len(filteredCommits) == 0 {
		fmt.Printf("No missing commits found. Branch '%s' is up to date with '%s'.\n", branch2, branch1)
		return
	}

	fmt.Printf("Found %d missing commit(s):\n\n", len(filteredCommits))

	// Display missing commits as one-liners with color
	for _, commit := range filteredCommits {
		fmt.Printf("%s%s%s %s (%s%s%s, %s%s%s)\n",
			ColorYellow, commit.Hash[:8], ColorReset,
			commit.Subject,
			ColorGreen, commit.Author, ColorReset,
			ColorCyan, commit.Date, ColorReset,
		)
	}

	// Sort commits by date (oldest first)
	sort.Slice(filteredCommits, func(i, j int) bool {
		return filteredCommits[i].Date < filteredCommits[j].Date
	})

	fmt.Printf("To apply these commits to branch '%s', you can (in order!):\n", branch2)
	fmt.Printf("1. Checkout '%s': git checkout %s\n", branch2, branch2)
	fmt.Printf("2. Cherry-pick commits in order (to avoid conflicts): ")
	hashes := make([]string, len(filteredCommits))
	for i, commit := range filteredCommits {
		hashes[i] = commit.Hash // full hash for safety
	}
	fmt.Printf("git cherry-pick %s\n", strings.Join(hashes, " "))
	fmt.Printf("3. Or merge '%s' into '%s': git merge %s\n", branch1, branch2, branch1)
}

// getMissingCommits finds commits in branch1 that are not in branch2 (by hash)
func getMissingCommits(branch1, branch2 string) ([]Commit, error) {
	// Use git log to find commits in branch1 but not in branch2
	// Format: hash<delim>subject<delim>author<delim>date
	cmd := exec.Command("git", "log", "--pretty=format:%H"+LogDelimiter+"%s"+LogDelimiter+"%an"+LogDelimiter+"%ad", "--date=short", branch1, "^"+branch2)

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get commit diff: %v", err)
	}

	if len(output) == 0 {
		return []Commit{}, nil
	}

	var commits []Commit
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, LogDelimiter)
		if len(parts) >= 4 {
			commit := Commit{
				Hash:    parts[0],
				Subject: parts[1],
				Author:  parts[2],
				Date:    parts[3],
			}
			commits = append(commits, commit)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading commit output: %v", err)
	}

	// Sort commits by date (oldest first)
	sort.Slice(commits, func(i, j int) bool {
		return commits[i].Date < commits[j].Date
	})

	return commits, nil
}

// getAllSubjects returns a set of all normalized commit subjects in a branch
func getAllSubjects(branch string) (map[string]struct{}, error) {
	cmd := exec.Command("git", "log", "--pretty=format:%s", branch)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get subjects: %v", err)
	}
	subjects := make(map[string]struct{})
	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		normSubj := NormalizeSubject(scanner.Text())
		subjects[normSubj] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading subject output: %v", err)
	}
	return subjects, nil
} 