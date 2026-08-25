package list

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	variablevalue "github.com/zeabur/cli/internal/cmd/variable/value"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	id            string
	name          string
	environmentID string
	showValues    bool
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
	cmd.Flags().BoolVar(&opts.showValues, "show-values", false, "Show full variable values (may expose secrets)")

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
	zctx := f.EffectiveContext()

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
		service, err := util.GetServiceByName(f.ApiClient, f.CurrentOwnerID(), f.Config.GetUsername(), f.CurrentProjectName(), f.CurrentProjectID(), opts.name)
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
	if !opts.showValues {
		variableList = maskVariables(variableList)
		readonlyVariableList = maskVariables(readonlyVariableList)
	}

	if len(variableList) == 0 && len(readonlyVariableList) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No variables found")
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(map[string]any{"variables": variableList, "readonlyVariables": readonlyVariableList})
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

func maskVariables(variables model.Variables) model.Variables {
	masked := make(model.Variables, 0, len(variables))
	for _, variable := range variables {
		maskedVariable := *variable
		maskedVariable.Value = variablevalue.Mask(variable.Value)
		masked = append(masked, &maskedVariable)
	}
	return masked
}
