package cmdutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zeabur/cli/pkg/model"
)

// ResolveWorkspaceArg turns user-supplied "<name|id>" into a concrete team
// from the caller's memberships. The same resolution rules cover both
// `zeabur workspace switch <arg>` and the global `--workspace <arg>` flag:
//
//   - 24-char hex (case-insensitive) → matched against team IDs in the
//     memberships slice. **Membership is enforced** — an ID not in
//     `teams` is rejected here, before any further backend call. Backend
//     RBAC remains the source of truth for every subsequent operation, so
//     this is a UX shortcut rather than a security gate.
//   - non-hex → matched by team name against the membership list.
//   - exactly one match → return that team.
//   - zero matches → "no workspace named ..." error.
//   - two or more matches → "ambiguous" error listing each candidate with
//     the concrete `<bin> workspace switch <id>` invocation that
//     disambiguates it. Team names are unconstrained; users must use the
//     ID to pick.
//
// The caller passes the membership list in (typically Factory.ListTeams) so
// the per-process cache is shared with the lazy startup verify and any
// downstream commands that surface roles.
func ResolveWorkspaceArg(teams []model.Team, arg string) (*model.Team, error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return nil, errors.New("workspace name or id is required")
	}

	if isObjectIDHex(arg) {
		target := strings.ToLower(arg)
		for i := range teams {
			if strings.ToLower(teams[i].ID) == target {
				return &teams[i], nil
			}
		}
		return nil, fmt.Errorf("no team with id %q in your memberships", arg)
	}

	var matches []model.Team
	for i := range teams {
		if teams[i].Name == arg {
			matches = append(matches, teams[i])
		}
	}
	switch len(matches) {
	case 0:
		return nil, fmt.Errorf("no workspace named %q", arg)
	case 1:
		return &matches[0], nil
	default:
		var b strings.Builder
		fmt.Fprintf(&b, "ambiguous: %d workspaces named %q, please switch by id:\n", len(matches), arg)
		for _, m := range matches {
			role := ""
			if m.MyRole != nil {
				role = " (" + m.MyRole.Display() + ")"
			}
			fmt.Fprintf(&b, "    %s workspace switch %s%s\n", invocationName(), m.ID, role)
		}
		return nil, errors.New(strings.TrimRight(b.String(), "\n"))
	}
}

// invocationName returns the basename of the running binary (e.g. "zeabur"
// for a native install or "zeabur" / "npx zeabur" — we can't recover the npx
// wrapper from inside Go, so we only ever return the basename of os.Args[0]).
// Falls back to "zeabur" if os.Args is empty (impossible in normal use, but
// safe in tests).
func invocationName() string {
	if len(os.Args) == 0 || os.Args[0] == "" {
		return "zeabur"
	}
	name := filepath.Base(os.Args[0])
	if name == "" || name == "." || name == "/" {
		return "zeabur"
	}
	return name
}

// isObjectIDHex matches the 24-char hex shape (case-insensitive) that the
// backend uses for primitive.ObjectID. We don't go through
// primitive.ObjectIDFromHex to avoid pulling in the mongo driver from cli/.
func isObjectIDHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
