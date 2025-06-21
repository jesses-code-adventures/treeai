package cmd

import (
	"fmt"
	"github.com/jesses-code-adventures/opentree/opentree"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "opentree <worktree-name>",
	Short: "Tmux plugin for creating AI-dedicated git worktrees",
	Long: `OpenCode Trees is a tmux plugin that creates isolated git worktrees in .opencode-trees/ 
directories for AI-assisted development while maintaining clean separation from your main environment.

This tool requires tmux to be installed and is designed to work as a tmux plugin.`,
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
	opentree.CreateWorktree(args[0])
}
