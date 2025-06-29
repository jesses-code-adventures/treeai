package tmux

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/jesses-code-adventures/treeai/config"
	"github.com/jesses-code-adventures/treeai/git"
)

func CheckInstalled() error {
	_, err := exec.LookPath("tmux")
	if err != nil {
		return fmt.Errorf(`tmux is not installed or not in PATH

treeai requires tmux to be installed

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
func SessionName(gitRoot, worktreeName string) (string, error) {
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

func CreateAndSwitchSession(cfg *config.Config, worktreeName, prompt string) (string, error) {
	gitRoot, err := git.FindRoot()
	if err != nil {
		return "", err
	}

	if cfg == nil {
		cfg = config.New()
	}

	sessionName, err := SessionName(gitRoot, worktreeName)
	if err != nil {
		return "", err
	}

	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	if checkCmd.Run() == nil {
		return sessionName, fmt.Errorf("tmux session '%s' already exists", sessionName)
	}

	createCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName, "-c", cfg.WorktreePath(worktreeName))
	if err = createCmd.Run(); err != nil {
		return sessionName, fmt.Errorf("failed to create tmux session: %w", err)
	}

	// Send the binary command to the shell
	sendCmd := exec.Command("tmux", "send-keys", "-t", sessionName+":0", cfg.Bin, "Enter")
	if err = sendCmd.Run(); err != nil {
		return sessionName, fmt.Errorf("failed to send binary command: %w", err)
	}

	// Create additional windows with custom commands
	for _, command := range cfg.Commands {
		windowCmd := exec.Command("tmux", "new-window", "-t", sessionName, "-c", cfg.WorktreePath(worktreeName), "bash", "-c", command)
		if err = windowCmd.Run(); err != nil {
			return sessionName, fmt.Errorf("failed to create window with command '%s': %w", command, err)
		}
	}

	// Always select window 0 (the specified binary) as the default focused window
	selectCmd := exec.Command("tmux", "select-window", "-t", sessionName+":0")
	if err = selectCmd.Run(); err != nil {
		return sessionName, fmt.Errorf("failed to select binary window: %w", err)
	}

	// If a prompt is provided, send it to opencode and don't switch/attach to the session
	if prompt != "" {
		sendCmd := exec.Command("tmux", "send-keys", "-t", sessionName+":0", prompt, "Enter")
		if err = sendCmd.Run(); err != nil {
			return sessionName, fmt.Errorf("failed to send prompt to opencode: %w", err)
		}
		return sessionName, nil
	}

	currentSession, err := GetCurrentSession()
	if err != nil {
		return sessionName, err
	}

	if currentSession != "" {
		// We're inside tmux, switch to the new session
		switchCmd := exec.Command("tmux", "switch-client", "-t", sessionName)
		if err := switchCmd.Run(); err != nil {
			return sessionName, fmt.Errorf("failed to switch to tmux session: %w", err)
		}
	} else {
		// We're outside tmux, attach to the new session
		attachCmd := exec.Command("tmux", "attach-session", "-t", sessionName)
		attachCmd.Stdin = os.Stdin
		attachCmd.Stdout = os.Stdout
		attachCmd.Stderr = os.Stderr
		if err := attachCmd.Run(); err != nil {
			return sessionName, fmt.Errorf("failed to attach to tmux session: %w", err)
		}
	}

	return sessionName, nil
}

func CreateAndSwitchToWindow(cfg *config.Config, worktreeName, prompt string) (string, error) {
	if cfg == nil {
		cfg = config.New()
	}

	windowName := worktreeName
	createCmd := exec.Command("tmux", "new-window", "-n", windowName, "-c", cfg.WorktreePath(worktreeName))
	if err := createCmd.Run(); err != nil {
		return windowName, fmt.Errorf("failed to create tmux window: %w", err)
	}

	// Send the binary command to the shell in the new window
	sendCmd := exec.Command("tmux", "send-keys", "-t", windowName, cfg.Bin, "Enter")
	if err := sendCmd.Run(); err != nil {
		return windowName, fmt.Errorf("failed to send binary command: %w", err)
	}

	// If a prompt is provided, send it to opencode
	if prompt != "" {
		sendCmd := exec.Command("tmux", "send-keys", "-t", windowName, prompt, "Enter")
		if err := sendCmd.Run(); err != nil {
			return windowName, fmt.Errorf("failed to send prompt to opencode: %w", err)
		}
	}

	return windowName, nil
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

func KillSession(sessionName string) error {
	checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
	if checkCmd.Run() != nil {
		return nil // Session doesn't exist, nothing to kill
	}

	killCmd := exec.Command("tmux", "kill-session", "-t", sessionName)
	if err := killCmd.Run(); err != nil {
		return fmt.Errorf("failed to kill tmux session '%s': %w", sessionName, err)
	}

	return nil
}
