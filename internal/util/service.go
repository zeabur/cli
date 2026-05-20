package util

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

// GetServiceByName resolves a service by name within the active workspace.
//
//   - Personal workspace (ownerID == ""): uses the cheap
//     `service(owner, projectName, name)` query.
//   - Team workspace (ownerID != ""): the personal query keys on the
//     caller's username and would silently look at the personal account.
//     Use the project-scoped `ListAllServices(projectID)` and match by name
//     locally — projectID is unique across owners and team-safe.
//
// In a team workspace the caller must have a project context (projectID),
// because services have no top-level owner-scoped lookup and a service name
// is only unique within a project. Personal path is unchanged.
func GetServiceByName(client api.Client, ownerID, personalUsername, projectName, projectID, serviceName string) (*model.Service, error) {
	ctx := context.Background()
	if ownerID == "" {
		service, err := client.GetService(ctx, "", personalUsername, projectName, serviceName)
		if err != nil {
			return nil, fmt.Errorf("get service<%s> failed: %w", serviceName, err)
		}
		return service, nil
	}

	if projectID == "" {
		return nil, fmt.Errorf("cannot resolve service by name in a team workspace without a project context — set a project first (e.g. `zeabur context set project --id <project-id>`)")
	}
	services, err := client.ListAllServices(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("list services in project: %w", err)
	}
	for _, s := range services {
		if s.Name == serviceName {
			return s, nil
		}
	}
	return nil, fmt.Errorf("no service named %q in this project", serviceName)
}

func AddServiceParam(cmd *cobra.Command, id, name *string) {
	cmd.Flags().StringVar(id, "id", "", "Service ID")
	cmd.Flags().StringVarP(name, "name", "n", "", "Service name")
}
