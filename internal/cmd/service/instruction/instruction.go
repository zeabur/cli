package instruction

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdInstruction(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "instruction",
		Short: "Instruction for prebuiltservice",
		Long:  `Instruction for prebuilt service`,
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstruction(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runInstruction(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()
	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	if err := paramCheck(opts); err != nil {
		return err
	}

	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	instructions, err := f.ApiClient.ServiceInstructions(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}

	for _, instruction := range instructions {
		fmt.Printf("%s: %s\n", instruction.Title, instruction.Content)
	}

	return nil
}

func paramCheck(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("service id or name is required")
	}

	return nil
}
