package create

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id            string
	name          string
	environmentID string
	keys          map[string]string
	skipConfirm   bool
	inputDone     bool
}

func NewCmdCreateVariable(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create variable(s)",
		Long:  `Create variable(s) for a service`,
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateVariable(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringToStringVarP(&opts.keys, "key", "k", nil, "Key value pair of the variable")

	return cmd
}

func runCreateVariable(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runCreateVariableInteractive(f, opts)
	} else {
		return runCreateVariableNonInteractive(f, opts)
	}
}

func runCreateVariableInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	opts.keys = make(map[string]string)

	for !opts.inputDone {
		varInput, err := f.Prompter.Input("Enter variable key value pair (key=value)", "")
		if err != nil {
			return err
		}
		keyValue := strings.Split(varInput, "=")
		if len(keyValue) != 2 {
			return fmt.Errorf("invalid input")
		}
		opts.keys[keyValue[0]] = keyValue[1]

		doneConfirm, err := f.Prompter.Confirm("Are you done entering variables?", false)
		if err != nil {
			return err
		}
		opts.inputDone = doneConfirm
	}

	return runCreateVariableNonInteractive(f, opts)
}

func runCreateVariableNonInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Creating variables of service: %s...", opts.name)),
	)
	s.Start()

	varList, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}

	varMap := varList.ToMap()
	for k, v := range opts.keys {
		if _, ok := varMap[k]; ok {
			s.Stop()
			return fmt.Errorf("variable %s already exists", k)
		}
		varMap[k] = v
	}
	createVarResult, err := f.ApiClient.UpdateVariables(context.Background(), opts.id, opts.environmentID, varMap)
	if err != nil {
		s.Stop()
		return err
	}
	if !createVarResult {
		s.Stop()
		return fmt.Errorf("failed to create variables of service: %s", opts.name)
	}
	s.Stop()

	f.Log.Infof("Successfully created variables of service: %s", opts.name)
	
	table := make([][]string, 0, len(varMap))
	for k, v := range varMap {
		table = append(table, []string{k, v})
	}
	f.Printer.Table([]string{"Key", "Value"}, table)

	return nil
}
