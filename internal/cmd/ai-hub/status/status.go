package status

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

func NewCmdStatus(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "status",
		Short: "Show AI Hub tenant status",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStatus(f, opts)
		},
	}

	return cmd
}

func runStatus(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching AI Hub status..."),
	)
	s.Start()
	tenant, err := f.ApiClient.GetAIHubTenant(context.Background())
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(tenant)
	}

	balanceDollars := float64(tenant.Balance) / 100000.0

	f.Log.Infof("Balance: $%.2f", balanceDollars)
	f.Log.Infof("Keys: %d", len(tenant.Keys))
	f.Log.Infof("Provider: %s", tenant.Provider)

	if tenant.AutoRechargeThreshold > 0 || tenant.AutoRechargeAmount > 0 {
		f.Log.Infof("Auto-Recharge: threshold $%.2f, amount $%.2f",
			float64(tenant.AutoRechargeThreshold)/100000.0,
			float64(tenant.AutoRechargeAmount)/100000.0,
		)
	} else {
		f.Log.Infof("Auto-Recharge: disabled")
	}

	return nil
}
