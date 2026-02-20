package instruction

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdInstruction(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "instruction",
		Short: "Instruction for prebuiltservice",
		Long:  `Instruction for prebuilt service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runInstruction(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runInstruction(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive && opts.id == "" && opts.name == "" {
		zctx := f.Config.GetContext()
		if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
			ProjectCtx:    zctx,
			ServiceID:     &opts.id,
			ServiceName:   &opts.name,
			EnvironmentID: &opts.environmentID,
			CreateNew:     false,
			FilterFunc: func(service *model.Service) bool {
				return service.Template == "PREBUILT"
			},
		}); err != nil {
			return err
		}
	}

	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("service id or name is required")
	}

	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.id)
		if err != nil {
			return err
		}
		opts.environmentID = envID
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
