package cmd

import (
	"fmt"
	"os"

	"github.com/jesses-code-adventures/treeai/config"
	l "github.com/jesses-code-adventures/treeai/logger"
	"github.com/jesses-code-adventures/treeai/treeai"
	"github.com/spf13/cobra"
)

var data string
var merge bool
var silent bool
var commands []string
var bin string
var prompt string
var gitignore bool
var window bool
var debug bool

var rootCmd = &cobra.Command{
	Use:   "treeai <worktree-name>",
	Short: "Tmux plugin for creating AI-dedicated git worktrees",
	Long: `treeai is a tmux plugin that creates isolated git worktrees for AI-assisted development while maintaining clean separation from your main environment.

This tool requires tmux to be installed and is designed to work as a tmux plugin.`,
	Args: cobra.ExactArgs(1),
	Run:  handleCommand,
}

func init() {
	rootCmd.Flags().BoolVar(&merge, "merge", false, "merge the worktree branch back to main and clean up")
	rootCmd.Flags().BoolVar(&silent, "silent", false, "suppress all output")
	rootCmd.Flags().BoolVar(&gitignore, "gitignore", false, "use .gitignore instead of .git/info/exclude to exclude worktrees from git")
	rootCmd.Flags().BoolVar(&window, "window", false, "open a new tmux window with the worktree, instead of a session")
	rootCmd.Flags().BoolVar(&debug, "debug", false, "use .gitignore instead of .git/info/exclude to exclude worktrees from git")
	rootCmd.Flags().StringArrayVar(&commands, "command", []string{}, "add an additional tmux window per command, running each command in a new window")
	rootCmd.Flags().StringVar(&bin, "bin", "opencode", "binary to launch in the tmux session")
	rootCmd.Flags().StringVar(&prompt, "prompt", "", "send a prompt to opencode in the new session")
	rootCmd.Flags().StringVar(&data, "data", os.ExpandEnv("$HOME/.local/share/treeai"), "path to data directory")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleCommand(cmd *cobra.Command, args []string) {
	branchName := args[0]
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %v\n", err)
		os.Exit(1)
	}
	cfg.ApplyFlags(bin, silent, data, commands, gitignore, debug, window)
	l.Init(cfg)

	if merge && len(commands) > 0 {
		fmt.Fprintf(os.Stderr, "Error: cannot create a window when merging\n")
		os.Exit(1)
	}

	if merge && prompt != "" {
		fmt.Fprintf(os.Stderr, "Error: cannot use --prompt flag when merging\n")
		os.Exit(1)
	}

	if merge && bin != "opencode" {
		fmt.Fprintf(os.Stderr, "Error: cannot use --bin-name flag when merging\n")
		os.Exit(1)
	}

	if merge {
		treeai.MergeWorktree(cfg, branchName)
	} else {
		treeai.CreateWorktree(cfg, branchName, prompt)
	}
}
