package update

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
	keys          map[string]string
	skipConfirm   bool
	inputDone     bool
}

func NewCmdUpdateVariable(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update variable(s)",
		Long:  `update variable(s) for a service`,
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

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Querying variables of service: %s...", opts.name)),
	)
	s.Start()

	varList, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	varMap := varList.ToMap()
	var keyTable, selectTable []string
	for k, v := range varMap {
		keyTable = append(keyTable, k)
		selectTable = append(selectTable, fmt.Sprintf("%s = %s", k, v))
		opts.keys[k] = v
	}

	s.Stop()

	for !opts.inputDone {
		updateVarSelect, err := f.Prompter.Select("Select variable to update", "", selectTable)
		if err != nil {
			return err
		}

		varInput, err := f.Prompter.Input("Enter selected variable value modified: ", "")
		if err != nil {
			return err
		}
		opts.keys[keyTable[updateVarSelect]] = varInput

		doneConfirm, err := f.Prompter.Confirm("Are you done entering value of variable?", false)
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
		spinner.WithSuffix(fmt.Sprintf(" Updating variables of service: %s...", opts.name)),
	)
	createVarResult, err := f.ApiClient.UpdateVariables(context.Background(), opts.id, opts.environmentID, opts.keys)
	if err != nil {
		s.Stop()
		return err
	}
	if !createVarResult {
		s.Stop()
		return fmt.Errorf("failed to update variables of service: %s", opts.name)
	}
	s.Stop()

	f.Log.Infof("Successfully updated variables of service: %s", opts.name)

	var table [][]string
	for k, v := range opts.keys {
		table = append(table, []string{k, v})
	}
	f.Printer.Table([]string{"Key", "Value"}, table)

	return nil
}
