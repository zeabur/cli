package get

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/fill"
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
		Use:   "get",
		Short: "Get a service, if environment is specified, get the service details in the environment",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}

	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.ServiceByName(fill.ServiceByNameOptions{
		ProjectCtx:  f.Config.GetContext(),
		ServiceID:   &opts.id,
		ServiceName: &opts.name,
	}); err != nil {
		return err
	}

	return runGetNonInteractive(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("--id or --name is required")
	}

	var (
		t   model.Tabler
		err error
	)

	if opts.environmentID == "" {
		t, err = getServiceBrief(f.ApiClient, opts.id)
	} else {
		t, err = getServiceDetails(f.ApiClient, opts.id, opts.environmentID)
	}

	if err != nil {
		return err
	}

	f.Printer.Table(t.Header(), t.Rows())

	return nil
}

func getServiceBrief(client api.ServiceAPI, id string) (t model.Tabler, err error) {
	ctx := context.Background()
	service, err := client.GetService(ctx, id, "", "", "")
	if err != nil {
		return nil, fmt.Errorf("get service failed: %w", err)
	}

	return service, nil
}

func getServiceDetails(client api.ServiceAPI, id, environmentID string) (t model.Tabler, err error) {
	ctx := context.Background()
	serviceDetail, err := client.GetServiceDetailByEnvironment(ctx, id, "", "", "", environmentID)
	if err != nil {
		return nil, fmt.Errorf("get service failed: %w", err)
	}

	return serviceDetail, nil
}
