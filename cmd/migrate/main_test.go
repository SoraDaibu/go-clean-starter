package main

import (
	"io"
	"testing"
)

func TestMigrationFlags(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want options
	}{
		{nil, options{direction: "up"}},
		{[]string{"--direction=down", "--interactive"}, options{direction: "down", interactive: true}},
		{[]string{"-direction", "up"}, options{direction: "up"}},
	} {
		got, err := parseFlags(tt.args, io.Discard)
		if err != nil || got != tt.want {
			t.Fatalf("parseFlags(%q) = %+v, %v; want %+v", tt.args, got, err, tt.want)
		}
	}
}

func TestHelpAndInvalidFlags(t *testing.T) {
	for _, arg := range []string{"-h", "--help"} {
		if err := run([]string{arg}, io.Discard); err != nil {
			t.Fatal(err)
		}
	}
	for _, args := range [][]string{{"up"}, {"--direction=sideways"}, {"--direction"}, {"--unknown"}, {"--interactive=maybe"}} {
		if err := run(args, io.Discard); err == nil {
			t.Fatalf("accepted %q", args)
		}
	}
}
