package util

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

// GetProjectByName resolves a project by name within the active workspace.
//
//   - Personal workspace (ownerID == ""): uses the cheap
//     `project(owner: $personalUsername, name: $projectName)` query against
//     the caller's account.
//   - Team workspace (ownerID != ""): there is no `project(ownerID, name)`
//     query, so walk the owner's project list and match by name locally.
//
// Without the team branch a `--name`-based lookup in a team workspace would
// silently fall through to the caller's personal account — either failing
// to find the project, or worse, returning a same-named personal project
// instead of the intended team one. The personal path is unchanged, so
// existing personal-workspace callers see no behavior change.
func GetProjectByName(client api.Client, ownerID, personalUsername, projectName string) (*model.Project, error) {
	ctx := context.Background()
	if ownerID == "" {
		project, err := client.GetProject(ctx, "", personalUsername, projectName)
		if err != nil {
			return nil, fmt.Errorf("get project<%s> failed: %w", projectName, err)
		}
		return project, nil
	}

	projects, err := client.ListAllProjects(ctx, ownerID)
	if err != nil {
		return nil, fmt.Errorf("list projects in workspace: %w", err)
	}
	for _, p := range projects {
		if p.Name == projectName {
			return p, nil
		}
	}
	return nil, fmt.Errorf("no project named %q in this workspace", projectName)
}

func AddProjectParam(cmd *cobra.Command, id, name *string) {
	cmd.Flags().StringVar(id, "id", "", "Project ID")
	cmd.Flags().StringVarP(name, "name", "n", "", "Project name")
}
