package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindRoot(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() (string, func())
		wantErr bool
	}{
		{
			name: "finds git root in current directory",
			setup: func() (string, func()) {
				tmpDir, err := os.MkdirTemp("", "git-test")
				if err != nil {
					t.Fatal(err)
				}
				gitDir := filepath.Join(tmpDir, ".git")
				if err := os.Mkdir(gitDir, 0755); err != nil {
					t.Fatal(err)
				}
				oldDir, _ := os.Getwd()
				os.Chdir(tmpDir)
				return tmpDir, func() {
					os.Chdir(oldDir)
					os.RemoveAll(tmpDir)
				}
			},
			wantErr: false,
		},
		{
			name: "finds git root in parent directory",
			setup: func() (string, func()) {
				tmpDir, err := os.MkdirTemp("", "git-test")
				if err != nil {
					t.Fatal(err)
				}
				gitDir := filepath.Join(tmpDir, ".git")
				if err := os.Mkdir(gitDir, 0755); err != nil {
					t.Fatal(err)
				}
				subDir := filepath.Join(tmpDir, "subdir")
				if err := os.Mkdir(subDir, 0755); err != nil {
					t.Fatal(err)
				}
				oldDir, _ := os.Getwd()
				os.Chdir(subDir)
				return tmpDir, func() {
					os.Chdir(oldDir)
					os.RemoveAll(tmpDir)
				}
			},
			wantErr: false,
		},
		{
			name: "returns error when not in git repository",
			setup: func() (string, func()) {
				tmpDir, err := os.MkdirTemp("", "no-git-test")
				if err != nil {
					t.Fatal(err)
				}
				oldDir, _ := os.Getwd()
				os.Chdir(tmpDir)
				return "", func() {
					os.Chdir(oldDir)
					os.RemoveAll(tmpDir)
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedRoot, cleanup := tt.setup()
			defer cleanup()

			got, err := FindRoot()
			if (err != nil) != tt.wantErr {
				t.Errorf("FindRoot() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				gotAbs, _ := filepath.EvalSymlinks(got)
				expectedAbs, _ := filepath.EvalSymlinks(expectedRoot)
				if gotAbs != expectedAbs {
					t.Errorf("FindRoot() = %v, want %v", gotAbs, expectedAbs)
				}
			}
		})
	}
}

func TestUpdateIgnore(t *testing.T) {
	tests := []struct {
		name           string
		existingIgnore string
		wantContent    string
		wantErr        bool
	}{
		{
			name:           "adds to empty gitignore",
			existingIgnore: "",
			wantContent:    ".opencode-trees/\n",
			wantErr:        false,
		},
		{
			name:           "adds to existing gitignore with newline",
			existingIgnore: "*.log\n",
			wantContent:    "*.log\n.opencode-trees/\n",
			wantErr:        false,
		},
		{
			name:           "adds to existing gitignore without newline",
			existingIgnore: "*.log",
			wantContent:    "*.log\n.opencode-trees/\n",
			wantErr:        false,
		},
		{
			name:           "skips if already present",
			existingIgnore: "*.log\n.opencode-trees/\n*.tmp\n",
			wantContent:    "*.log\n.opencode-trees/\n*.tmp\n",
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir, err := os.MkdirTemp("", "gitignore-test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			gitignorePath := filepath.Join(tmpDir, ".gitignore")
			if tt.existingIgnore != "" {
				if err = os.WriteFile(gitignorePath, []byte(tt.existingIgnore), 0644); err != nil {
					t.Fatal(err)
				}
			}

			err = UpdateIgnore(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateIgnore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				content, err := os.ReadFile(gitignorePath)
				if err != nil {
					t.Fatal(err)
				}
				if string(content) != tt.wantContent {
					t.Errorf("UpdateIgnore() content = %q, want %q", string(content), tt.wantContent)
				}
			}
		})
	}
}

func TestHasUncommittedChanges(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "git-status-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	hasChanges, err := HasUncommittedChanges(tmpDir)
	if err == nil {
		t.Error("HasUncommittedChanges() should return error for non-git directory")
	}
	if hasChanges {
		t.Error("HasUncommittedChanges() should return false for non-git directory")
	}
}
