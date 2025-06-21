package treeai

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jesses-code-adventures/treeai/git"
	"github.com/jesses-code-adventures/treeai/tmux"
)

func CreateWorktree(worktreeName string) {
	if err := tmux.CheckInstalled(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	gitRoot, err := git.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	opencodeTrees := filepath.Join(gitRoot, ".opencode-trees")
	if err2 := os.MkdirAll(opencodeTrees, 0755); err2 != nil {
		fmt.Fprintf(os.Stderr, "Error creating .opencode-trees directory: %v\n", err2)
		os.Exit(1)
	}

	worktreePath := filepath.Join(opencodeTrees, worktreeName)

	if _, err3 := os.Stat(worktreePath); err3 == nil {
		fmt.Fprintf(os.Stderr, "Error: worktree '%s' already exists\n", worktreeName)
		os.Exit(1)
	}

	if err4 := git.CreateWorktree(gitRoot, worktreePath, worktreeName); err4 != nil {
		fmt.Fprintf(os.Stderr, "Error creating git worktree: %v\n", err4)
		os.Exit(1)
	}

	if err5 := git.Updateignore(gitRoot); err5 != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to update .gitignore: %v\n", err5)
	}

	sessionName, err := tmux.CreateSessionName(gitRoot, worktreeName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tmux session name: %v\n", err)
		os.Exit(1)
	}

	if err := tmux.CreateAndSwitchSession(sessionName, worktreePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating tmux session: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created worktree: %s\n", worktreePath)
	fmt.Printf("Created tmux session: %s\n", sessionName)
}

func MergeWorktree(worktreeName string) {
	gitRoot, err := git.FindRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	opencodeTrees := filepath.Join(gitRoot, ".opencode-trees")
	worktreePath := filepath.Join(opencodeTrees, worktreeName)

	hasChanges, err := git.HasUncommittedChanges(gitRoot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking git status in root: %v\n", err)
		os.Exit(1)
	}
	if hasChanges {
		fmt.Fprintf(os.Stderr, "Error: uncommitted changes in git root. Please commit or stash changes first.\n")
		os.Exit(1)
	}

	if _, err1 := os.Stat(worktreePath); os.IsNotExist(err1) {
		fmt.Fprintf(os.Stderr, "Error: worktree '%s' does not exist\n", worktreeName)
		os.Exit(1)
	}

	hasChanges, err = git.HasUncommittedChanges(worktreePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error checking git status in worktree: %v\n", err)
		os.Exit(1)
	}
	if hasChanges {
		fmt.Fprintf(os.Stderr, "Error: uncommitted changes in worktree '%s'. Please commit or stash changes first.\n", worktreeName)
		os.Exit(1)
	}

	fmt.Printf("Rebasing on main...\n")
	if err := git.RebaseOnMain(worktreePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error rebasing on main: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Merging branch: %s\n", worktreeName)
	if err := git.MergeBranch(gitRoot, worktreeName); err != nil {
		fmt.Fprintf(os.Stderr, "Error merging branch %s: %v\n", worktreeName, err)
		os.Exit(1)
	}

	fmt.Printf("Removing worktree: %s\n", worktreePath)
	if err := git.RemoveWorktree(gitRoot, worktreePath); err != nil {
		fmt.Fprintf(os.Stderr, "Error removing worktree: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Deleting branch: %s\n", worktreeName)
	if err := git.DeleteBranch(gitRoot, worktreeName); err != nil {
		fmt.Fprintf(os.Stderr, "Error deleting branch %s: %v\n", worktreeName, err)
		os.Exit(1)
	}

	fmt.Printf("Successfully merged and cleaned up worktree: %s\n", worktreeName)
}
