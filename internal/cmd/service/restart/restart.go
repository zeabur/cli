package restart

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

func NewCmdRestart(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "restart",
		Short: "restart a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestart(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runRestart(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runRestartInteractive(f, opts)
	} else {
		return runRestartNonInteractive(f, opts)
	}
}

func runRestartInteractive(f *cmdutil.Factory, opts *Options) error {
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

	return runRestartNonInteractive(f, opts)
}

func runRestartNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.id)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	// to show friendly message
	idOrName := opts.name
	if idOrName == "" {
		idOrName = opts.id
	}

	// double check
	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to deploy service <%s>?", idOrName), true)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	err := f.ApiClient.RestartService(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return fmt.Errorf("restart service failed: %w", err)
	}

	f.Log.Infof("Service <%s> restarted successfully", idOrName)

	return nil
}
