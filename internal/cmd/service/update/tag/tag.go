package tag

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id   string
	name string

	environmentID string

	tag string

	skipConfirm bool
}

func NewCmdTag(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "tag",
		Short: "Update image tag of a prebuilt service",
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringVarP(&opts.tag, "tag", "t", "latest", "The new tag of the image")

	return cmd
}

func runUpdate(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runInteractive(f, opts)
	} else {
		return runNonInteractive(f, opts)
	}
}

func runInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	varInput, err := f.Prompter.Input("Enter a new image tag", "latest")
	if err != nil {
		return err
	}

	opts.tag = varInput

	return runNonInteractive(f, opts)
}

func runNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if err := checkParams(opts); err != nil {
		return err
	}

	// if name is set, get service id by name
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	idOrName := opts.name
	if idOrName == "" {
		idOrName = opts.id
	}

	// double check
	if f.Interactive && !opts.skipConfirm {
		confirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to update image tag of service <%s>?", idOrName), true)
		if err != nil {
			return err
		}
		if !confirm {
			return nil
		}
	}

	err := f.ApiClient.UpdateImageTag(context.Background(), opts.id, opts.environmentID, opts.tag)
	if err != nil {
		return err
	}

	return nil
}

func checkParams(opts *Options) error {
	if opts.id == "" && opts.name == "" {
		return fmt.Errorf("--id or --name is required")
	}

	if opts.environmentID == "" {
		return fmt.Errorf("--env-id is required")
	}

	if opts.tag == "" {
		return fmt.Errorf("--tag is required")
	}

	return nil
}
