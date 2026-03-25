package get

import (
	"context"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get <id>",
		Short: "Get details of an API key",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, args[0])
		},
	}
	return cmd
}

func runGet(f *cmdutil.Factory, id string) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching API key..."),
	)
	s.Start()
	key, err := f.ApiClient.GetZSendAPIKey(context.Background(), id)
	s.Stop()
	if err != nil {
		return err
	}

	if f.JSON {
		return f.Printer.JSON(key)
	}

	domains := ""
	if len(key.Domains) > 0 {
		for i, d := range key.Domains {
			if i > 0 {
				domains += ", "
			}
			domains += d
		}
	}

	f.Printer.Table(
		[]string{"Field", "Value"},
		[][]string{
			{"ID", key.ID},
			{"Name", key.Name},
			{"Permission", key.Permission},
			{"Domains", domains},
			{"Created At", key.CreatedAt.String()},
		},
	)
	return nil
}
