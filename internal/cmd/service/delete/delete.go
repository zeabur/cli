package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id   string
	name string

	environmentID string

	skipConfirm bool
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a service",
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runDeleteInteractive(f, opts)
	}
	return runDeleteNonInteractive(f, opts)
}

func runDeleteInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    zctx,
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	return runDeleteNonInteractive(f, opts)
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.environmentID == "" {
		projectID := f.Config.GetContext().GetProject().GetID()
		envID, err := util.ResolveEnvironmentID(f.ApiClient, projectID)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	if err := checkParams(opts); err != nil {
		return err
	}

	// if name is set, get service id by name
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	// to show friendly message
	idOrName := opts.name
	if idOrName == "" {
		idOrName = opts.id
	}

	// double check
	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to delete service <%s>?", idOrName), false)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	} else if !f.Interactive && !opts.skipConfirm {
		f.Log.Info("Please use --yes to confirm deletion without interactive prompt")
		return nil
	}

	err := f.ApiClient.DeleteService(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return fmt.Errorf("delete service failed: %w", err)
	}

	f.Log.Infof("Service <%s> deleted successfully", idOrName)

	return nil
}

func checkParams(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("--id or --name is required")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("--env-id is required")
	}

	return nil
}
