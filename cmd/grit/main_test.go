package main

import "testing"

func TestParseVersionOutput(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		want   string
	}{
		{
			name:  "standard output",
			input: "grit version 3.5.0\n",
			want:  "3.5.0",
		},
		{
			name:  "version prefixed with v",
			input: "grit version v3.5.0\n",
			want:  "3.5.0",
		},
		{
			name:  "mixed case command words",
			input: "GRIT VERSION 3.6.1\n",
			want:  "3.6.1",
		},
		{
			name:  "invalid format",
			input: "version 3.5.0",
			want:  "",
		},
		{
			name:  "empty",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseVersionOutput(tt.input)
			if got != tt.want {
				t.Fatalf("parseVersionOutput(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
