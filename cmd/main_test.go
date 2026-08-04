package main

import (
	"os"
	"testing"
)

func TestRun_UnknownSubcommandReturnsNonZero(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"zeabur", "this-subcommand-does-not-exist-zeabur-cli-test"}

	if code := run(); code == 0 {
		t.Fatalf("run() = %d, want non-zero when Cobra returns an error", code)
	}
}

func TestRun_VersionReturnsZero(t *testing.T) {
	oldArgs := os.Args
	t.Cleanup(func() { os.Args = oldArgs })

	os.Args = []string{"zeabur", "version"}

	if code := run(); code != 0 {
		t.Fatalf("run() = %d, want 0 for zeabur version", code)
	}
}
