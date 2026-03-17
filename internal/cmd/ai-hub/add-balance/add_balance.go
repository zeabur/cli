package addbalance

import (
	"context"
	"fmt"
	"strconv"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	amount   int
	provider string
}

func NewCmdAddBalance(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "add-balance",
		Short: "Add balance to AI Hub",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runAddBalance(f, opts)
		},
	}

	cmd.Flags().IntVar(&opts.amount, "amount", 0, "Amount to add in dollars")
	cmd.Flags().StringVar(&opts.provider, "provider", "litellm", "Payment provider")

	return cmd
}

func runAddBalance(f *cmdutil.Factory, opts Options) error {
	if opts.amount == 0 && f.Interactive {
		input, err := f.Prompter.Input("Amount to add (in dollars): ", "")
		if err != nil {
			return err
		}
		amount, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("invalid amount: %w", err)
		}
		opts.amount = amount
	}

	if opts.amount <= 0 {
		return fmt.Errorf("--amount is required and must be greater than 0")
	}

	amountMillicents := opts.amount * 100000

	var providerPtr *string
	if opts.provider != "" {
		providerPtr = &opts.provider
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Adding balance..."),
	)
	s.Start()
	result, err := f.ApiClient.AddAIHubBalance(context.Background(), amountMillicents, providerPtr)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(result)
	}

	newBalanceDollars := float64(result.NewBalance) / 100000.0
	f.Log.Infof("Balance added successfully! New balance: $%.2f", newBalanceDollars)

	return nil
}
