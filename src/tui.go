package gittools

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/jroimartin/gocui"
)

type TUI struct {
	gui     *gocui.Gui
	commits []Commit
	current int
	branch1 string
	branch2 string
}

func FindMissingTUI(branch1, branch2 string) {
	// Check if we're in a Git repository
	if !IsGitRepo() {
		fmt.Printf("Error: Not in a Git repository\n")
		return
	}

	// Check if branches exist
	if !BranchExists(branch1) {
		fmt.Printf("Error: Branch '%s' does not exist\n", branch1)
		return
	}

	if !BranchExists(branch2) {
		fmt.Printf("Error: Branch '%s' does not exist\n", branch2)
		return
	}

	// Get commits that are in branch1 but not in branch2 (by hash)
	missingCommits, err := getMissingCommits(branch1, branch2)
	if err != nil {
		fmt.Printf("Error getting missing commits: %v\n", err)
		return
	}

	// Get all commit subjects from branch2 for subject-based comparison (normalized)
	branch2Subjects, err := getAllSubjects(branch2)
	if err != nil {
		fmt.Printf("Error getting subjects from branch2: %v\n", err)
		return
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

	// Start TUI
	startTUI(filteredCommits, branch1, branch2)
}

func startTUI(commits []Commit, branch1, branch2 string) {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	// Enable colors and mouse support
	g.Cursor = true
	g.Mouse = true

	tui := &TUI{
		gui:     g,
		commits: commits,
		current: 0,
		branch1: branch1,
		branch2: branch2,
	}

	g.SetManagerFunc(tui.layout)
	
	if err := tui.setKeybindings(); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func (t *TUI) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	
	// Header view - make it taller and enable wrapping for help text
	headerHeight := 3
	if maxX < 100 { // For narrow terminals, make header even taller
		headerHeight = 4
	}
	if v, err := g.SetView("header", 0, 0, maxX-1, headerHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
		v.Wrap = true  // Enable text wrapping
		fmt.Fprintf(v, "Git Tools - Missing Commits: '%s' -> '%s' (%d commits)\n", 
			t.branch1, t.branch2, len(t.commits))
		fmt.Fprintf(v, "Controls: q:quit ↑↓/jk:navigate ←→/hl:scroll Enter:focus PgUp/PgDn/Space:patch")
	}

	// Commit list view (left side) - adjust for new header height
	listWidth := maxX / 2
	if v, err := g.SetView("list", 0, headerHeight+1, listWidth, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Commits"
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		v.Wrap = false  // Disable wrapping for horizontal scrolling
		t.updateCommitList(v)
		t.setCursor(v, t.current)
		
		if _, err := g.SetCurrentView("list"); err != nil {
			return err
		}
	}

	// Detail view (right side) - adjust for new header height
	if v, err := g.SetView("detail", listWidth+1, headerHeight+1, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "Commit Details (git show)"
		v.Wrap = false
		v.Autoscroll = false
		v.Frame = true
		if len(t.commits) > 0 {
			t.updateCommitDetail(v, t.current)
		}
	}

	return nil
}

func (t *TUI) updateCommitList(v *gocui.View) {
	v.Clear()
	for _, commit := range t.commits {
		fmt.Fprintf(v, "%s %s (%s, %s)\n", commit.Hash[:8], 
			commit.Subject, commit.Author, commit.Date)
	}
}

func (t *TUI) setCursor(v *gocui.View, line int) {
	if err := v.SetCursor(0, line); err != nil {
		// If we can't set cursor, try setting origin
		v.SetOrigin(0, 0)
		v.SetCursor(0, line)
	}
}

func (t *TUI) updateCommitDetail(v *gocui.View, index int) {
	v.Clear()
	v.SetOrigin(0, 0) // Reset scroll position when switching commits
	if index >= len(t.commits) {
		return
	}
	
	commit := t.commits[index]
	
	// Get full commit with patch (like git log -p) with color
	fullPatch, err := getCommitFullPatch(commit.Hash)
	if err != nil {
		fmt.Fprintf(v, "Error getting commit details: %v", err)
		return
	}
	
	// Display the full colored patch
	fmt.Fprint(v, fullPatch)
	
	// Add cherry-pick instruction at the end
	fmt.Fprintf(v, "\n%s--- Cherry-pick command ---%s\n", "\033[1;36m", "\033[0m")
	fmt.Fprintf(v, "%sgit cherry-pick %s%s\n", "\033[1;32m", commit.Hash, "\033[0m")
}

func (t *TUI) setKeybindings() error {
	if err := t.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}
	
	// Global Page Up/Down for patch view scrolling (works from any pane)
	if err := t.gui.SetKeybinding("", gocui.KeyPgup, gocui.ModNone, t.pageUpPatch); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("", gocui.KeyPgdn, gocui.ModNone, t.pageDownPatch); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("", gocui.KeySpace, gocui.ModNone, t.pageDownPatch); err != nil {
		return err
	}
	
	// List navigation
	if err := t.gui.SetKeybinding("list", gocui.KeyArrowUp, gocui.ModNone, t.cursorUp); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", gocui.KeyArrowDown, gocui.ModNone, t.cursorDown); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", 'k', gocui.ModNone, t.cursorUp); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", 'j', gocui.ModNone, t.cursorDown); err != nil {
		return err
	}
	// Horizontal scrolling for long commit messages
	if err := t.gui.SetKeybinding("list", gocui.KeyArrowLeft, gocui.ModNone, t.scrollListLeft); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", gocui.KeyArrowRight, gocui.ModNone, t.scrollListRight); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", 'h', gocui.ModNone, t.scrollListLeft); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", 'l', gocui.ModNone, t.scrollListRight); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("list", gocui.KeyEnter, gocui.ModNone, t.showCommit); err != nil {
		return err
	}
	
	// Detail view navigation (line by line)
	if err := t.gui.SetKeybinding("detail", gocui.KeyArrowUp, gocui.ModNone, t.scrollDetailUp); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("detail", gocui.KeyArrowDown, gocui.ModNone, t.scrollDetailDown); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("detail", 'k', gocui.ModNone, t.scrollDetailUp); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("detail", 'j', gocui.ModNone, t.scrollDetailDown); err != nil {
		return err
	}
	if err := t.gui.SetKeybinding("detail", gocui.KeyEsc, gocui.ModNone, t.backToList); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (t *TUI) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if t.current > 0 {
		t.current--
		t.setCursor(v, t.current)
		if detailView, err := g.View("detail"); err == nil {
			t.updateCommitDetail(detailView, t.current)
		}
	}
	return nil
}

