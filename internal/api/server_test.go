package api

import (
	"strings"
	"testing"

	"github.com/donation-station/donation-station/internal/cpa"
)

func TestAuthFileUnavailableMessage(t *testing.T) {
	tests := []struct {
		name      string
		file      cpa.AuthFile
		wantBlock bool
	}{
		{
			name:      "disabled credential is blocked",
			file:      cpa.AuthFile{Disabled: true},
			wantBlock: true,
		},
		{
			name:      "unavailable credential is blocked",
			file:      cpa.AuthFile{Unavailable: true, StatusMessage: "刷新失败"},
			wantBlock: true,
		},
		{
			name:      "missing project id message is blocked",
			file:      cpa.AuthFile{Status: "ok", StatusMessage: "Antigravity 凭证缺少 project_id。请重新登录或刷新凭证以发现项目。"},
			wantBlock: true,
		},
		{
			name:      "normal credential is allowed",
			file:      cpa.AuthFile{Status: "ok"},
			wantBlock: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := authFileUnavailableMessage(tt.file)
			if tt.wantBlock && strings.TrimSpace(got) == "" {
				t.Fatalf("expected credential to be blocked")
			}
			if !tt.wantBlock && got != "" {
				t.Fatalf("expected credential to be allowed, got %q", got)
			}
		})
	}
}

func TestFindAuthFileByEmail(t *testing.T) {
	files := []cpa.AuthFile{
		{Email: "user@example.com", Provider: "codex", ID: "codex-file"},
		{Email: "User@Example.com", Provider: "gemini-cli", ID: "gemini-file"},
	}

	got := findAuthFileByEmail(files, "user@example.com", "gemini_cli")
	if got == nil || got.ID != "gemini-file" {
		t.Fatalf("expected normalized provider match, got %#v", got)
	}

	got = findAuthFileByEmail(files, "user@example.com", "")
	if got == nil {
		t.Fatalf("expected email fallback match")
	}
}
