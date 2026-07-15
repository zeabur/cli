package update

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
	rawKeys       []string
	keys          map[string]string
	updatedKeys   []string // tracks only the keys the user explicitly changed
	skipConfirm   bool
	inputDone     bool
}

// maskValue masks a variable value for display, showing only the first 3
// characters followed by asterisks. Short values are fully masked.
func maskValue(v string) string {
	if len(v) <= 3 {
		return "***"
	}
	return v[:3] + "***"
}

func NewCmdUpdateVariable(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "update variable(s)",
		Long:  `update variable(s) for a service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if !f.Interactive {
				keys, err := cmdutil.ParseKeyValuePairs(opts.rawKeys)
				if err != nil {
					return err
				}
				opts.keys = keys
			}
			return runUpdateVariable(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringArrayVarP(&opts.rawKeys, "key", "k", nil, "Key value pair of the variable")

	return cmd
}

func runUpdateVariable(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		opts.keys = make(map[string]string)
		return runUpdateVariableInteractive(f, opts)
	}
	return runUpdateVariableNonInteractive(f, opts)
}

func runUpdateVariableInteractive(f *cmdutil.Factory, opts *Options) error {
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
		selectTable = append(selectTable, fmt.Sprintf("%s = %s", k, maskValue(v)))
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
		selectedKey := keyTable[updateVarSelect]
		opts.keys[selectedKey] = varInput

		seen := false
		for _, k := range opts.updatedKeys {
			if k == selectedKey {
				seen = true
				break
			}
		}
		if !seen {
			opts.updatedKeys = append(opts.updatedKeys, selectedKey)
		}

		doneConfirm, err := f.Prompter.Confirm("Are you done entering value of variable?", false)
		if err != nil {
			return err
		}
		opts.inputDone = doneConfirm
	}

	return runUpdateVariableNonInteractive(f, opts)
}

func runUpdateVariableNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	// Remember which keys the user explicitly requested to update,
	// so we only show those in the output (not the full merged set).
	// In interactive mode, updatedKeys is already populated.
	if len(opts.updatedKeys) == 0 {
		for k := range opts.keys {
			opts.updatedKeys = append(opts.updatedKeys, k)
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Updating variables of service: %s...", opts.name)),
	)
	s.Start()

	// Fetch existing variables and merge with user-provided keys,
	// so that unmentioned variables are preserved (not deleted).
	// In interactive mode, opts.keys already contains the full
	// variable map from the interactive flow.
	varList, _, err := f.ApiClient.ListVariables(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		s.Stop()
		return err
	}
	mergedVars := varList.ToMap()
	for k, v := range opts.keys {
		mergedVars[k] = v
	}
	opts.keys = mergedVars

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

	if f.JSON {
		out := make([]map[string]string, 0, len(opts.updatedKeys))
		for _, k := range opts.updatedKeys {
			out = append(out, map[string]string{"Key": k})
		}
		return f.Printer.JSON(out)
	}

	f.Log.Infof("Successfully updated variables of service: %s\n", opts.name)

	table := make([][]string, 0, len(opts.updatedKeys))
	for _, k := range opts.updatedKeys {
		table = append(table, []string{k})
	}
	f.Printer.Table([]string{"Key"}, table)

	return nil
}
