package expose

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/util"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id   string // use id or name to specify service
	name string

	environmentID string // environmentID is required
}

func NewCmdExpose(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "expose",
		Short: "Expose a service temporarily",
		Long: `Expose a service temporarily, default 3600 seconds.
example:
      zeabur service expose # cli will try to get service from context or prompt to select one
	  zeabur service expose --id xxxxx --env-id xxxx # use id and env-id to expose service
      zeabur service expose --name xxxxx --env-id xxxx # if project context is set, use name, env-id to expose service
`,
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runExpose(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

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
	projectID := f.Config.GetContext().GetProject().GetID()

	tempTCPPort, err := f.ApiClient.ExposeService(ctx, opts.id, opts.environmentID, projectID, opts.name)
	if err != nil {
		return fmt.Errorf("failed to expose service: %w", err)
	}

	f.Log.Infof("Service is exposed on port %d, and will be closed after %d seconds\n",
		tempTCPPort.NodePort, tempTCPPort.RemainSeconds)

	return nil
}

func runExposeInteractive(f *cmdutil.Factory, opts *Options) error {
	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(
		f.Config.GetContext(), &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	return runExposeNonInteractive(f, opts)
}

func paramCheck(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("please specify --id or --name")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("please specify --env-id")
	}

	return nil
}
