package util

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/config"
	"github.com/zeabur/cli/pkg/model"
)

func GetServiceByName(config config.Config, client api.Client, serviceName string) (service *model.Service, err error) {
	ownerName := config.GetUsername()
	projectName := config.GetContext().GetProject().GetName()
	service, err = client.GetService(context.Background(), "", ownerName, projectName, serviceName)
	if err != nil {
		return nil, fmt.Errorf("get service<%s> failed: %w", serviceName, err)
	}

	return service, nil
}

func AddServiceParam(cmd *cobra.Command, id, name *string) {
	cmd.Flags().StringVar(id, "id", "", "Service ID")
	cmd.Flags().StringVarP(name, "name", "n", "", "Service name")
}
