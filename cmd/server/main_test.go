package main

import (
	"bytes"
	"io"
	"testing"
)

func TestHelpAndVersion(t *testing.T) {
	for _, arg := range []string{"-h", "--help", "--version"} {
		var out bytes.Buffer
		if err := run([]string{arg}, &out); err != nil {
			t.Fatal(err)
		}
		if out.Len() == 0 {
			t.Fatal("expected output")
		}
	}
}

func TestInvalidFlags(t *testing.T) {
	for _, args := range [][]string{{"--unknown"}, {"serve"}, {"extra"}} {
		if err := run(args, io.Discard); err == nil {
			t.Fatalf("accepted %q", args)
		}
	}
}
