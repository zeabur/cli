package env

import (
	"context"
	"fmt"
	"os"

	"github.com/briandowns/spinner"
	"github.com/hashicorp/go-envparse"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id            string
	name          string
	environmentID string
	envFilename   string
}

func NewCmdEnvVariable(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "env",
		Short: "update variables from .env",
		Long:  "overwrite variables from a .env file",
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateVariableByEnv(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().StringVarP(&opts.envFilename, "file", "f", ".env", "Path to the .env file")

	return cmd
}

func runUpdateVariableByEnv(f *cmdutil.Factory, opts *Options) error {
	if _, err := os.Stat(opts.envFilename); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", opts.envFilename)
		}

		return fmt.Errorf("file cannot open: %s (%w)", opts.envFilename, err)
	}

	return runUpdateVariableNonInteractive(f, opts)
}

func runUpdateVariableNonInteractive(f *cmdutil.Factory, opts *Options) error {
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

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Updating variables of service: %s...", opts.name)),
	)

	// read files from .env
	envFile, err := os.Open(opts.envFilename)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}
	defer func() {
		_ = envFile.Close()
	}()

	envMap, err := envparse.Parse(envFile)
	if err != nil {
		return fmt.Errorf("parse env file: %w", err)
	}

	createVarResult, err := f.ApiClient.UpdateVariables(context.Background(), opts.id, opts.environmentID, envMap)
	if err != nil {
		s.Stop()
		return err
	}
	if !createVarResult {
		s.Stop()
		return fmt.Errorf("failed to update variables of service: %s", opts.name)
	}
	s.Stop()

	f.Log.Infof("Successfully updated variables of service: %s\n\tRestart your service manually to apply the changes.\n", opts.name)

	table := make([][]string, 0, len(envMap))
	for k, v := range envMap {
		table = append(table, []string{k, v})
	}
	f.Printer.Table([]string{"Key", "Value"}, table)

	return nil
}
