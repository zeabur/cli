// Package switchcmd implements `zeabur workspace switch <name|id>`. The
// package is named switchcmd because `switch` is a Go keyword.
package switchcmd

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/zcontext"
)

// NewCmdSwitch builds `zeabur workspace switch`.
func NewCmdSwitch(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "switch <name|id>",
		Short: "Switch to a team workspace",
		Long: `Switch the CLI's workspace to a team.

The argument may be the team's name or its full 24-character ObjectID. Team
names are not unique; when the name resolves to more than one team the
command exits with an error and prints the concrete invocation for each
candidate so you can pick by ID.

To return to the personal workspace use ` + "`zeabur workspace clear`" + ` —
` + "`switch personal`" + ` always looks for a team literally named "personal".

Switching clears the persisted project, environment, and service context,
because resource IDs do not overlap between workspaces.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f, args[0])
		},
	}
}

func run(f *cmdutil.Factory, arg string) error {
	teams, err := f.ListTeams(context.Background())
	if err != nil {
		return fmt.Errorf("list teams: %w", err)
	}
	team, err := cmdutil.ResolveWorkspaceArg(teams, arg)
	if err != nil {
		return err
	}

	cctx := f.Config.GetContext()
	prevProject := cctx.GetProject()

	cctx.SetWorkspace(&zcontext.Workspace{
		ID:   team.ID,
		Name: team.Name,
		Kind: zcontext.WorkspaceKindTeam,
	})
	// Resource IDs don't overlap between workspaces, so any pinned
	// project/environment/service context from the previous workspace is now
	// stale. Clear it and tell the user, so the next interactive command
	// re-prompts rather than silently 404-ing.
	cctx.ClearAll()

	role := ""
	if team.MyRole != nil {
		role = ", " + team.MyRole.Display()
	}
	fmt.Printf("Switched to workspace %q [%s] (team%s).\n", team.Name, team.ID, role)
	if !prevProject.Empty() {
		fmt.Printf("Project context cleared (was: %s).\n", prevProject.GetName())
	}
	return nil
}
