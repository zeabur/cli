package list

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/pkg/fill"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	id            string
	name          string
	environmentID string
}

func NewCmdListDomains(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}
	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:     "list",
		Short:   "list domains",
		Long:    `List domains of a service`,
		Args:    cobra.NoArgs,
		Aliases: []string{"ls"},
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetService(), &opts.id, &opts.name),
			util.DefaultIDByContext(zctx.GetEnvironment(), &opts.environmentID),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runListDomains(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)

	return cmd
}

func runListDomains(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runListDomainsInteractive(f, opts)
	} else {
		return runListDomainsNonInteractive(f, opts)
	}
}

func runListDomainsInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    zctx,
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	return runListDomainsNonInteractive(f, opts)
}

func runListDomainsNonInteractive(f *cmdutil.Factory, opts *Options) error {
	zctx := f.Config.GetContext()

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    zctx,
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Fetching domains of service %s ...", opts.name)),
	)
	s.Start()
	domainList, err := f.ApiClient.ListDomains(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	s.Stop()

	if len(domainList) == 0 {
		f.Log.Infof("No domains found")
		return nil
	}

	f.Printer.Table(domainList.Header(), domainList.Rows())

	return nil
}
