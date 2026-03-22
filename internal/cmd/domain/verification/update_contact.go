package verification

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type updateContactOptions struct {
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

func newCmdUpdateContact(f *cmdutil.Factory) *cobra.Command {
	opts := &updateContactOptions{}

	cmd := &cobra.Command{
		Use:   "update-contact",
		Short: "Update the registrant contact info on a domain (changing email triggers new ICANN verification)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateContact(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
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

func runUpdateContact(f *cmdutil.Factory, opts *updateContactOptions) error {
	ctx := context.Background()

	if opts.id == "" {
		if !f.Interactive {
			return fmt.Errorf("--id is required")
		}
		domains, err := f.ApiClient.ListRegisteredDomains(ctx)
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
		idx, err := f.Prompter.Select("Select domain to update registrant contact", "", options)
		if err != nil {
			return err
		}
		opts.id = domains[idx].ID
	}

	var originalEmail string

	if f.Interactive {
		// Pre-fill from current registrant profile if available
		domain, err := f.ApiClient.GetRegisteredDomain(ctx, opts.id)
		if err == nil && domain.RegistrantProfile != nil {
			p := domain.RegistrantProfile
			originalEmail = p.Email
			if opts.firstName == "" {
				opts.firstName, _ = f.Prompter.Input("First name: ", p.FirstName)
			}
			if opts.lastName == "" {
				opts.lastName, _ = f.Prompter.Input("Last name: ", p.LastName)
			}
			if opts.email == "" {
				opts.email, _ = f.Prompter.Input("Email: ", p.Email)
			}
			if opts.phone == "" {
				opts.phone, _ = f.Prompter.Input("Phone (e.g. +1.5551234567): ", p.Phone)
			}
			if opts.address1 == "" {
				opts.address1, _ = f.Prompter.Input("Address: ", p.Address1)
			}
			if opts.address2 == "" {
				opts.address2, _ = f.Prompter.Input("Address line 2 (optional): ", "")
			}
			if opts.city == "" {
				opts.city, _ = f.Prompter.Input("City: ", p.City)
			}
			if opts.state == "" {
				opts.state, _ = f.Prompter.Input("State/Province: ", p.State)
			}
			if opts.country == "" {
				opts.country, _ = f.Prompter.Input("Country (e.g. US): ", p.Country)
			}
			if opts.postalCode == "" {
				opts.postalCode, _ = f.Prompter.Input("Postal code: ", p.PostalCode)
			}
			if opts.organization == "" {
				opts.organization, _ = f.Prompter.Input("Organization (optional): ", p.Organization)
			}
		} else {
			if opts.firstName == "" {
				opts.firstName, _ = f.Prompter.Input("First name: ", "")
			}
			if opts.lastName == "" {
				opts.lastName, _ = f.Prompter.Input("Last name: ", "")
			}
			if opts.email == "" {
				opts.email, _ = f.Prompter.Input("Email: ", "")
			}
			if opts.phone == "" {
				opts.phone, _ = f.Prompter.Input("Phone (e.g. +1.5551234567): ", "")
			}
			if opts.address1 == "" {
				opts.address1, _ = f.Prompter.Input("Address: ", "")
			}
			if opts.address2 == "" {
				opts.address2, _ = f.Prompter.Input("Address line 2 (optional): ", "")
			}
			if opts.city == "" {
				opts.city, _ = f.Prompter.Input("City: ", "")
			}
			if opts.state == "" {
				opts.state, _ = f.Prompter.Input("State/Province: ", "")
			}
			if opts.country == "" {
				opts.country, _ = f.Prompter.Input("Country (e.g. US): ", "")
			}
			if opts.postalCode == "" {
				opts.postalCode, _ = f.Prompter.Input("Postal code: ", "")
			}
			if opts.organization == "" {
				opts.organization, _ = f.Prompter.Input("Organization (optional): ", "")
			}
		}
	}

	// Validate required fields
	for _, check := range []struct{ val, flag string }{
		{opts.firstName, "--first-name"},
		{opts.lastName, "--last-name"},
		{opts.email, "--email"},
		{opts.phone, "--phone"},
		{opts.address1, "--address1"},
		{opts.city, "--city"},
		{opts.state, "--state"},
		{opts.country, "--country"},
		{opts.postalCode, "--postal-code"},
	} {
		if check.val == "" {
			return fmt.Errorf("%s is required", check.flag)
		}
	}

	input := model.UpdateRegistrantContactInput{
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

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Updating registrant contact..."),
	)
	s.Start()
	err := f.ApiClient.UpdateRegistrantContact(ctx, opts.id, input)
	s.Stop()
	if err != nil {
		return fmt.Errorf("update registrant contact failed: %w", err)
	}

	f.Log.Infof("Registrant contact updated successfully")
	if opts.email != originalEmail {
		f.Log.Infof("Note: changing the email triggers a new ICANN verification flow")
	}

	return nil
}
