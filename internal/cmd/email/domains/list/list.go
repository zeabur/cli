package list

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct{}

func NewCmdList(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List email domains",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(f, opts)
		},
	}

	return cmd
}

func runList(f *cmdutil.Factory, opts Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching domains..."),
	)
	s.Start()
	reply, err := f.ApiClient.ListZSendDomains(context.Background(), nil, nil)
	s.Stop()
	if err != nil {
		return err
	}

	domains := make(model.ZSendDomains, 0, len(reply.Domains))
	for i := range reply.Domains {
		domains = append(domains, &reply.Domains[i])
	}

	if f.JSON {
		if len(domains) == 0 {
			return f.Printer.JSON([]any{})
		}
		return f.Printer.JSON(domains)
	}

	f.Printer.Table(domains.Header(), domains.Rows())

	return nil
}
