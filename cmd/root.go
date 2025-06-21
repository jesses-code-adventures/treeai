package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "opentree <worktree-name>",
	Short: "Create AI-dedicated git worktrees for OpenCode collaboration",
	Long: `OpenCode Trees creates isolated git worktrees in .opencode-trees/ directories
for AI-assisted development while maintaining clean separation from your main environment.`,
	Args: cobra.ExactArgs(1),
	Run:  createWorktree,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func createWorktree(cmd *cobra.Command, args []string) {
	worktreeName := args[0]
	
	// Find git root directory
	gitRoot, err := findGitRoot()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	
	// Create .opencode-trees directory if it doesn't exist
	opencodeTrees := filepath.Join(gitRoot, ".opencode-trees")
	if err := os.MkdirAll(opencodeTrees, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating .opencode-trees directory: %v\n", err)
		os.Exit(1)
	}
	
	// Create worktree path
	worktreePath := filepath.Join(opencodeTrees, worktreeName)
	
	// Check if worktree already exists
	if _, err := os.Stat(worktreePath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: worktree '%s' already exists\n", worktreeName)
		os.Exit(1)
	}
	
	// Create git worktree
	if err := createGitWorktree(gitRoot, worktreePath, worktreeName); err != nil {
		fmt.Fprintf(os.Stderr, "Error creating git worktree: %v\n", err)
		os.Exit(1)
	}
	
	// Update .gitignore
	if err := updateGitignore(gitRoot); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to update .gitignore: %v\n", err)
	}
	
	fmt.Printf("✅ Created worktree: %s\n", worktreePath)
	fmt.Printf("💡 To start working: cd %s\n", worktreePath)
}

func findGitRoot() (string, error) {
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

func createGitWorktree(gitRoot, worktreePath, branchName string) error {
	// Create new branch and worktree
	cmd := exec.Command("git", "worktree", "add", "-b", branchName, worktreePath)
	cmd.Dir = gitRoot
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git worktree add failed: %w\nOutput: %s", err, string(output))
	}
	
	return nil
}

func updateGitignore(gitRoot string) error {
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