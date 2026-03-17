package autorecharge

import (
	"context"
	"fmt"
	"strconv"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	threshold int
	amount    int
}

func NewCmdAutoRecharge(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "auto-recharge",
		Short: "Configure AI Hub auto-recharge settings",
		Long:  "Configure auto-recharge settings. Set both --threshold and --amount to 0 to disable auto-recharge.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAutoRecharge(f, opts)
		},
	}

	cmd.Flags().IntVar(&opts.threshold, "threshold", -1, "Balance threshold in dollars to trigger recharge (0 to disable)")
	cmd.Flags().IntVar(&opts.amount, "amount", -1, "Amount in dollars to recharge (0 to disable)")

	return cmd
}

func runAutoRecharge(f *cmdutil.Factory, opts Options) error {
	if f.Interactive && (opts.threshold < 0 || opts.amount < 0) {
		return runAutoRechargeInteractive(f, opts)
	}
	return runAutoRechargeNonInteractive(f, opts)
}

func runAutoRechargeInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.threshold < 0 {
		input, err := f.Prompter.Input("Recharge threshold (in dollars, 0 to disable): ", "")
		if err != nil {
			return err
		}
		threshold, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid threshold: %w", err)
		}
		opts.threshold = threshold
	}

	if opts.amount < 0 {
		input, err := f.Prompter.Input("Recharge amount (in dollars, 0 to disable): ", "")
		if err != nil {
			return err
		}
		amount, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid amount: %w", err)
		}
		opts.amount = amount
	}

	return updateAutoRecharge(f, opts)
}

func runAutoRechargeNonInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.threshold < 0 || opts.amount < 0 {
		return fmt.Errorf("both --threshold and --amount are required")
	}

	return updateAutoRecharge(f, opts)
}

func updateAutoRecharge(f *cmdutil.Factory, opts Options) error {
	thresholdMillicents := opts.threshold * 100000
	amountMillicents := opts.amount * 100000

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Updating auto-recharge settings..."),
	)
	s.Start()
	result, err := f.ApiClient.UpdateAIHubAutoRechargeSettings(context.Background(), thresholdMillicents, amountMillicents)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(result)
	}

	if result.AutoRechargeThreshold == 0 && result.AutoRechargeAmount == 0 {
		f.Log.Infof("Auto-recharge disabled")
	} else {
		f.Log.Infof("Auto-recharge updated: threshold $%.2f, amount $%.2f",
			float64(result.AutoRechargeThreshold)/100000.0,
			float64(result.AutoRechargeAmount)/100000.0,
		)
	}

	return nil
}
