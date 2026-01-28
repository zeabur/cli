package util

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/model"
)

func GetProjectByName(config config.Config, client api.Client, projectName string) (project *model.Project, err error) {
	ownerName := config.GetUsername()
	project, err = client.GetProject(context.Background(), "", ownerName, projectName)
	if err != nil {
		return nil, fmt.Errorf("get project<%s> failed: %w", projectName, err)
	}

	return project, nil
}

func AddProjectParam(cmd *cobra.Command, id, name *string) {
	cmd.Flags().StringVar(id, "id", "", "Project ID")
	cmd.Flags().StringVarP(name, "name", "n", "", "Project name")
}
