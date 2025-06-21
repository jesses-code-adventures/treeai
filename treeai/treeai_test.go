package treeai

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSetupWorktreeDirectory(t *testing.T) {
	tests := []struct {
		name         string
		worktreeName string
		setupFunc    func(string) error
		wantErr      bool
		errContains  string
	}{
		{
			name:         "creates new worktree directory successfully",
			worktreeName: "test-branch",
			setupFunc:    nil,
			wantErr:      false,
		},
		{
			name:         "returns error when worktree already exists",
			worktreeName: "existing-branch",
			setupFunc: func(gitRoot string) error {
				opencodeTrees := filepath.Join(gitRoot, ".opencode-trees")
				if err := os.MkdirAll(opencodeTrees, 0755); err != nil {
					return err
				}
				existingWorktree := filepath.Join(opencodeTrees, "existing-branch")
				return os.Mkdir(existingWorktree, 0755)
			},
			wantErr:     true,
			errContains: "already exists",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "worktree-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			if tt.setupFunc != nil {
				if err := tt.setupFunc(tmpDir); err != nil {
					t.Fatal(err)
				}
			}

			got, err := setupWorktreeDirectory(tmpDir, tt.worktreeName)
			if (err != nil) != tt.wantErr {
				t.Errorf("setupWorktreeDirectory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && tt.errContains != "" {
				if err == nil || err.Error() != "worktree '"+tt.worktreeName+"' already exists" {
					t.Errorf("setupWorktreeDirectory() error = %v, want error containing %v", err, tt.errContains)
				}
			}

			if !tt.wantErr {
				expectedPath := filepath.Join(tmpDir, ".opencode-trees", tt.worktreeName)
				if got != expectedPath {
					t.Errorf("setupWorktreeDirectory() = %v, want %v", got, expectedPath)
				}

				if _, err := os.Stat(filepath.Join(tmpDir, ".opencode-trees")); os.IsNotExist(err) {
					t.Error("setupWorktreeDirectory() should create .opencode-trees directory")
				}
			}
		})
	}
}

func TestValidateMergePrerequisites(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "merge-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	worktreePath := filepath.Join(tmpDir, ".opencode-trees", "test-branch")
	if err := os.MkdirAll(worktreePath, 0755); err != nil {
		t.Fatal(err)
	}

	err = validateMergePrerequisites(tmpDir, worktreePath, "test-branch")
	if err == nil {
		t.Error("validateMergePrerequisites() should return error for non-git directory")
	}
}
