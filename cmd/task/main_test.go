package main

import (
	"io"
	"testing"
)

func TestTaskFlags(t *testing.T) {
	for _, tt := range []struct {
		args []string
		want options
	}{
		{nil, options{sourceDir: "./internal/task/item/data"}},
		{[]string{"--source-dir=items", "--dry-run"}, options{sourceDir: "items", dryRun: true}},
		{[]string{"-source-dir", "items", "-dry-run=false"}, options{sourceDir: "items"}},
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
	for _, args := range [][]string{{"import"}, {"--unknown"}, {"--source-dir"}, {"--dry-run=maybe"}, {"extra", "--dry-run"}} {
		if err := run(args, io.Discard); err == nil {
			t.Fatalf("accepted %q", args)
		}
	}
}
