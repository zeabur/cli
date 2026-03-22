package getregistered

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdGetRegistered(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get-registered",
		Short: "Get registered domain details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" {
		if !f.Interactive {
			return fmt.Errorf("--id is required")
		}
		domains, err := f.ApiClient.ListRegisteredDomains(context.Background())
		if err != nil {
			return fmt.Errorf("list registered domains failed: %w", err)
		}
		if len(domains) == 0 {
			return fmt.Errorf("no registered domains found")
		}

		options := make([]string, len(domains))
		for i, d := range domains {
			options[i] = d.Domain
		}
		idx, err := f.Prompter.Select("Select a domain", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	domain, err := f.ApiClient.GetRegisteredDomain(context.Background(), opts.id)
	if err != nil {
		return fmt.Errorf("get registered domain failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	f.Printer.Table(domain.Header(), domain.Rows())

	if domain.RegistrantProfile != nil {
		p := domain.RegistrantProfile
		f.Log.Infof("")
		f.Log.Infof("Registrant Profile:")
		f.Printer.Table(
			[]string{"Name", "Email", "Phone", "Country"},
			[][]string{{p.FirstName + " " + p.LastName, p.Email, p.Phone, p.Country}},
		)
	}

	return nil
}
