package exec

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id   string
	name string

	environmentID string

	command []string
}

func NewCmdExec(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "exec -- <command> [args...]",
		Short: "Execute a command in a service container",
		Long:  "Execute a command in a running service's container. The command and its arguments should be specified after \"--\".",
		Example: `  # List files in the service container
  zeabur service exec -- ls -la

  # Run a shell command
  zeabur service exec -- sh -c "echo hello"

  # Specify service by name
  zeabur service exec --name my-svc --env-id xxx -- cat /etc/hostname`,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.command = cmd.Flags().Args()
			if len(opts.command) == 0 {
				return fmt.Errorf("command is required, use -- to separate it from flags")
			}
			return runExec(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runExec(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runExecInteractive(f, opts)
	}
	return runExecNonInteractive(f, opts)
}

func runExecInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name == "" {
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
	}

	return runExecNonInteractive(f, opts)
}

func runExecNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	result, err := f.ApiClient.ExecuteCommand(context.Background(), opts.id, opts.environmentID, opts.command)
	if err != nil {
		return fmt.Errorf("execute command failed: %w", err)
	}

	fmt.Print(result.Output)

	if result.ExitCode != 0 {
		os.Exit(result.ExitCode)
	}

	return nil
}
