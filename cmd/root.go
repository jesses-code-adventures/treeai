package cmd

import (
	"fmt"
	"github.com/jesses-code-adventures/treeai/treeai"
	"github.com/spf13/cobra"
	"os"
)

var mergeFlag bool
var silentFlag bool

var rootCmd = &cobra.Command{
	Use:   "treeai <worktree-name>",
	Short: "Tmux plugin for creating AI-dedicated git worktrees",
	Long: `treeai is a tmux plugin that creates isolated git worktrees in .opencode-trees/ 
directories for AI-assisted development while maintaining clean separation from your main environment.

This tool requires tmux to be installed and is designed to work as a tmux plugin.`,
	Args: cobra.ExactArgs(1),
	Run:  handleCommand,
}

func init() {
	rootCmd.Flags().BoolVar(&mergeFlag, "merge", false, "merge the worktree branch back to main and clean up")
	rootCmd.Flags().BoolVar(&silentFlag, "silent", false, "suppress all output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCommand(cmd *cobra.Command, args []string) {
	if mergeFlag {
		treeai.MergeWorktree(args[0], silentFlag)
	} else {
		treeai.CreateWorktree(args[0], silentFlag)
	}
}
