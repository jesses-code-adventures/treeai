# TreeAI

A cli application & tmux plugin for tight integration between tmux, git worktrees and (Opencode)[https://github.com/sst/OpenCode].

## Overview

### Create a worktree & tmux session

Call `treeai branch-name` to create a new branch & worktree called `branch-name` in the `.treeai` directory of your project. A tmux session with then be created and switched to, with windows for `nvim`, `bash` and `opencode` set up and `opencode` focused. By default, you can toggle between your main tmux session and the `opencode` session using `<prefix>L` to alternate between two tmux sessions. This allows you to assign opencode some work and quickly switch back to what you were doing.

### Merge your worktree and clean up git environment

When you're satisfied with what `opencode` has implemented, merge the worktree from your `main` directory by calling `treeai branch-name --merge`. This will check out the worktree, rebase it against `main`, merge it into `main`, prune worktrees and delete merged git branches.

## Installation

```bash
# Clone the repository
git clone https://github.com/jesses-code-adventures/treeai.git
cd treeai
make install
```

### ~/.tmux.conf

After installing `treeai`, add the following to your `~/.tmux.conf` file, then reload your configuration:

```tmux
# OpenCode Trees
bind-key o command-prompt -p "worktree name:" "run-shell 'treeai %%'" # binds creation of a new worktree to `<prefix>o`
```

## License

MIT License - see [LICENSE](LICENSE) for details.
