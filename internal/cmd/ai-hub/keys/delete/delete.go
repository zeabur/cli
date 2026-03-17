package delete

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	keyID string
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete an AI Hub API key",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.keyID, "key-id", "", "ID of the key to delete")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runDeleteInteractive(f, opts)
	}
	return runDeleteNonInteractive(f, opts)
}

func runDeleteInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.keyID == "" {
		// Fetch keys to let user select
		tenant, err := f.ApiClient.GetAIHubTenant(context.Background())
		if err != nil {
			return err
		}

		if len(tenant.Keys) == 0 {
			return fmt.Errorf("no API keys found")
		}

		options := make([]string, 0, len(tenant.Keys))
		for _, key := range tenant.Keys {
			label := key.KeyID
			if key.Alias != "" {
				label = fmt.Sprintf("%s (%s)", key.Alias, key.KeyID)
			}
			options = append(options, label)
		}

		selected, err := f.Prompter.Select("Select a key to delete:", "", options)
		if err != nil {
			return err
		}

		opts.keyID = tenant.Keys[selected].KeyID
	}

	confirmed, err := f.Prompter.Confirm("Are you sure you want to delete this key?", false)
	if err != nil {
		return err
	}
	if !confirmed {
		return nil
	}

	return deleteKey(f, opts)
}

func runDeleteNonInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.keyID == "" {
		return fmt.Errorf("--key-id is required")
	}

	return deleteKey(f, opts)
}

func deleteKey(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Deleting AI Hub key..."),
	)
	s.Start()
	err := f.ApiClient.DeleteAIHubKey(context.Background(), opts.keyID)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{"status": "success", "keyID": opts.keyID})
	}

	f.Log.Infof("Key %s deleted successfully", opts.keyID)

	return nil
}
