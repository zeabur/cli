// Package workspace contains the `zeabur workspace` command group, which
// drives the personal-vs-team-workspace switcher the CLI uses for
// directory-level operations (project list / create, deploy with no linked
// project). Specific resources (services / deployments / variables) are
// addressed by ID and stay workspace-independent.
package workspace

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"

	workspaceClearCmd "github.com/zeabur/cli/internal/cmd/workspace/clear"
	workspaceCurrentCmd "github.com/zeabur/cli/internal/cmd/workspace/current"
	workspaceListCmd "github.com/zeabur/cli/internal/cmd/workspace/list"
	workspaceSwitchCmd "github.com/zeabur/cli/internal/cmd/workspace/switch"
)

// NewCmdWorkspace builds the `zeabur workspace` parent command.
func NewCmdWorkspace(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "workspace",
		Short: "Manage the personal / team workspace the CLI acts under",
		Long: `Manage the workspace the CLI uses for project list / create / deploy.

A workspace is either "personal" (the caller's own account) or a team. The
choice affects only directory-level operations: which projects are listed
and which owner a newly-created project is filed under. Operations on a
specific project, service, or deployment use that resource's own owner.

Switching workspaces clears the persisted project / environment / service
context because resource IDs do not overlap between workspaces.`,
	}

	cmd.AddCommand(workspaceListCmd.NewCmdList(f))
	cmd.AddCommand(workspaceCurrentCmd.NewCmdCurrent(f))
	cmd.AddCommand(workspaceSwitchCmd.NewCmdSwitch(f))
	cmd.AddCommand(workspaceClearCmd.NewCmdClear(f))

	return cmd
}
