package treeai

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jesses-code-adventures/treeai/git"
	"github.com/jesses-code-adventures/treeai/tmux"
)

func exitWithError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func printf(silent bool, format string, args ...any) {
	if !silent {
		fmt.Printf(format, args...)
	}
}

func fprintf(silent bool, w io.Writer, format string, args ...any) {
	if !silent {
		fmt.Fprintf(w, format, args...)
	}
}

func CreateWorktree(worktreeName string, silent bool) {
	if err := tmux.CheckInstalled(); err != nil {
		exitWithError("Error: %v\n", err)
	}

	gitRoot, err := git.FindRoot()
	if err != nil {
		exitWithError("Error: %v\n", err)
	}

	worktreePath, err := setupWorktreeDirectory(gitRoot, worktreeName)
	if err != nil {
		exitWithError("Error: %v\n", err)
	}

	if err := git.CreateWorktree(gitRoot, worktreePath, worktreeName); err != nil {
		exitWithError("Error creating git worktree: %v\n", err)
	}

	if err := git.UpdateIgnore(gitRoot); err != nil {
		fprintf(silent, os.Stderr, "Warning: failed to update .gitignore: %v\n", err)
	}

	sessionName, err := tmux.CreateSessionName(gitRoot, worktreeName)
	if err != nil {
		exitWithError("Error creating tmux session name: %v\n", err)
	}

	if err := tmux.CreateAndSwitchSession(sessionName, worktreePath); err != nil {
		exitWithError("Error creating tmux session: %v\n", err)
	}

	printf(silent, "Created worktree: %s\n", worktreePath)
	printf(silent, "Created tmux session: %s\n", sessionName)
}

func setupWorktreeDirectory(gitRoot, worktreeName string) (string, error) {
	opencodeTrees := filepath.Join(gitRoot, ".opencode-trees")
	if err := os.MkdirAll(opencodeTrees, 0755); err != nil {
		return "", fmt.Errorf("creating .opencode-trees directory: %w", err)
	}

	worktreePath := filepath.Join(opencodeTrees, worktreeName)
	if _, err := os.Stat(worktreePath); err == nil {
		return "", fmt.Errorf("worktree '%s' already exists", worktreeName)
	}

	return worktreePath, nil
}

func MergeWorktree(worktreeName string, silent bool) {
	gitRoot, err := git.FindRoot()
	if err != nil {
		exitWithError("Error: %v\n", err)
	}

	worktreePath := filepath.Join(gitRoot, ".opencode-trees", worktreeName)

	if err := validateMergePrerequisites(gitRoot, worktreePath, worktreeName); err != nil {
		exitWithError("Error: %v\n", err)
	}

	printf(silent, "Rebasing on main...\n")
	if err := git.RebaseOnMain(worktreePath); err != nil {
		exitWithError("Error rebasing on main: %v\n", err)
	}

	printf(silent, "Merging branch: %s\n", worktreeName)
	if err := git.MergeBranch(gitRoot, worktreeName); err != nil {
		exitWithError("Error merging branch %s: %v\n", worktreeName, err)
	}

	printf(silent, "Removing worktree: %s\n", worktreePath)
	if err := git.RemoveWorktree(gitRoot, worktreePath); err != nil {
		exitWithError("Error removing worktree: %v\n", err)
	}

	printf(silent, "Deleting branch: %s\n", worktreeName)
	if err := git.DeleteBranch(gitRoot, worktreeName); err != nil {
		exitWithError("Error deleting branch %s: %v\n", worktreeName, err)
	}

	sessionName, err := tmux.CreateSessionName(gitRoot, worktreeName)
	if err != nil {
		printf(silent, "Warning: Could not determine tmux session name: %v\n", err)
	} else {
		printf(silent, "Killing tmux session: %s\n", sessionName)
		if err := tmux.KillSession(sessionName); err != nil {
			printf(silent, "Warning: Could not kill tmux session '%s': %v\n", sessionName, err)
		}
	}

	printf(silent, "Successfully merged and cleaned up worktree: %s\n", worktreeName)
}

func validateMergePrerequisites(gitRoot, worktreePath, worktreeName string) error {
	hasChanges, err := git.HasUncommittedChanges(gitRoot)
	if err != nil {
		return fmt.Errorf("checking git status in root: %w", err)
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in git root. Please commit or stash changes first")
	}

	if _, err := os.Stat(worktreePath); os.IsNotExist(err) {
		return fmt.Errorf("worktree '%s' does not exist", worktreeName)
	}

	hasChanges, err = git.HasUncommittedChanges(worktreePath)
	if err != nil {
		return fmt.Errorf("checking git status in worktree: %w", err)
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in worktree '%s'. Please commit or stash changes first", worktreeName)
	}

	return nil
}
