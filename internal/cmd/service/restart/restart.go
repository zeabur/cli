package restart

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
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
		Short: "Restart a service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRestart(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.name, "name", ctx.GetService().GetName(), "Service name")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Environment ID")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runRestart(f *cmdutil.Factory, opts *Options) error {
	// when params are missing, we need to use project id to get service id and environment id
	if checkParams(opts) != nil {
		projectCtx := f.Config.GetContext().GetProject()
		if projectCtx.Empty() {
			return fmt.Errorf("please set context project first")
		}
	}

	if f.Interactive {
		return runRestartInteractive(f, opts)
	} else {
		return runRestartNonInteractive(f, opts)
	}
}

func runRestartInteractive(f *cmdutil.Factory, opts *Options) error {

	projectID := f.Config.GetContext().GetProject().GetID()

	// fill in service id
	if opts.id == "" && opts.name == "" {
		serviceInfo, _, err := f.Selector.SelectService(projectID)
		if err != nil {
			return err
		}
		opts.id = serviceInfo.GetID()
	}

	// fill in environment id
	if opts.environmentID == "" {
		environmentInfo, _, err := f.Selector.SelectEnvironment(projectID)
		if err != nil {
			return err
		}
		opts.environmentID = environmentInfo.GetID()
	}

	return runRestartNonInteractive(f, opts)
}

func runRestartNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := checkParams(opts); err != nil {
		return err
	}

	// if only name is provided, we need to get service id by project id first
	if opts.id == "" && opts.name != "" {
		projectCtx := f.Config.GetContext().GetProject()
		if projectCtx.Empty() {
			return fmt.Errorf("when using service name, please set context project first")
		}
		ownerName := f.Config.GetUsername()
		projectName := projectCtx.GetName()
		service, err := f.ApiClient.GetService(context.Background(), "", ownerName, projectName, opts.name)
		if err != nil {
			return fmt.Errorf("get service by name failed: %w", err)
		}
		opts.id = service.ID
	}

	// double check
	if f.Interactive && !opts.skipConfirm {
		idOrName := opts.name
		if idOrName == "" {
			idOrName = opts.id
		}
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to restart service <%s>?", idOrName), true)
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

	return nil
}

func checkParams(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("--id or --name is required")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("--environment-id is required")
	}

	return nil
}
