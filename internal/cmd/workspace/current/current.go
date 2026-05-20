// Package current implements `zeabur workspace current`.
package current

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdCurrent builds `zeabur workspace current`.
func NewCmdCurrent(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "current",
		Short: "Show the workspace the CLI is currently acting under",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f)
		},
	}
}

func run(f *cmdutil.Factory) error {
	// Use the effective workspace (CurrentWorkspace) so the report honours
	// a `--workspace` flag override for the current invocation, matching
	// what every other command in the process sees. Reading the persisted
	// workspace directly would silently lie about the override.
	ws := f.CurrentWorkspace()
	if ws.IsPersonal() {
		label := f.Config.GetUser()
		if label == "" {
			label = f.Config.GetUsername()
		}
		if label == "" {
			label = "(you)"
		}
		fmt.Printf("personal  (%s)\n", label)
		return nil
	}

	// For a team workspace also fetch the freshest role from the backend so
	// `current` reports the live role rather than whatever was cached when
	// the workspace was last switched in. Don't fail the command on
	// network error — the persisted name/id is still useful — but log it
	// at debug so it isn't completely silent (PLA-1590 review N3).
	role := ""
	teams, err := f.ListTeams(context.Background())
	if err != nil {
		f.Log.Debugf("workspace current: list teams failed, omitting role: %v", err)
	} else {
		for _, t := range teams {
			if t.ID == ws.ID && t.MyRole != nil {
				role = "  " + t.MyRole.Display()
				break
			}
		}
	}
	fmt.Printf("%s  [%s]  team%s\n", ws.Name, ws.ID, role)
	return nil
}
