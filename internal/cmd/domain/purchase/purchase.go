package purchase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	domain        string
	registrantID  string
	skipConfirm   bool
}

func NewCmdPurchase(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "purchase [domain]",
		Short: "Purchase a domain",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.domain = args[0]
			}
			return runPurchase(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.registrantID, "registrant-id", "", "Registrant profile ID")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runPurchase(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runPurchaseInteractive(f, opts)
	}
	return runPurchaseNonInteractive(f, opts)
}

func runPurchaseInteractive(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if opts.domain == "" {
		domain, err := f.Prompter.Input("Domain to purchase (e.g. example.com): ", "")
		if err != nil {
			return err
		}
		opts.domain = domain
	}

	// Check availability
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Checking availability..."),
	)
	s.Start()
	result, err := f.ApiClient.CheckDomainRegistrationAvailability(ctx, opts.domain)
	s.Stop()
	if err != nil {
		return fmt.Errorf("check domain availability failed: %w", err)
	}

	if !result.Available {
		return fmt.Errorf("domain %s is not available for registration", opts.domain)
	}

	f.Log.Infof("Domain: %s", result.Domain)
	if result.Price != nil {
		f.Log.Infof("Price: $%.2f/yr", float64(*result.Price)/100)
	}

	if opts.registrantID == "" {
		profiles, err := f.ApiClient.ListRegistrantProfiles(ctx)
		if err != nil {
			return fmt.Errorf("list registrant profiles failed: %w", err)
		}

		if len(profiles) == 0 {
			f.Log.Infof("No registrant profiles found. Please create one first:")
			f.Log.Infof("  zeabur domain registrant create")
			return fmt.Errorf("no registrant profile available")
		}

		options := make([]string, len(profiles))
		for i, p := range profiles {
			def := ""
			if p.IsDefault {
				def = " (default)"
			}
			options[i] = fmt.Sprintf("%s %s <%s>%s", p.FirstName, p.LastName, p.Email, def)
		}

		idx, err := f.Prompter.Select("Select registrant profile", "", options)
		if err != nil {
			return err
		}
		opts.registrantID = profiles[idx].ID
	}

	return runPurchaseNonInteractive(f, opts)
}

func runPurchaseNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.domain == "" {
		return fmt.Errorf("domain argument is required")
	}
	if opts.registrantID == "" {
		return fmt.Errorf("--registrant-id is required")
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(
			fmt.Sprintf("Purchase %s?", opts.domain),
			false,
		)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Purchasing %s...", opts.domain)),
	)
	s.Start()
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()
	result, err := f.ApiClient.PurchaseDomain(ctx, opts.domain, opts.registrantID)
	s.Stop()
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "bind a credit card") || strings.Contains(errMsg, "insufficient balance") {
			f.Log.Errorf("Purchase failed: %s", errMsg)
			f.Log.Infof("Please bind a credit card or top up your balance at: https://zeabur.com/account/billing")
			return fmt.Errorf("payment required")
		}
		return fmt.Errorf("purchase domain failed: %w", err)
	}

	if f.JSON {
		return f.Printer.JSON(result)
	}

	f.Log.Infof("Domain %s purchased successfully!", result.RegisteredDomain.Domain)
	if result.PaymentAmountFromBalance != nil && *result.PaymentAmountFromBalance > 0 {
		f.Log.Infof("  Paid from balance: $%.2f", float64(*result.PaymentAmountFromBalance)/100)
	}
	if result.PaymentAmountFromPaymentMethod != nil && *result.PaymentAmountFromPaymentMethod > 0 {
		f.Log.Infof("  Paid from card: $%.2f", float64(*result.PaymentAmountFromPaymentMethod)/100)
	}

	return nil
}
