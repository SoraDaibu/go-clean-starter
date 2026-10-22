package cmd

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"io"
	"strings"
	"testing"
)

func TestParseCommand(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want commandOptions
	}{
		{"serve", []string{"serve"}, commandOptions{name: "serve"}},
		{"environment migration", []string{"migrate"}, commandOptions{name: "migrate"}},
		{"interactive up", []string{"migrate", "up"}, commandOptions{name: "migrate up", direction: "up"}},
		{"interactive down", []string{"migrate", "down"}, commandOptions{name: "migrate down", direction: "down"}},
		{"import defaults", []string{"task", "import"}, commandOptions{name: "task import", sourceDir: "./internal/task/item/data"}},
		{"import flags", []string{"task", "import", "--source-dir", "/tmp/items", "--dry-run"}, commandOptions{name: "task import", sourceDir: "/tmp/items", dryRun: true}},
		{"single dash flags", []string{"task", "import", "-source-dir=items", "-dry-run=false"}, commandOptions{name: "task import", sourceDir: "items"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseCommand(tt.args, io.Discard)
			if err != nil || got != tt.want {
				t.Fatalf("parseCommand(%q) = %+v, %v; want %+v", tt.args, got, err, tt.want)
			}
		})
	}
}

func TestInvalidCommands(t *testing.T) {
	for _, args := range [][]string{
		{"unknown"}, {"--unknown"}, {"serve", "--unknown"}, {"serve", "extra"},
		{"migrate", "sideways"}, {"migrate", "up", "extra"},
		{"task", "unknown"}, {"task", "import", "--source-dir"},
		{"task", "import", "--dry-run=maybe"}, {"task", "import", "extra", "--dry-run"},
		{"--version", "serve"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			err := Run(context.Background(), args, "test", io.Discard)
			if err == nil {
				t.Fatalf("Run(%q) accepted invalid arguments", args)
			}
		})
	}
}

func TestHelpWithoutServices(t *testing.T) {
	for _, args := range [][]string{
		nil, {"-h"}, {"--help"}, {"help"}, {"serve", "--help"},
		{"migrate", "--help"}, {"migrate", "up", "-h"}, {"migrate", "down", "--help"},
		{"task"}, {"task", "--help"}, {"task", "import", "--help"}, {"help", "task", "import"},
	} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			var output bytes.Buffer
			_, err := parseCommand(args, &output)
			if !errors.Is(err, flag.ErrHelp) {
				t.Fatalf("parseCommand(%q) = %v; want flag.ErrHelp", args, err)
			}
			if !strings.Contains(output.String(), "Usage:") {
				t.Fatalf("missing usage: %q", output.String())
			}
			if err := Run(context.Background(), args, "test", io.Discard); err != nil {
				t.Fatalf("help failed: %v", err)
			}
		})
	}
}

func TestVersion(t *testing.T) {
	for _, arg := range []string{"--version", "-version", "-v"} {
		var output bytes.Buffer
		if err := Run(context.Background(), []string{arg}, "v1.2.3 - abc123", &output); err != nil {
			t.Fatal(err)
		}
		if got := output.String(); got != "v1.2.3 - abc123\n" {
			t.Fatalf("version output = %q", got)
		}
	}
}
