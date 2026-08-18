package main

import "testing"

func TestParseCommand(t *testing.T) {
	tests := []struct {
		arg  string
		want command
	}{
		{arg: "init", want: cmdInit},
		{arg: "list", want: cmdList},
		{arg: "completion", want: cmdCompletion},
		{arg: "prod-app", want: cmdConnect},
	}

	for _, tt := range tests {
		t.Run(tt.arg, func(t *testing.T) {
			if got := parseCommand(tt.arg); got != tt.want {
				t.Fatalf("parseCommand(%q) = %v, want %v", tt.arg, got, tt.want)
			}
		})
	}
}
