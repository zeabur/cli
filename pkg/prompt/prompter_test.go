package prompt_test

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/zeabur/cli/pkg/prompt"
)

// withPipeStdin replaces os.Stdin with the read end of a pipe whose write end
// stays open, like a CI step or an agent shell that never sends input.
func withPipeStdin(t *testing.T) {
	t.Helper()

	r, w, err := os.Pipe()
	require.NoError(t, err)

	orig := os.Stdin
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = orig
		_ = w.Close()
		_ = r.Close()
	})
}

func requireFailsFast(t *testing.T, call func() error) {
	t.Helper()

	done := make(chan error, 1)
	go func() { done <- call() }()

	select {
	case err := <-done:
		require.True(t, errors.Is(err, prompt.ErrNonInteractive), "got %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("prompt blocked on a non-terminal stdin")
	}
}

func TestPrompterFailsWithoutTerminal(t *testing.T) {
	p := prompt.New()

	t.Run("Select", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			_, err := p.Select("Select a service", "a", []string{"a", "b"})
			return err
		})
	})

	t.Run("MultiSelect", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			_, err := p.MultiSelect("Select services", nil, []string{"a", "b"})
			return err
		})
	})

	t.Run("Input", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			_, err := p.Input("Name", "")
			return err
		})
	})

	t.Run("InputWithHelp", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			_, err := p.InputWithHelp("Name", "", "help")
			return err
		})
	})

	t.Run("Confirm", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			_, err := p.Confirm("Continue?", false)
			return err
		})
	})

	t.Run("ConfirmDeletion", func(t *testing.T) {
		withPipeStdin(t)
		requireFailsFast(t, func() error {
			return p.ConfirmDeletion("my-service")
		})
	})
}
