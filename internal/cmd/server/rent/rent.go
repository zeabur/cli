package rent

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
	provider string
	region   string
	plan     string

	skipConfirm bool
}

func NewCmdRent(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "rent",
		Short: "Rent a new dedicated server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRent(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.provider, "provider", "", "Server provider code")
	cmd.Flags().StringVar(&opts.region, "region", "", "Server region ID")
	cmd.Flags().StringVar(&opts.plan, "plan", "", "Server plan name")
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runRent(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runRentInteractive(f, opts)
	}
	return runRentNonInteractive(f, opts)
}

func runRentInteractive(f *cmdutil.Factory, opts *Options) error {
	ctx := context.Background()

	if opts.provider == "" {
		providers, err := f.ApiClient.ListDedicatedServerProviders(ctx)
		if err != nil {
			return fmt.Errorf("list providers failed: %w", err)
		}
		if len(providers) == 0 {
			return fmt.Errorf("no providers available")
		}

		options := make([]string, len(providers))
		for i, p := range providers {
			options[i] = p.Name
		}

		idx, err := f.Prompter.Select("Select a provider", "", options)
		if err != nil {
			return err
		}
		opts.provider = providers[idx].Code
	}

	if opts.region == "" {
		regions, err := f.ApiClient.ListDedicatedServerRegions(ctx, opts.provider)
		if err != nil {
			return fmt.Errorf("list regions failed: %w", err)
		}
		if len(regions) == 0 {
			return fmt.Errorf("no regions available for provider %s", opts.provider)
		}

		options := make([]string, len(regions))
		for i, r := range regions {
			if r.City != "" && r.Country != "" {
				options[i] = fmt.Sprintf("%s (%s, %s)", r.Name, r.City, r.Country)
			} else if r.Country != "" {
				options[i] = fmt.Sprintf("%s (%s)", r.Name, r.Country)
			} else {
				options[i] = r.Name
			}
		}

		idx, err := f.Prompter.Select("Select a region", "", options)
		if err != nil {
			return err
		}
		opts.region = regions[idx].ID
	}

	if opts.plan == "" {
		plans, err := f.ApiClient.ListDedicatedServerPlans(ctx, opts.provider, opts.region)
		if err != nil {
			return fmt.Errorf("list plans failed: %w", err)
		}
		if len(plans) == 0 {
			return fmt.Errorf("no plans available for provider %s in region %s", opts.provider, opts.region)
		}

		options := make([]string, len(plans))
		for i, p := range plans {
			available := ""
			if !p.Available {
				available = " [unavailable]"
			}
			gpu := ""
			if p.GPU != nil {
				gpu = fmt.Sprintf(", GPU: %s", *p.GPU)
			}
			options[i] = fmt.Sprintf("%s - %d CPU, %d GB RAM, %d GB Disk%s - $%d/mo%s",
				p.Name, p.CPU, p.Memory, p.Disk, gpu, p.Price, available)
		}

		idx, err := f.Prompter.Select("Select a plan", "", options)
		if err != nil {
			return err
		}
		if !plans[idx].Available {
			return fmt.Errorf("plan %s is currently unavailable", plans[idx].Name)
		}
		opts.plan = plans[idx].Name
	}

	return runRentNonInteractive(f, opts)
}

func runRentNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.provider == "" {
		return fmt.Errorf("--provider is required")
	}
	if opts.region == "" {
		return fmt.Errorf("--region is required")
	}
	if opts.plan == "" {
		return fmt.Errorf("--plan is required")
	}

	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(
			fmt.Sprintf("Rent server with provider=%s, region=%s, plan=%s?", opts.provider, opts.region, opts.plan),
			false,
		)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	serverID, err := f.ApiClient.RentServer(context.Background(), opts.provider, opts.region, opts.plan)
	if err != nil {
		errMsg := err.Error()
		if strings.Contains(errMsg, "bind a credit card or recharge credits") || strings.Contains(errMsg, "insufficient balance") {
			f.Log.Errorf("Rent server failed: %s", errMsg)
			f.Log.Infof("Please bind a credit card or top up your balance at: https://zeabur.com/account/billing")
			return fmt.Errorf("payment required")
		}
		return fmt.Errorf("rent server failed: %w", err)
	}

	f.Log.Infof("Server rented successfully, ID: %s", serverID)
	f.Log.Infof("Waiting for server provisioning...")

	return waitForServerInit(f, serverID)
}

func waitForServerInit(f *cmdutil.Factory, serverID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Initializing server..."),
	)
	s.Start()
	defer s.Stop()

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	seenEvents := 0
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server provisioning timed out after 10 minutes, use `zeabur server get %s` to check status", serverID)
		case <-ticker.C:
		}

		server, err := f.ApiClient.GetServer(ctx, serverID)
		if err != nil {
			continue
		}

		for i := seenEvents; i < len(server.Events); i++ {
			ev := server.Events[i]
			s.Stop()
			f.Log.Infof("[%s] %s", ev.Severity, ev.Message)
			s.Start()

			if ev.Message == "Server initialized" {
				s.Stop()
				f.Log.Infof("Server is ready!")
				return nil
			}

			if strings.HasPrefix(ev.Message, "Failed to rent server") {
				s.Stop()
				f.Log.Errorf("Server provisioning failed: %s", ev.Message)
				f.Log.Infof("The charge has been refunded to your balance.")
				return fmt.Errorf("server provisioning failed")
			}
		}
		seenEvents = len(server.Events)

		s.Suffix = fmt.Sprintf(" Initializing server... (%d events)", seenEvents)
	}
}
