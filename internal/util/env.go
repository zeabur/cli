package util

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/api"
)

// AddEnvParam todo: support name
func AddEnvParam(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVar(id, "id", "", "Environment ID")
}

func AddEnvOfServiceParam(cmd *cobra.Command, id *string) {
	cmd.Flags().StringVar(id, "env-id", "", "Environment ID of service")
}

// ResolveEnvironmentID resolves the environment ID from the project ID
// by listing environments and returning the first one.
// Every project has exactly one environment since environments are deprecated.
func ResolveEnvironmentID(client api.Client, projectID string) (string, error) {
	if projectID == "" {
		return "", fmt.Errorf("project ID is required to resolve environment ID; please set project context with `zeabur context set project`")
	}

	environments, err := client.ListEnvironments(context.Background(), projectID)
	if err != nil {
		return "", fmt.Errorf("failed to list environments for project %s: %w", projectID, err)
	}

	if len(environments) == 0 {
		return "", fmt.Errorf("no environment found for project %s", projectID)
	}

	return environments[0].ID, nil
}
