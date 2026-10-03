package auth_test

import (
	"os"
	"testing"

	"github.com/zeabur/cli/pkg/auth"
	"golang.org/x/term"
)

func TestCanCompleteBrowserLogin_MatchesIsTerminal(t *testing.T) {
	want := term.IsTerminal(int(os.Stdin.Fd()))
	if got := auth.CanCompleteBrowserLogin(); got != want {
		t.Fatalf("CanCompleteBrowserLogin() = %v, term.IsTerminal(stdin) = %v", got, want)
	}
}
