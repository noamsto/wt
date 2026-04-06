package prompt

import "testing"

func TestBestMatch(t *testing.T) {
	tests := []struct {
		name  string
		items []string
		query string
		want  string
	}{
		{
			name:  "empty query returns first item",
			items: []string{"/repo/.worktrees/feat-a", "/repo/.worktrees/feat-b"},
			query: "",
			want:  "/repo/.worktrees/feat-a",
		},
		{
			name:  "empty query empty items returns empty",
			items: []string{},
			query: "",
			want:  "",
		},
		{
			name:  "exact match",
			items: []string{"/repo/.worktrees/feat-a", "/repo/.worktrees/feat-b"},
			query: "/repo/.worktrees/feat-b",
			want:  "/repo/.worktrees/feat-b",
		},
		{
			name:  "substring match",
			items: []string{"/repo/.worktrees/feat-a", "/repo/.worktrees/fix-bug"},
			query: "fix",
			want:  "/repo/.worktrees/fix-bug",
		},
		{
			name:  "case insensitive match",
			items: []string{"/repo/.worktrees/Feature-A", "/repo/.worktrees/fix-bug"},
			query: "feature",
			want:  "/repo/.worktrees/Feature-A",
		},
		{
			name:  "no match returns first item",
			items: []string{"/repo/.worktrees/feat-a", "/repo/.worktrees/feat-b"},
			query: "nonexistent",
			want:  "/repo/.worktrees/feat-a",
		},
		{
			name:  "no match empty items returns empty",
			items: []string{},
			query: "anything",
			want:  "",
		},
		{
			name:  "multiple matches returns first",
			items: []string{"/repo/.worktrees/feat-a", "/repo/.worktrees/feat-b"},
			query: "feat",
			want:  "/repo/.worktrees/feat-a",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := bestMatch(tt.items, tt.query)
			if got != tt.want {
				t.Errorf("bestMatch() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConfirm_SkipPrompt(t *testing.T) {
	if !Confirm("test?", true) {
		t.Error("Confirm with skipPrompt=true should return true")
	}
}
