package get

import (
	"context"
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
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
		Use:     "get",
		Short:   "Get deployment, if deployment-id is not specified, use serviceID/serviceName and environmentID to get the deployment",
		PreRunE: util.NeedProjectContext(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	zctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.deploymentID, "deployment-id", "", "Deployment ID")
	cmd.Flags().StringVar(&opts.serviceID, "service-id", zctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.serviceName, "service-name", zctx.GetService().GetName(), "Service Name")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", zctx.GetEnvironment().GetID(), "Environment ID")

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
		_, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.serviceID, &opts.serviceName, &opts.environmentID)
		if err != nil {
			return err
		}
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) (err error) {
	if err = paramCheck(opts); err != nil {
		return err
	}

	var deployment *model.Deployment

	// If deployment id is provided, get deployment by deployment id
	if opts.deploymentID != "" {
		deployment, err = getDeploymentByID(f, opts.deploymentID)
	} else {
		// or, get deployment by service id and environment id

		// If service id is not provided, get service id by service name
		if opts.serviceID == "" {
			var service *model.Service
			if service, err = util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName); err != nil {
				return fmt.Errorf("failed to get service: %w", err)
			} else {
				opts.serviceID = service.ID
			}
		}
		deployment, err = getDeploymentByServiceAndEnvironment(f, opts.serviceID, opts.environmentID)
	}

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

func paramCheck(opts *Options) error {
	if opts.deploymentID != "" {
		return nil
	}

	if opts.serviceID == "" && opts.serviceName == "" {
		return errors.New("when deployment-id is not specified, service-id or service-name is required")
	}

	if opts.environmentID == "" {
		return errors.New("when deployment-id is not specified, environment-id is required")
	}

	return nil
}
