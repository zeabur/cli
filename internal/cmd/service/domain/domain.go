package domain

import (
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id            string
	name          string
	domainName    string
	environmentID string
	skipConfirm   bool
}

func NewCmdCreateDomain(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "create-domain",
		Short: "create a domain",
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateDomain(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")

	return cmd
}

func runCreateDomain(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runCreateDomainInteractive(f, opts)
	} else {
		return runCreateDomainNonInteractive(f, opts)
	}
}

func runCreateDomainInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	return runCreateDomainNonInteractive(f, opts)
}

func runCreateDomainNonInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(zctx, &opts.id, &opts.name, &opts.environmentID); err != nil {
		return err
	}

	return nil
}
