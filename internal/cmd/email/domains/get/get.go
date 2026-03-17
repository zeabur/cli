package get

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	id string
}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get email domain details and DNS records",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.id, "id", "", "Domain ID")

	return cmd
}

func runGet(f *cmdutil.Factory, opts Options) error {
	if f.Interactive {
		return runGetInteractive(f, opts)
	}
	return runGetNonInteractive(f, opts)
}

func runGetInteractive(f *cmdutil.Factory, opts Options) error {
	if opts.id == "" {
		id, err := f.Prompter.Input("Domain ID: ", "")
		if err != nil {
			return err
		}
		opts.id = id
	}

	if err := paramCheck(opts); err != nil {
		return err
	}
	return getDomain(f, opts)
}

func runGetNonInteractive(f *cmdutil.Factory, opts Options) error {
	if err := paramCheck(opts); err != nil {
		return err
	}
	return getDomain(f, opts)
}

func getDomain(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching domain..."),
	)
	s.Start()
	domain, err := f.ApiClient.GetZSendDomain(context.Background(), opts.id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(domain)
	}

	statusMsg := ""
	if domain.StatusMsg != nil {
		statusMsg = *domain.StatusMsg
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", domain.ID},
			{"Domain", domain.Value},
			{"Region", domain.Region},
			{"Status", domain.Status},
			{"Status Message", statusMsg},
		},
	)

	if len(domain.Records) > 0 {
		f.Log.Infof("\nDNS Records:")
		rows := make([][]string, 0, len(domain.Records))
		for _, r := range domain.Records {
			rows = append(rows, []string{
				r.Category,
				r.Type,
				r.Name,
				r.Content,
				r.Status,
			})
		}
		f.Printer.Table([]string{"Category", "Type", "Name", "Content", "Status"}, rows)
	}

	return nil
}

func paramCheck(opts Options) error {
	if opts.id == "" {
		return fmt.Errorf("domain ID is required")
	}
	return nil
}