func (t *TUI) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if t.current < len(t.commits)-1 {
		t.current++
		t.setCursor(v, t.current)
		if detailView, err := g.View("detail"); err == nil {
			t.updateCommitDetail(detailView, t.current)
		}
	}
	return nil
}

func (t *TUI) showCommit(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("detail"); err != nil {
		return err
	}
	return nil
}

func (t *TUI) backToList(g *gocui.Gui, v *gocui.View) error {
	if _, err := g.SetCurrentView("list"); err != nil {
		return err
	}
	return nil
}

func (t *TUI) scrollDetailUp(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if oy > 0 {
		if err := v.SetOrigin(ox, oy-1); err != nil {
			return err
		}
	}
	return nil
}

func (t *TUI) scrollDetailDown(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if err := v.SetOrigin(ox, oy+1); err != nil {
		return err
	}
	return nil
}

func (t *TUI) pageUpPatch(g *gocui.Gui, v *gocui.View) error {
	// Always scroll the detail view, regardless of current focus
	detailView, err := g.View("detail")
	if err != nil {
		return err
	}
	
	_, maxY := detailView.Size()
	ox, oy := detailView.Origin()
	newY := oy - maxY + 3  // Leave some overlap
	if newY < 0 {
		newY = 0
	}
	if err := detailView.SetOrigin(ox, newY); err != nil {
		return err
	}
	return nil
}

func (t *TUI) pageDownPatch(g *gocui.Gui, v *gocui.View) error {
	// Always scroll the detail view, regardless of current focus
	detailView, err := g.View("detail")
	if err != nil {
		return err
	}
	
	_, maxY := detailView.Size()
	ox, oy := detailView.Origin()
	newY := oy + maxY - 3  // Leave some overlap
	if err := detailView.SetOrigin(ox, newY); err != nil {
		return err
	}
	return nil
}

func getCommitFullPatch(hash string) (string, error) {
	// Use git show with color to get full patch information
	cmd := exec.Command("git", "show", "--color=always", "--stat", "--patch", hash)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func getCommitDiffStats(hash string) (string, error) {
	cmd := exec.Command("git", "show", "--stat", "--format=", hash)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func (t *TUI) scrollListLeft(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if ox > 0 {
		if err := v.SetOrigin(ox-3, oy); err != nil {
			return err
		}
	}
	return nil
}

func (t *TUI) scrollListRight(g *gocui.Gui, v *gocui.View) error {
	ox, oy := v.Origin()
	if err := v.SetOrigin(ox+3, oy); err != nil {
		return err
	}
	return nil
} 