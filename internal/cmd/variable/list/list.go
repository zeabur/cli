package list

import (
	"context"
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdListVariables(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list environment variables",
		Long:    `List environment variables of a service`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
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

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	return runListVariablesNonInteractive(f, opts)
}

func runListVariablesNonInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Fetching environment variablesof service %s ...", opts.name)),
	)
	s.Start()
	variableList, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	s.Stop()

	if len(variableList) == 0 {
		f.Log.Infof("No variables found")
		return nil
	}

	f.Printer.Table(variableList.Header(), variableList.Rows())

	return nil
}
