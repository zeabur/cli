package deploy

import (
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	name string
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy local codes to Zeabur with one command",
		PreRunE: util.NeedProjectContextWhenNonInteractive(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Service name")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runDeployInteractive(f, opts)
	} else {
		return runDeployNonInteractive(f, opts)
	}
}

func runDeployNonInteractive(f *cmdutil.Factory, opts *Options) error {
	f.Log.Info("Deploying...")
	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching repository information..."),
	)
	s.Start()

	f.Log.Info("Deploying interactive...")

	s.Stop()
	return nil
}
