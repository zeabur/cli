package list

import (
	"context"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	serviceID     string
	serviceName   string
	environmentID string
}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List deployments",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.serviceID, "service-id", "", "Service ID")
	cmd.Flags().StringVar(&opts.serviceName, "service-name", "", "Service Name")
	cmd.Flags().StringVar(&opts.environmentID, "env-id", "", "Environment ID")

	return cmd
}

func runList(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runListInteractive(f, opts)
	} else {
		return runListNonInteractive(f, opts)
	}
}

func runListInteractive(f *cmdutil.Factory, opts *Options) error {
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

	return runListNonInteractive(f, opts)
}

func runListNonInteractive(f *cmdutil.Factory, opts *Options) error {
	// Resolve service ID from name
	if opts.serviceID == "" && opts.serviceName != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.serviceName)
		if err != nil {
			return fmt.Errorf("failed to get service: %w", err)
		}
		opts.serviceID = service.ID
	}

	if opts.serviceID == "" {
		return errors.New("--service-id or --service-name is required")
	}

	// Resolve environment from service's project
	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.serviceID)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	deployments, err := f.ApiClient.ListAllDeployments(context.Background(), opts.serviceID, opts.environmentID)
	if err != nil {
		return fmt.Errorf("failed to list deployments: %w", err)
	}

	if len(deployments) == 0 {
		f.Log.Info("No deployments found")
		return nil
	}

	f.Printer.Table(deployments.Header(), deployments.Rows())

	return nil
}
