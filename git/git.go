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
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
	cmd.Dir = gitRoot

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func UpdateIgnore(gitRoot string) error {
	gitignorePath := filepath.Join(gitRoot, ".gitignore")

	content := ""
	if data, err := os.ReadFile(gitignorePath); err == nil {
		content = string(data)
	}

	if strings.Contains(content, ".opencode-trees") {
		return nil
	}

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

func RebaseOnMain(workingDir string) error {
	if hasConflicts, err := checkRebaseConflicts(workingDir); err != nil {
		return fmt.Errorf("failed to check for rebase conflicts: %w", err)
	} else if hasConflicts {
		return fmt.Errorf("rebase conflicts detected. resolve conflicts manually first")
	}

	cmd := exec.Command("git", "rebase", "main")
	cmd.Dir = workingDir

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to rebase on main: %w\nOutput: %s", err, string(output))
	}

	return nil
}

func checkRebaseConflicts(workingDir string) (bool, error) {
	currentBranch, err := GetCurrentBranch(workingDir)
	if err != nil {
		return false, fmt.Errorf("getting current branch: %w", err)
	}

	cmd := exec.Command("git", "merge-tree", "main", currentBranch)
	cmd.Dir = workingDir

	output, err := cmd.Output()
	if err != nil {
		return true, nil
	}

	conflictMarkers := []string{"<<<<<<<", "=======", ">>>>>>>"}
	outputStr := string(output)
	for _, marker := range conflictMarkers {
		if strings.Contains(outputStr, marker) {
			return true, nil
		}
	}

	return false, nil
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

func HasUncommittedChanges(dir string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("failed to check git status: %w", err)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}
