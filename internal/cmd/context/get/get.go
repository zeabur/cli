package get

import (
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct{}

func NewCmdGet(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get Contexts",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGet(f, opts)
		},
	}

	return cmd
}

func runGet(f *cmdutil.Factory, opts *Options) error {
	// Use the effective context so `--workspace` override truthfully shows
	// "(not set)" for inner context — displaying the persisted team-A
	// project under an override to team-B would mislead the user into
	// thinking it's available there (PLA-1590 B+).
	ctx := f.EffectiveContext()
	project := ctx.GetProject()
	environment := ctx.GetEnvironment()
	service := ctx.GetService()

	header := []string{"Context", "Name", "ID"}
	data := [][]string{
		{"Project", project.GetName(), project.GetID()},
		{"Environment", environment.GetName(), environment.GetID()},
		{"Service", service.GetName(), service.GetID()},
	}

	for _, line := range data {
		if line[1] == "" {
			line[1] = "<not set>"
		}
		if line[2] == "" {
			line[2] = "<not set>"
		}
	}

	if f.JSON {
		out := make([]map[string]string, len(data))
		for i, row := range data {
			out[i] = map[string]string{}
			for j, h := range header {
				out[i][h] = row[j]
			}
		}
		return f.Printer.JSON(out)
	}

	f.Printer.Table(header, data)

	// Human-readable mode also tells the user *why* everything is unset when
	// they're running under a `--workspace` override, so they don't
	// misread it as a config bug. JSON mode stays structurally clean
	// (no prose mixed into the payload) so scripts keep parsing it.
	if f.HasWorkspaceOverride() {
		f.Log.Info("Note: --workspace is one-shot; persisted project/service/environment context is not used.")
	}

	return nil
}
