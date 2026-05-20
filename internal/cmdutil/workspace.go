package cmdutil

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

// ResolveWorkspaceArg turns user-supplied "<name|id>" into a workspace ID by
// asking the backend's `teams` query for the caller's memberships. The same
// resolution rules apply to both `zeabur workspace switch <arg>` and the
// global `--workspace <arg>` flag:
//
//   - 24-char hex → treated as a team ObjectID directly. The caller is still
//     a member only if that team appears in the `teams` reply, but for the
//     flag path we trust the input: backend RBAC will reject if it isn't.
//   - non-hex → matched by team name against the membership list.
//   - exactly one match → return that team's ID.
//   - zero matches → "no workspace named ..." error.
//   - two or more matches → "ambiguous" error listing each candidate with the
//     concrete `zeabur workspace switch <id>` invocation that disambiguates
//     it. Team names are not unique; users must use the ID to pick.
//
// Returns the resolved Team so callers can both record the name in config
// (for `workspace switch`) and surface it in the deploy-to-team hint.
func ResolveWorkspaceArg(ctx context.Context, client api.Client, arg string) (*model.Team, error) {
	arg = strings.TrimSpace(arg)
	if arg == "" {
		return nil, errors.New("workspace name or id is required")
	}

	teams, err := client.ListTeams(ctx)
	if err != nil {
		return nil, fmt.Errorf("list teams: %w", err)
	}

	if isObjectIDHex(arg) {
		for i := range teams {
			if teams[i].ID == arg {
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
			fmt.Fprintf(&b, "    zeabur workspace switch %s%s\n", m.ID, role)
		}
		return nil, errors.New(strings.TrimRight(b.String(), "\n"))
	}
}

// isObjectIDHex matches the 24-char lowercase/uppercase hex shape that the
// backend uses for primitive.ObjectID. We don't go through
// primitive.ObjectIDFromHex to avoid pulling in the mongo driver from cli/.
func isObjectIDHex(s string) bool {
	if len(s) != 24 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil
}
