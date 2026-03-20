package delete

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id          string
	skipConfirm bool
}

func NewCmdDelete(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a registrant profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDelete(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registrant profile ID")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runDelete(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if opts.id == "" {
		if !f.Interactive {
			return fmt.Errorf("--id is required")
		}
		profiles, err := f.ApiClient.ListRegistrantProfiles(ctx)
		if err != nil {
			return fmt.Errorf("list registrant profiles failed: %w", err)
		}
		if len(profiles) == 0 {
			return fmt.Errorf("no registrant profiles found")
		}
		options := make([]string, len(profiles))
		for i, p := range profiles {
			options[i] = fmt.Sprintf("%s %s <%s>", p.FirstName, p.LastName, p.Email)
		}
		idx, err := f.Prompter.Select("Select profile to delete", "", options)
		if err != nil {
			return err
		}
		opts.id = profiles[idx].ID
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm("Delete this registrant profile?", false)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	err := f.ApiClient.DeleteRegistrantProfile(ctx, opts.id)
	if err != nil {
		return fmt.Errorf("delete registrant profile failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{"status": "deleted"})
	}

	f.Log.Infof("Registrant profile deleted")
	return nil
}
