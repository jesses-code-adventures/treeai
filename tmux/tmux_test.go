package tmux

import (
	"os"
	"testing"
)

func TestCreateSessionName(t *testing.T) {
	tests := []struct {
		name            string
		gitRoot         string
		worktreeName    string
		currentSession  string
		expectedPattern string
	}{
		{
			name:            "creates session name from git root",
			gitRoot:         "/home/user/myproject",
			worktreeName:    "feature-branch",
			currentSession:  "",
			expectedPattern: "myproject-feature-branch",
		},
		{
			name:            "creates session name from current session",
			gitRoot:         "/home/user/myproject",
			worktreeName:    "feature-branch",
			currentSession:  "existing-session",
			expectedPattern: "existing-session-feature-branch",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTmux := os.Getenv("TMUX")
			defer os.Setenv("TMUX", oldTmux)

			if tt.currentSession != "" {
				os.Setenv("TMUX", "dummy-value")
			} else {
				os.Unsetenv("TMUX")
			}

			got, err := SessionName(tt.gitRoot, tt.worktreeName)
			if err != nil && tt.currentSession == "" {
				t.Errorf("CreateSessionName() error = %v", err)
				return
			}

			if tt.currentSession == "" && got != tt.expectedPattern {
				t.Errorf("CreateSessionName() = %v, want %v", got, tt.expectedPattern)
			}
		})
	}
}

func TestGetCurrentSession(t *testing.T) {
	tests := []struct {
		name       string
		tmuxEnv    string
		wantResult string
		wantErr    bool
	}{
		{
			name:       "returns empty when not in tmux",
			tmuxEnv:    "",
			wantResult: "",
			wantErr:    false,
		},
		{
			name:       "detects tmux session",
			tmuxEnv:    "dummy-value",
			wantResult: "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			oldTmux := os.Getenv("TMUX")
			defer os.Setenv("TMUX", oldTmux)

			if tt.tmuxEnv == "" {
				os.Unsetenv("TMUX")
			} else {
				os.Setenv("TMUX", tt.tmuxEnv)
			}

			got, err := GetCurrentSession()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCurrentSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.wantResult {
				t.Errorf("GetCurrentSession() = %v, want %v", got, tt.wantResult)
			}
		})
	}
}

func TestCheckInstalled(t *testing.T) {
	err := CheckInstalled()
	if err != nil {
		t.Logf("tmux not installed: %v", err)
	} else {
		t.Log("tmux is installed")
	}
}
