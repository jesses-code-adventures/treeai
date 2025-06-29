package treeai

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/jesses-code-adventures/treeai/config"
	"github.com/jesses-code-adventures/treeai/git"
	"github.com/jesses-code-adventures/treeai/logger"
	"github.com/jesses-code-adventures/treeai/tmux"
)

func exitWithError(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}

func CopyFile(srcPath, dstPath string) error {
	if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	srcFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("error opening source file: %v", err)
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("error creating destination file: %v", err)
	}
	defer dstFile.Close()

	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}
	return nil
}

func CreateWorktree(cfg *config.Config, worktreeName, prompt string) {
	if cfg == nil {
		cfg = config.New()
	}
	logger.Init(cfg)
	l := logger.Logger
	if err := tmux.CheckInstalled(); err != nil {
		exitWithError("Error: %v\n", err)
	}

	gitRoot, err := git.FindRoot()
	if err != nil {
		exitWithError("Error: %v\n", err)
	}
	l.Debug(fmt.Sprintf("gitRoot: %s", gitRoot))

	worktreePath, err := setupWorktreeDirectory(cfg, worktreeName)
	if err != nil {
		exitWithError("Error: %v\n", err)
	}
	l.Debug(fmt.Sprintf("worktreePath: %s", worktreePath))

	if err = git.CreateWorktree(gitRoot, worktreePath, worktreeName); err != nil {
		exitWithError("Error creating git worktree: %v\n", err)
	}

	if len(cfg.Copy) > 0 {
		for _, file := range cfg.Copy {
			srcPath := filepath.Join(gitRoot, file)
			dstPath := filepath.Join(worktreePath, file)
			if err = CopyFile(srcPath, dstPath); err != nil {
				exitWithError("Error copying file: %v\n", err)
			}
		}
	}

	// TODO: might not need this if using data dir
	if err = git.UpdateIgnore(gitRoot, cfg.Gitignore); err != nil {
		l.Warn(fmt.Sprintf("Warning: failed to update .gitignore: %v\n", err))
	}

	if cfg.Bin == "opencode" {
		cfg.Bin = cfg.Bin + " " + worktreePath
	}

	if cfg.Window {
		s, err := tmux.CreateAndSwitchToWindow(cfg, worktreeName, prompt)
		if err != nil {
			exitWithError("Error creating tmux window: %v\n", err)
		}
		l.Info(fmt.Sprintf("Created tmux window: %s\n", s))
	} else {
		s, err := tmux.CreateAndSwitchSession(cfg, worktreeName, prompt)
		if err != nil {
			exitWithError("Error creating tmux session: %v\n", err)
		}
		l.Info(fmt.Sprintf("Created tmux session: %s\n", s))
	}

	l.Info(fmt.Sprintf("Created worktree: %s\n", worktreePath))
}

func setupWorktreeDirectory(cfg *config.Config, worktreeName string) (string, error) {
	dataDir := filepath.Join(cfg.Data)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return "", fmt.Errorf("error creating data directory: %w", err)
	}

	worktreePath := filepath.Join(dataDir, worktreeName)
	if _, err := os.Stat(worktreePath); err == nil {
		return "", fmt.Errorf("worktree '%s' already exists", worktreeName)
	}

	return worktreePath, nil
}

func MergeWorktree(cfg *config.Config, worktreeName string) {
	if cfg == nil {
		cfg = config.New()
	}
	logger.Init(cfg)
	l := logger.Logger

	gitRoot, err := git.FindRoot()
	if err != nil {
		exitWithError("Error: %v\n", err)
	}

	worktreePath := filepath.Join(cfg.Data, worktreeName)

	if err = validateMergePrerequisites(gitRoot, worktreePath, worktreeName); err != nil {
		exitWithError("Error: %v\n", err)
	}

	currentBranch, err := git.GetCurrentBranch(gitRoot)
	if err != nil {
		exitWithError("Error: %v\n", err)
	}

	l.Info(fmt.Sprintf("Rebasing on %s...\n", currentBranch))
	if err = git.RebaseOnBranch(worktreePath, currentBranch); err != nil {
		exitWithError("Error rebasing on %s: %v\n", currentBranch, err)
	}

	l.Info(fmt.Sprintf("Merging branch: %s\n", worktreeName))
	if err = git.MergeBranch(gitRoot, worktreeName); err != nil {
		exitWithError("Error merging branch %s: %v\n", worktreeName, err)
	}

	l.Info(fmt.Sprintf("Removing worktree: %s\n", worktreePath))
	if err = git.RemoveWorktree(gitRoot, worktreePath); err != nil {
		exitWithError("Error removing worktree: %v\n", err)
	}

	l.Info(fmt.Sprintf("Deleting branch: %s\n", worktreeName))
	if err = git.DeleteBranch(gitRoot, worktreeName); err != nil {
		exitWithError("Error deleting branch %s: %v\n", worktreeName, err)
	}

	sessionName, err := tmux.SessionName(gitRoot, worktreeName)
	if err != nil {
		l.Error(fmt.Sprintf("Warning: Could not determine tmux session name: %v\n", err))
		return
	} else {
		l.Info(fmt.Sprintf("Killing tmux session: %s\n", sessionName))
		if err = tmux.KillSession(sessionName); err != nil {
			l.Error("Warning: Could not kill tmux session '%s': %v\n", sessionName, err)
			return
		}
	}

	l.Info(fmt.Sprintf("Successfully merged and cleaned up worktree: %s\n", worktreeName))
}

func validateMergePrerequisites(gitRoot, worktreePath, worktreeName string) error {
	hasChanges, err := git.HasUncommittedChanges(gitRoot)
	if err != nil {
		return fmt.Errorf("checking git status in root: %w", err)
	}
	if hasChanges {
		return fmt.Errorf("uncommitted changes in git root. Please commit or stash changes first")
	}

	if _, err = os.Stat(worktreePath); os.IsNotExist(err) {
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
