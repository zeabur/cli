package create

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
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

func NewCmdCreate(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a registrant profile",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreate(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&opts.email, "email", "", "Email address")
	cmd.Flags().StringVar(&opts.phone, "phone", "", "Phone number (E.164 format, e.g. +1.5551234567)")
	cmd.Flags().StringVar(&opts.address1, "address1", "", "Address line 1")
	cmd.Flags().StringVar(&opts.address2, "address2", "", "Address line 2")
	cmd.Flags().StringVar(&opts.city, "city", "", "City")
	cmd.Flags().StringVar(&opts.state, "state", "", "State/Province")
	cmd.Flags().StringVar(&opts.country, "country", "", "Country (ISO 3166-1 alpha-2, e.g. US)")
	cmd.Flags().StringVar(&opts.postalCode, "postal-code", "", "Postal code")
	cmd.Flags().StringVar(&opts.organization, "organization", "", "Organization (optional)")

	return cmd
}

func runCreate(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runCreateInteractive(f, opts)
	}
	return runCreateNonInteractive(f, opts)
}

func runCreateInteractive(f *cmdutil.Factory, opts *Options) error {
	var err error

	if opts.firstName == "" {
		opts.firstName, err = f.Prompter.Input("First name: ", "")
		if err != nil {
			return err
		}
	}
	if opts.lastName == "" {
		opts.lastName, err = f.Prompter.Input("Last name: ", "")
		if err != nil {
			return err
		}
	}
	if opts.email == "" {
		opts.email, err = f.Prompter.Input("Email: ", "")
		if err != nil {
			return err
		}
	}
	if opts.phone == "" {
		opts.phone, err = f.Prompter.Input("Phone (e.g. +1.5551234567): ", "")
		if err != nil {
			return err
		}
	}
	if opts.address1 == "" {
		opts.address1, err = f.Prompter.Input("Address: ", "")
		if err != nil {
			return err
		}
	}
	if opts.address2 == "" {
		opts.address2, err = f.Prompter.Input("Address line 2 (optional): ", "")
		if err != nil {
			return err
		}
	}
	if opts.city == "" {
		opts.city, err = f.Prompter.Input("City: ", "")
		if err != nil {
			return err
		}
	}
	if opts.state == "" {
		opts.state, err = f.Prompter.Input("State/Province: ", "")
		if err != nil {
			return err
		}
	}
	if opts.country == "" {
		opts.country, err = f.Prompter.Input("Country (e.g. US): ", "")
		if err != nil {
			return err
		}
	}
	if opts.postalCode == "" {
		opts.postalCode, err = f.Prompter.Input("Postal code: ", "")
		if err != nil {
			return err
		}
	}
	if opts.organization == "" {
		opts.organization, err = f.Prompter.Input("Organization (optional): ", "")
		if err != nil {
			return err
		}
	}

	return runCreateNonInteractive(f, opts)
}

func runCreateNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.firstName == "" {
		return fmt.Errorf("--first-name is required")
	}
	if opts.lastName == "" {
		return fmt.Errorf("--last-name is required")
	}
	if opts.email == "" {
		return fmt.Errorf("--email is required")
	}
	if opts.phone == "" {
		return fmt.Errorf("--phone is required")
	}
	if opts.address1 == "" {
		return fmt.Errorf("--address1 is required")
	}
	if opts.city == "" {
		return fmt.Errorf("--city is required")
	}
	if opts.state == "" {
		return fmt.Errorf("--state is required")
	}
	if opts.country == "" {
		return fmt.Errorf("--country is required")
	}
	if opts.postalCode == "" {
		return fmt.Errorf("--postal-code is required")
	}

	input := model.CreateRegistrantProfileInput{
		FirstName:  opts.firstName,
		LastName:   opts.lastName,
		Email:      opts.email,
		Phone:      opts.phone,
		Address1:   opts.address1,
		City:       opts.city,
		State:      opts.state,
		Country:    opts.country,
		PostalCode: opts.postalCode,
	}
	if opts.address2 != "" {
		input.Address2 = &opts.address2
	}
	if opts.organization != "" {
		input.Organization = &opts.organization
	}

	profile, err := f.ApiClient.CreateRegistrantProfile(context.Background(), input)
	if err != nil {
		return fmt.Errorf("create registrant profile failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(profile)
	}

	f.Log.Infof("Registrant profile created: %s %s <%s>", profile.FirstName, profile.LastName, profile.Email)
	return nil
}
