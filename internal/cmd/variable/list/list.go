package list

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdListVariables(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list environment variables",
		Long:    `List environment variables of a service`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListVariables(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runListVariables(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runListVariablesInteractive(f, opts)
	} else {
		return runListVariablesNonInteractive(f, opts)
	}
}

func runListVariablesInteractive(f *cmdutil.Factory, opts *Options) error {
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

	return runListVariablesNonInteractive(f, opts)
}

func runListVariablesNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Fetching environment variables of service %s ...", opts.name)),
	)
	s.Start()
	variableList, readonlyVariableList, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	s.Stop()

	if len(variableList) == 0 && len(readonlyVariableList) == 0 {
		f.Log.Infof("No variables found")
		return nil
	}

	f.Log.Infof("Variables of service: %s\n", opts.name)
	f.Printer.Table(variableList.Header(), variableList.Rows())

	if len(readonlyVariableList) != 0 {
		fmt.Println()
		f.Log.Infof("Readonly variables of service: %s\n", opts.name)
		f.Printer.Table(readonlyVariableList.Header(), readonlyVariableList.Rows())
	}

	return nil
}
