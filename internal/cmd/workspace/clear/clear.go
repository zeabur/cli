// Package clear implements `zeabur workspace clear`. Clear is the ONLY way
// to return to the personal workspace — `workspace switch personal` is
// intentionally interpreted as "find a team literally named personal" because
// team names are unconstrained.
package clear

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdClear builds `zeabur workspace clear`.
func NewCmdClear(f *cmdutil.Factory) *cobra.Command {
	return &cobra.Command{
		Use:   "clear",
		Short: "Switch back to the personal workspace",
		Long: `Return to the personal workspace.

This is the only way to switch to personal: `+"`workspace switch personal`"+`
always looks for a team literally named "personal". This also clears the
persisted project, environment, and service context because resource IDs do
not overlap between workspaces.`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return run(f)
		},
	}
}

func run(f *cmdutil.Factory) error {
	cctx := f.Config.GetContext()
	prev := cctx.GetWorkspace()
	prevProject := cctx.GetProject()

	cctx.ClearWorkspace()
	cctx.ClearAll()

	if prev.IsPersonal() {
		fmt.Println("Already on personal workspace.")
	} else {
		fmt.Printf("Switched to personal workspace (was: %s).\n", prev.Name)
	}
	if !prevProject.Empty() {
		fmt.Printf("Project context cleared (was: %s).\n", prevProject.GetName())
	}
	return nil
}
