package delete

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
	deleteKeys    []string
	keys          map[string]string
	skipConfirm   bool
	inputDone     bool
}

func NewCmdDeleteVariable(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "delete variable(s)",
		Long:    `delete variable(s) for a service`,
		Aliases: []string{"del"},
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteVariable(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringArray("delete-keys", opts.deleteKeys, "Key value pair of the variable")

	return cmd
}

func runDeleteVariable(f *cmdutil.Factory, opts *Options) error {
	opts.keys = make(map[string]string)

	if f.Interactive {
		return runDeleteVariableInteractive(f, opts)
	} else {
		return runDeleteVariableNonInteractive(f, opts)
	}
}

func runDeleteVariableInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Querying variables of service: %s...", opts.name)),
	)
	s.Start()

	varList, _, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	varMap := varList.ToMap()
	keyTable := make([]string, 0, len(varMap))
	selectTable := make([]string, 0, len(varMap))
	for k, v := range varMap {
		keyTable = append(keyTable, k)
		selectTable = append(selectTable, fmt.Sprintf("%s = %s", k, v))
		opts.keys[k] = v
	}

	s.Stop()

	for !opts.inputDone {
		updateVarSelect, err := f.Prompter.Select("Select variable to delete", "", selectTable)
		if err != nil {
			return err
		}

		opts.keys[keyTable[updateVarSelect]] = ""

		doneConfirm, err := f.Prompter.Confirm("Are you done selecting variable(s) to be deleted?", false)
		if err != nil {
			return err
		}
		opts.inputDone = doneConfirm
	}

	return runDeleteVariableNonInteractive(f, opts)
}

func runDeleteVariableNonInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	for _, v := range opts.deleteKeys {
		opts.keys[v] = ""
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Deleting variables of service: %s...", opts.name)),
	)

	for k, v := range opts.keys {
		if v == "" {
			delete(opts.keys, k)
		}
	}

	createVarResult, err := f.ApiClient.UpdateVariables(context.Background(), opts.id, opts.environmentID, opts.keys)
	if err != nil {
		s.Stop()
		return err
	}
	if !createVarResult {
		s.Stop()
		return fmt.Errorf("failed to delete variables of service: %s", opts.name)
	}
	s.Stop()

	f.Log.Infof("Successfully deleted variables of service: %s\n", opts.name)

	table := make([][]string, 0, len(opts.keys))
	for k, v := range opts.keys {
		table = append(table, []string{k, v})
	}
	f.Printer.Table([]string{"Key", "Value"}, table)

	return nil
}
