package util

import (
	"context"
	"fmt"

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
