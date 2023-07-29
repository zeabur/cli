package get

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string
	name string

	environmentID string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	cmd := &cobra.Command{
		Use:     "get",
		Short:   "Get a service, if environment is specified, get the service details in the environment",
		PreRunE: util.NeedProjectContext(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.name, "name", ctx.GetService().GetName(), "Service name")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Environment ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}

	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.ServiceByName(f.Config.GetContext(), &opts.id, &opts.name); err != nil {
		return err
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	projectName := f.Config.GetContext().GetProject().GetName()
	username := f.Config.GetUsername()

	var (
		t   model.Tabler
		err error
	)

	if opts.environmentID == "" {
		t, err = getServiceBrief(f.ApiClient, opts.id, username, projectName, opts.name)
	} else {
		t, err = getServiceDetails(f.ApiClient, opts.id, username, projectName, opts.name, opts.environmentID)
	}

	if err != nil {
		return err
	}

	f.Printer.Table(t.Header(), t.Rows())

	return nil
}

func getServiceBrief(client api.ServiceAPI, id, username, projectName, name string) (t model.Tabler, err error) {
	ctx := context.Background()
	service, err := client.GetService(ctx, id, username, projectName, name)
	if err != nil {
		return nil, fmt.Errorf("get service failed: %w", err)
	}

	return service, nil
}

func getServiceDetails(client api.ServiceAPI, id, username, projectID, name, environmentID string) (t model.Tabler, err error) {
	ctx := context.Background()
	serviceDetail, err := client.GetServiceDetailByEnvironment(ctx, id, username, projectID, name, environmentID)
	if err != nil {
		return nil, fmt.Errorf("get service failed: %w", err)
	}

	return serviceDetail, nil
}
