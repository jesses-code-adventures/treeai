package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %w", err)
	}

	for {
		gitDir := filepath.Join(dir, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("not in a git repository")
		}
		dir = parent
	}
}

func CreateWorktree(gitRoot, worktreePath, branchName string) error {
	// Create new branch and worktree
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func Updateignore(gitRoot string) error {
	gitignorePath := filepath.Join(gitRoot, ".gitignore")

	// Read existing .gitignore
	content := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = string(data)
	}

	// Check if .opencode-trees is already ignored
	if strings.Contains(content, ".opencode-trees") {
		return nil
	}

	// Add .opencode-trees to .gitignore
	if content != "" && !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	content += ".opencode-trees/\n"

	return os.WriteFile(gitignorePath, []byte(content), 0644)
}

func GetCurrentBranch(gitRoot string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = gitRoot

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func SwitchBranch(gitRoot, branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to switch to branch %s: %w\nOutput: %s", branchName, err, string(output))
	}

	return nil
}

func RebaseOnMain(gitRoot string) error {
	cmd := exec.Command("git", "rebase", "main")
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rebase on main: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func MergeBranch(gitRoot, branchName string) error {
	cmd := exec.Command("git", "merge", branchName)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to merge branch %s: %w\nOutput: %s", branchName, err, string(output))
	}

	return nil
}

func RemoveWorktree(gitRoot, worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to remove worktree %s: %w\nOutput: %s", worktreePath, err, string(output))
	}

	return nil
}

func DeleteBranch(gitRoot, branchName string) error {
	cmd := exec.Command("git", "branch", "-d", branchName)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to delete branch %s: %w\nOutput: %s", branchName, err, string(output))
	}

	return nil
}
