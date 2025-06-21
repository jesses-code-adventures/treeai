# OpenCode Trees

A powerful tool for AI-assisted development that creates isolated git worktrees for AI collaboration while maintaining clean separation from your main development environment.

## Overview

OpenCode Trees automates the creation of AI-dedicated development environments by:
- Creating isolated git worktrees in `.opencode-trees/` directories
- Spawning tmux sessions with OpenCode running in AI environments  
- Providing seamless hotkey switching between human and AI development contexts
- Maintaining git ignore patterns to keep AI artifacts separate

## Features

### Core Functionality
- **Automatic Worktree Management**: Creates and manages `.opencode-trees` directories in your git repositories
- **Branch Isolation**: Each AI session gets its own git branch and worktree
- **Tmux Integration**: Spawns dedicated tmux sessions for AI environments
- **Session Switching**: Hotkey support for rapid context switching between human and AI sessions
- **Git Integration**: Automatically handles `.gitignore` entries for AI directories

### Workflow Benefits
- **Clean Separation**: Keep AI experiments separate from main development
- **Context Preservation**: Each AI session maintains its own state and history
- **Rapid Switching**: Toggle between human and AI contexts with a single keypress
- **Branch Management**: Easily manage feature branches dedicated to AI collaboration

## Installation

```bash
# Clone the repository
git clone https://github.com/your-username/opencode-trees.git
cd opencode-trees

# Install dependencies
./install.sh

# Optional: Install as tmux plugin
# Add to ~/.tmux.conf:
# set -g @plugin 'your-username/opencode-trees'
```

## Usage

### Command Line Interface

```bash
# Create AI worktree for for branch ai-refactor
opencode-trees ai-refactor

# Create AI worktree with custom session name
opencode-trees --session custom-name branch-name
```

### TODO: Tmux Plugin Usage

## Configuration

### TODO: Tmux Configuration

## Directory Structure

```
your-project/
├── .git/
├── .gitignore              # Auto-updated with .opencode-trees
├── .opencode-trees/        # AI worktrees directory
│   ├── feature-branch-1/   # Isolated worktree
│   ├── feature-branch-2/   # Another AI session
│   └── ...
├── src/                    # Your main project
└── ...
```

## How It Works

1. **Detection**: Scans upward to find the nearest `.git` directory
2. **Worktree Creation**: Creates isolated git worktree in `.opencode-trees/[branch-name]`
3. **Git Ignore**: Automatically adds `.opencode-trees` to `.gitignore`
4. **Session Management**: Creates tmux session and launches OpenCode
5. **Context Switching**: Enables rapid switching between human and AI sessions

## Contributing

We welcome contributions from both humans and LLMs! Here's how to get involved:

### For Human Contributors
- Fork the repository and create feature branches
- Add tests for new functionality
- Update documentation as needed

### For LLM Contributors
- Use the tool itself to create AI development sessions
- Contribute improvements to AI workflow integration
- Enhance automation and user experience features
- Help optimize tmux and git integration

## License

MIT License - see [LICENSE](LICENSE) for details.
