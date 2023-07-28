package expose

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string // use id to specify service

	projectID string // use projectID and serviceName to specify service
	name      string

	environmentID string // environmentID is required
}

func NewCmdExpose(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{
		projectID: f.Config.GetContext().GetProject().GetID(),
	}

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Expose a service temporarily",
		Long: `Expose a service temporarily, default 3600 seconds.
example:
      zeabur service expose # cli will try to get service from context or prompt to select one
	  zeabur service expose --id xxxxx --environment-id xxxx # use id and environment-id to expose service
      zeabur service expose --name xxxxx --environment-id xxxx # if project context is set, use name, environment-id to expose service
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpose(f, opts)
		},
	}

	ctx := f.Config.GetContext()

	cmd.Flags().StringVar(&opts.id, "id", ctx.GetService().GetID(), "Service ID")
	cmd.Flags().StringVar(&opts.name, "name", ctx.GetService().GetName(), "Service name")
	cmd.Flags().StringVar(&opts.environmentID, "environment-id", ctx.GetEnvironment().GetID(), "Service environment ID")

	return cmd
}

func runExpose(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runExposeInteractive(f, opts)
	} else {
		return runExposeNonInteractive(f, opts)
	}
}

func runExposeNonInteractive(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()

	tempTCPPort, err := f.ApiClient.ExposeService(ctx, opts.id, opts.environmentID, opts.projectID, opts.name)
	if err != nil {
		return fmt.Errorf("failed to expose service: %w", err)
	}

	f.Log.Infof("Service is exposed on port %d, and will be closed after %d seconds\n",
		tempTCPPort.NodePort, tempTCPPort.RemainSeconds)

	return nil
}

func runExposeInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(
		&opts.projectID, &opts.id, &opts.name, &opts.environmentID); err != nil {
	}

	return runExposeNonInteractive(f, opts)
}

func paramCheck(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("please specify --environment-id")
	}

	return nil
}
