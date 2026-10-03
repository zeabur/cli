package config_test

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/zeabur/cli/pkg/config"
)

func assertUnixMode(t *testing.T, path string, want os.FileMode) {
	t.Helper()
	if runtime.GOOS == "windows" {
		return
	}
	st, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := st.Mode().Perm(); got != want {
		t.Fatalf("%s mode = %o, want %o", path, got, want)
	}
}

func TestNew_ConfigFileMode600(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "zeabur", "cli.yaml")
	cfg := config.New(path)
	assertUnixMode(t, path, 0o600)

	cfg.SetTokenString("secret-token")
	if err := cfg.Write(); err != nil {
		t.Fatal(err)
	}
	assertUnixMode(t, path, 0o600)
}

func TestNew_TightensExisting0644(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cli.yaml")
	if err := os.WriteFile(path, []byte("token: old\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	_ = config.New(path)
	assertUnixMode(t, path, 0o600)
}
