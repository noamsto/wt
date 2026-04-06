package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags flags
		wantRest  []string
	}{
		{
			name:      "no args",
			args:      []string{},
			wantFlags: flags{},
			wantRest:  nil,
		},
		{
			name:      "single short flag -y",
			args:      []string{"-y"},
			wantFlags: flags{yes: true},
			wantRest:  nil,
		},
		{
			name:      "single short flag -q",
			args:      []string{"-q"},
			wantFlags: flags{quiet: true},
			wantRest:  nil,
		},
		{
			name:      "single short flag -n",
			args:      []string{"-n"},
			wantFlags: flags{noSwitch: true},
			wantRest:  nil,
		},
		{
			name:      "single short flag -i",
			args:      []string{"-i"},
			wantFlags: flags{interactive: true},
			wantRest:  nil,
		},
		{
			name:      "combined short flags -yqn",
			args:      []string{"-yqn"},
			wantFlags: flags{yes: true, quiet: true, noSwitch: true},
			wantRest:  nil,
		},
		{
			name:      "combined short flags -yqi",
			args:      []string{"-yqi"},
			wantFlags: flags{yes: true, quiet: true, interactive: true},
			wantRest:  nil,
		},
		{
			name:      "long flag --yes",
			args:      []string{"--yes"},
			wantFlags: flags{yes: true},
			wantRest:  nil,
		},
		{
			name:      "long flag --quiet",
			args:      []string{"--quiet"},
			wantFlags: flags{quiet: true},
			wantRest:  nil,
		},
		{
			name:      "long flag --no-switch",
			args:      []string{"--no-switch"},
			wantFlags: flags{noSwitch: true},
			wantRest:  nil,
		},
		{
			name:      "long flag --interactive",
			args:      []string{"--interactive"},
			wantFlags: flags{interactive: true},
			wantRest:  nil,
		},
		{
			name:      "mixed short and long flags",
			args:      []string{"-y", "--quiet", "-n"},
			wantFlags: flags{yes: true, quiet: true, noSwitch: true},
			wantRest:  nil,
		},
		{
			name:      "positional args only",
			args:      []string{"list"},
			wantFlags: flags{},
			wantRest:  []string{"list"},
		},
		{
			name:      "flags before subcommand",
			args:      []string{"-yqn", "feature-branch"},
			wantFlags: flags{yes: true, quiet: true, noSwitch: true},
			wantRest:  []string{"feature-branch"},
		},
		{
			name:      "flags around subcommand",
			args:      []string{"-y", "remove", "my-branch"},
			wantFlags: flags{yes: true},
			wantRest:  []string{"remove", "my-branch"},
		},
		{
			name:      "unknown short flag goes to rest",
			args:      []string{"-x"},
			wantFlags: flags{},
			wantRest:  []string{"-x"},
		},
		{
			name:      "combined with unknown — known flags set but arg goes to rest",
			args:      []string{"-yx"},
			wantFlags: flags{yes: true},
			wantRest:  []string{"-yx"},
		},
		{
			name:      "bare dash goes to rest",
			args:      []string{"-"},
			wantFlags: flags{},
			wantRest:  []string{"-"},
		},
		{
			name:      "unknown long flag goes to rest",
			args:      []string{"--unknown"},
			wantFlags: flags{},
			wantRest:  []string{"--unknown"},
		},
		{
			name:      "all four flags combined",
			args:      []string{"-yqni"},
			wantFlags: flags{yes: true, quiet: true, noSwitch: true, interactive: true},
			wantRest:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotRest := parseArgs(tt.args)

			if gotFlags != tt.wantFlags {
				t.Errorf("flags = %+v, want %+v", gotFlags, tt.wantFlags)
			}

			if len(gotRest) != len(tt.wantRest) {
				t.Fatalf("rest length = %d, want %d (got %v, want %v)", len(gotRest), len(tt.wantRest), gotRest, tt.wantRest)
			}
			for i := range gotRest {
				if gotRest[i] != tt.wantRest[i] {
					t.Errorf("rest[%d] = %q, want %q", i, gotRest[i], tt.wantRest[i])
				}
			}
		})
	}
}
