package gittools

import (
	"os/exec"
	"strings"
	"unicode"
)

// IsGitRepo checks if the current directory is a Git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// BranchExists checks if a branch exists (locally or remotely)
func BranchExists(branch string) bool {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	if cmd.Run() == nil {
		return true
	}
	cmd = exec.Command("git", "show-ref", "--verify", "--quiet", "refs/remotes/origin/"+branch)
	if cmd.Run() == nil {
		return true
	}
	cmd = exec.Command("git", "rev-parse", "--verify", branch)
	return cmd.Run() == nil
}

// NormalizeSubject trims, collapses whitespace, removes non-printable/control characters from a commit subject
func NormalizeSubject(subject string) string {
	subject = strings.TrimSpace(subject)
	// Remove non-printable/control characters
	clean := make([]rune, 0, len(subject))
	for _, r := range subject {
		if unicode.IsPrint(r) && !unicode.IsControl(r) {
			clean = append(clean, r)
		}
	}
	// Collapse all whitespace to single spaces
	return strings.Join(strings.Fields(string(clean)), " ")
} 