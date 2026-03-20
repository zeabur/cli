package update

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	id           string
	firstName    string
	lastName     string
	email        string
	phone        string
	address1     string
	address2     string
	city         string
	state        string
	country      string
	postalCode   string
	organization string
}

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update a registrant profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registrant profile ID")
	cmd.Flags().StringVar(&opts.firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&opts.email, "email", "", "Email address")
	cmd.Flags().StringVar(&opts.phone, "phone", "", "Phone number")
	cmd.Flags().StringVar(&opts.address1, "address1", "", "Address line 1")
	cmd.Flags().StringVar(&opts.address2, "address2", "", "Address line 2")
	cmd.Flags().StringVar(&opts.city, "city", "", "City")
	cmd.Flags().StringVar(&opts.state, "state", "", "State/Province")
	cmd.Flags().StringVar(&opts.country, "country", "", "Country")
	cmd.Flags().StringVar(&opts.postalCode, "postal-code", "", "Postal code")
	cmd.Flags().StringVar(&opts.organization, "organization", "", "Organization")

	return cmd
}

func runUpdate(f *cmdutil.Factory, opts *Options) error {
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
		idx, err := f.Prompter.Select("Select profile to update", "", options)
		if err != nil {
			return err
		}
		opts.id = profiles[idx].ID
	}

	input := model.UpdateRegistrantProfileInput{}
	if opts.firstName != "" {
		input.FirstName = &opts.firstName
	}
	if opts.lastName != "" {
		input.LastName = &opts.lastName
	}
	if opts.email != "" {
		input.Email = &opts.email
	}
	if opts.phone != "" {
		input.Phone = &opts.phone
	}
	if opts.address1 != "" {
		input.Address1 = &opts.address1
	}
	if opts.address2 != "" {
		input.Address2 = &opts.address2
	}
	if opts.city != "" {
		input.City = &opts.city
	}
	if opts.state != "" {
		input.State = &opts.state
	}
	if opts.country != "" {
		input.Country = &opts.country
	}
	if opts.postalCode != "" {
		input.PostalCode = &opts.postalCode
	}
	if opts.organization != "" {
		input.Organization = &opts.organization
	}

	profile, err := f.ApiClient.UpdateRegistrantProfile(ctx, opts.id, input)
	if err != nil {
		return fmt.Errorf("update registrant profile failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(profile)
	}

	f.Log.Infof("Registrant profile updated: %s %s <%s>", profile.FirstName, profile.LastName, profile.Email)
	return nil
}
