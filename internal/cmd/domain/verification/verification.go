package verification

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

func NewCmdVerification(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verification",
		Short: "Manage registrant verification for registered domains",
	}

	cmd.AddCommand(newCmdResend(f))
	cmd.AddCommand(newCmdStatus(f))
	cmd.AddCommand(newCmdUpdateContact(f))

	return cmd
}

func selectDomain(f *cmdutil.Factory, id *string) error {
	ctx := context.Background()
	if *id != "" {
		return nil
	}
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
	idx, err := f.Prompter.Select("Select domain", "", options)
	if err != nil {
		return err
	}
	*id = domains[idx].ID
	return nil
}

// --- resend ---

type resendOptions struct {
	id string
}

func newCmdResend(f *cmdutil.Factory) *cobra.Command {
	opts := &resendOptions{}
	cmd := &cobra.Command{
		Use:   "resend",
		Short: "Resend the ICANN registrant verification email",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runResend(f, opts)
		},
	}
	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
	return cmd
}

func runResend(f *cmdutil.Factory, opts *resendOptions) error {
	if err := selectDomain(f, &opts.id); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Resending verification email..."),
	)
	s.Start()
	err := f.ApiClient.ResendRegistrantVerificationEmail(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return fmt.Errorf("resend verification email failed: %w", err)
	}

	f.Log.Infof("Verification email resent successfully")
	return nil
}

// --- status ---

type statusOptions struct {
	id string
}

func newCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := &statusOptions{}
	cmd := &cobra.Command{
		Use:   "status",
		Short: "Check the ICANN registrant verification status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}
	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
	return cmd
}

func runStatus(f *cmdutil.Factory, opts *statusOptions) error {
	if err := selectDomain(f, &opts.id); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Checking verification status..."),
	)
	s.Start()
	domain, err := f.ApiClient.GetRegisteredDomain(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return fmt.Errorf("get registered domain failed: %w", err)
	}

	status := "unknown"
	if domain.RegistrantVerificationStatus != nil {
		status = *domain.RegistrantVerificationStatus
	}

	if f.JSON {
		return f.Printer.JSON(map[string]string{
			"domain": domain.Domain,
			"status": status,
		})
	}

	f.Log.Infof("Domain: %s", domain.Domain)
	f.Log.Infof("Verification Status: %s", status)
	return nil
}

// --- update-contact ---

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
		Short: "Update the registrant contact info on a registered domain",
		Long:  "Update the registrant contact info at OpenSRS. Changing the email triggers a new ICANN verification flow.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdateContact(f, opts)
		},
	}
	cmd.Flags().StringVar(&opts.id, "id", "", "Registered domain ID")
	cmd.Flags().StringVar(&opts.firstName, "first-name", "", "First name")
	cmd.Flags().StringVar(&opts.lastName, "last-name", "", "Last name")
	cmd.Flags().StringVar(&opts.email, "email", "", "Email address")
	cmd.Flags().StringVar(&opts.phone, "phone", "", "Phone number (e.g. +1.5551234567)")
	cmd.Flags().StringVar(&opts.address1, "address1", "", "Address line 1")
	cmd.Flags().StringVar(&opts.address2, "address2", "", "Address line 2")
	cmd.Flags().StringVar(&opts.city, "city", "", "City")
	cmd.Flags().StringVar(&opts.state, "state", "", "State/Province")
	cmd.Flags().StringVar(&opts.country, "country", "", "Country (ISO 3166-1 alpha-2)")
	cmd.Flags().StringVar(&opts.postalCode, "postal-code", "", "Postal code")
	cmd.Flags().StringVar(&opts.organization, "organization", "", "Organization (optional)")
	return cmd
}

func runUpdateContact(f *cmdutil.Factory, opts *updateContactOptions) error {
	if err := selectDomain(f, &opts.id); err != nil {
		return err
	}

	// Validate required fields
	for _, check := range []struct {
		val, name string
	}{
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
			return fmt.Errorf("%s is required", check.name)
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
	err := f.ApiClient.UpdateRegistrantContact(context.Background(), opts.id, input)
	s.Stop()
	if err != nil {
		return fmt.Errorf("update registrant contact failed: %w", err)
	}

	f.Log.Infof("Registrant contact updated successfully")
	if opts.email != "" {
		f.Log.Infof("A new verification email will be sent to %s", opts.email)
	}
	return nil
}
