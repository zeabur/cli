package get

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	serviceID     string
	serviceName   string
	environmentID string

	deploymentID string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get deployment, if deployment-id is not specified, use serviceID/serviceName and environmentID to get the deployment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.deploymentID, "deployment-id", "", "Deployment ID")
	cmd.Flags().StringVar(&opts.serviceID, "service-id", "", "Service ID")
	cmd.Flags().StringVar(&opts.serviceName, "service-name", "", "Service Name")
	cmd.Flags().StringVar(&opts.environmentID, "env-id", "", "Environment ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	} else {
		return runGetNonInteractive(f, opts)
	}
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.deploymentID == "" {
		zctx := f.Config.GetContext()
		_, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
			ProjectCtx:    zctx,
			ServiceID:     &opts.serviceID,
			ServiceName:   &opts.serviceName,
			EnvironmentID: &opts.environmentID,
			CreateNew:     false,
		})
		if err != nil {
			return err
		}
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) (err error) {
	// If deployment ID is provided, just use it directly
	if opts.deploymentID != "" {
		deployment, err := getDeploymentByID(f, opts.deploymentID)
		if err != nil {
			return err
		}
		f.Printer.Table(deployment.Header(), deployment.Rows())
		return nil
	}

	// Resolve service ID from name
	if opts.serviceID == "" && opts.serviceName != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		opts.serviceID = service.ID
	}

	if opts.serviceID == "" {
		return errors.New("--deployment-id or --service-id/--service-name is required")
	}

	// Resolve environment from service's project
	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.serviceID)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	deployment, err := getDeploymentByServiceAndEnvironment(f, opts.serviceID, opts.environmentID)
	if err != nil {
		return err
	}

	f.Printer.Table(deployment.Header(), deployment.Rows())

	return nil
}

func getDeploymentByID(f *cmdutil.Factory, deploymentID string) (*model.Deployment, error) {
	deployment, err := f.ApiClient.GetDeployment(context.Background(), deploymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	return deployment, nil
}

func getDeploymentByServiceAndEnvironment(f *cmdutil.Factory, serviceID, environmentID string) (*model.Deployment, error) {
	deployment, exist, err := f.ApiClient.GetLatestDeployment(context.Background(), serviceID, environmentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}

	if !exist {
		return nil, errors.New("no deployment found")
	}

	return deployment, nil
}
