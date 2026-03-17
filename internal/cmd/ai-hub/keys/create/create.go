package create

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	alias string
}

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new AI Hub API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.alias, "alias", "", "Alias for the API key")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts Options) error {
	if opts.alias == "" && f.Interactive {
		alias, err := f.Prompter.Input("Key alias (optional): ", "")
		if err != nil {
			return err
		}
		opts.alias = alias
	}

	var aliasPtr *string
	if opts.alias != "" {
		aliasPtr = &opts.alias
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Creating AI Hub key..."),
	)
	s.Start()
	result, err := f.ApiClient.CreateAIHubKey(context.Background(), aliasPtr)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(result)
	}

	f.Log.Infof("Key created successfully!")
	f.Log.Infof("Key ID: %s", result.Key.KeyID)
	if result.Key.Alias != "" {
		f.Log.Infof("Alias: %s", result.Key.Alias)
	}
	fmt.Printf("API Key: %s\n", result.APIKey)
	f.Log.Infof("WARNING: This API key will only be shown once. Please save it now.")

	return nil
}
