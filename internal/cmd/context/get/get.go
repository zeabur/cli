package get

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	project := f.Config.GetContext().GetProject()
	if project.Empty() {
		f.Log.Info("Project: None")
	} else {
		f.Log.Infof("Project: Name: %s, ID: %s", project.GetName(), project.GetID())
	}

	environment := f.Config.GetContext().GetEnvironment()
	if environment.Empty() {
		f.Log.Info("Environment: None")
	} else {
		f.Log.Infof("Environment: Name: %s, ID: %s", environment.GetName(), environment.GetID())
	}

	service := f.Config.GetContext().GetService()
	if service.Empty() {
		f.Log.Info("Service: None")
	} else {
		f.Log.Infof("Service: Name: %s, ID: %s", service.GetName(), service.GetID())
	}

	return nil
}
