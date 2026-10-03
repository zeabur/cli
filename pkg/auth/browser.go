package auth

import (
	"os"

	"golang.org/x/term"
)

// CanCompleteBrowserLogin reports whether we should start the implicit
// browser callback. /dev/null is a char device, so ModeCharDevice is the
// wrong test; term.IsTerminal is false for pipes and /dev/null.
//
// OpenURL can still return nil on this host (DISPLAY=:1, Chromium starts,
// dbus dies). Combined with a non-terminal stdin that is the hang:
// WaitForToken blocks forever.
func CanCompleteBrowserLogin() bool {
	return term.IsTerminal(int(os.Stdin.Fd()))
}