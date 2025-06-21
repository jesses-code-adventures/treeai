package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func CheckInstalled() error {
	_, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf(`tmux is not installed or not in PATH

TreeAI requires tmux to be installed

To install tmux:
  • macOS: brew install tmux
  • Ubuntu/Debian: sudo apt install tmux  
  • CentOS/RHEL: sudo yum install tmux
  • Arch Linux: sudo pacman -S tmux

After installing tmux, you can use this tool to integrate git worktrees with opencode & tmux sessions`)
	}
	return nil
}

func GetCurrentSession() (string, error) {
	tmuxSession := os.Getenv("TMUX")
	if tmuxSession == "" {
		return "", nil // Not in a tmux session
	}

	cmd := exec.Command("tmux", "display-message", "-p", "#S")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current tmux session: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// CreateSessionName creates a tmux session name based on the git root directory and the target worktree name
func CreateSessionName(gitRoot, worktreeName string) (string, error) {
	currentSession, err := GetCurrentSession()
	if err != nil {
		return "", err
	}

	var baseSessionName string
	if currentSession != "" {
		baseSessionName = currentSession
	} else {
		baseSessionName = filepath.Base(gitRoot)
	}

	return fmt.Sprintf("%s-%s", baseSessionName, worktreeName), nil
}

func CreateAndSwitchSession(sessionName, worktreePath string) error {
	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	if checkCmd.Run() == nil {
		return fmt.Errorf("tmux session '%s' already exists", sessionName)
	}

	createCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", worktreePath, "nvim")
	if err := createCmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux session: %w", err)
	}

	window2Cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-c", worktreePath)
	if err := window2Cmd.Run(); err != nil {
		return fmt.Errorf("failed to create second window: %w", err)
	}

	window3Cmd := exec.Command("tmux", "new-window", "-t", sessionName, "-c", worktreePath, "opencode")
	if err := window3Cmd.Run(); err != nil {
		return fmt.Errorf("failed to create third window with opencode: %w", err)
	}

	selectCmd := exec.Command("tmux", "select-window", "-t", sessionName+":2")
	if err := selectCmd.Run(); err != nil {
		return fmt.Errorf("failed to select first window: %w", err)
	}

	currentSession, err := GetCurrentSession()
	if err != nil {
		return err
	}

	if currentSession != "" {
		// We're inside tmux, switch to the new session
		switchCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchCmd.Run(); err != nil {
			return fmt.Errorf("failed to switch to tmux session: %w", err)
		}
	} else {
		// We're outside tmux, attach to the new session
		attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
		attachCmd.Stdin = os.Stdin
		attachCmd.Stdout = os.Stdout
		attachCmd.Stderr = os.Stderr
		if err := attachCmd.Run(); err != nil {
			return fmt.Errorf("failed to attach to tmux session: %w", err)
		}
	}

	return nil
}

func SwitchToSession(sessionName string) error {
	currentSession, err := GetCurrentSession()
	if err != nil {
		return err
	}

	if currentSession == "" {
		return fmt.Errorf("not currently in a tmux session")
	}

	if currentSession == sessionName {
		return nil
	}

	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	if checkCmd.Run() != nil {
		return fmt.Errorf("tmux session '%s' does not exist", sessionName)
	}

	switchCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
	if err := switchCmd.Run(); err != nil {
		return fmt.Errorf("failed to switch to tmux session '%s': %w", sessionName, err)
	}

	return nil
}
