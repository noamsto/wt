package main

import (
	"testing"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		wantFlags flags
		wantHelp  bool
	}{
		{
			name:      "no args",
			args:      []string{},
			wantFlags: flags{},
		},
		{
			name:      "single short flag -y",
			args:      []string{"-y"},
			wantFlags: flags{yes: true},
		},
		{
			name:      "single short flag -q",
			args:      []string{"-q"},
			wantFlags: flags{quiet: true},
		},
		{
			name:      "single short flag -i",
			args:      []string{"-i"},
			wantFlags: flags{interactive: true},
		},
		{
			name:      "combined short flags -yqi",
			args:      []string{"-yqi"},
			wantFlags: flags{yes: true, quiet: true, interactive: true},
		},
		{
			name:      "long flag --yes",
			args:      []string{"--yes"},
			wantFlags: flags{yes: true},
		},
		{
			name:      "long flag --quiet",
			args:      []string{"--quiet"},
			wantFlags: flags{quiet: true},
		},
		{
			name:      "long flag --interactive",
			args:      []string{"--interactive"},
			wantFlags: flags{interactive: true},
		},
		{
			name:      "help short flag -h",
			args:      []string{"-h"},
			wantFlags: flags{},
			wantHelp:  true,
		},
		{
			name:      "help long flag --help",
			args:      []string{"--help"},
			wantFlags: flags{},
			wantHelp:  true,
		},
		{
			name:      "help subcommand",
			args:      []string{"help"},
			wantFlags: flags{},
			wantHelp:  true,
		},
		{
			name:      "mixed short and long flags",
			args:      []string{"-y", "--quiet"},
			wantFlags: flags{yes: true, quiet: true},
		},
		{
			name:      "all flags combined",
			args:      []string{"-yqi"},
			wantFlags: flags{yes: true, quiet: true, interactive: true},
		},
		{
			name:      "in-popup flag",
			args:      []string{"-i", "--in-popup"},
			wantFlags: flags{interactive: true, inPopup: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFlags, gotHelp := parseArgs(tt.args)

			if gotFlags != tt.wantFlags {
				t.Errorf("flags = %+v, want %+v", gotFlags, tt.wantFlags)
			}

			if gotHelp != tt.wantHelp {
				t.Errorf("help = %v, want %v", gotHelp, tt.wantHelp)
			}
		})
	}
}
