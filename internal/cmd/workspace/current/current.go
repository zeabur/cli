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
	ws := f.Config.GetContext().GetWorkspace()
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
	// the workspace was last switched in.
	role := ""
	teams, err := f.ApiClient.ListTeams(context.Background())
	if err == nil {
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
