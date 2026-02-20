package delete

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
	domainName    string
	skipConfirm   bool
}

func NewCmdDeleteDomain(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "delete",
		Short:   "Delete domain",
		Long:    `Delete domain of a service`,
		Aliases: []string{"del"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeleteDomain(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringVar(&opts.domainName, "domain", "", "Domain name")

	return cmd
}

func runDeleteDomain(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runDeleteDomainInteractive(f, opts)
	} else {
		return runDeleteDomainNonInteractive(f, opts)
	}
}

func runDeleteDomainInteractive(f *cmdutil.Factory, opts *Options) error {
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

	domainNameList := make([]string, len(domainList))
	for i, domain := range domainList {
		domainNameList[i] = domain.Domain
	}
	deleteDomainSelection, err := f.Prompter.Select("Select domain to delete", "", domainNameList)
	if err != nil {
		return err
	}
	opts.domainName = domainNameList[deleteDomainSelection]

	if opts.skipConfirm {
		return runDeleteDomainNonInteractive(f, opts)
	}

	deleteConfirm, err := f.Prompter.Confirm(fmt.Sprintf("Are you sure to delete domain %s", opts.domainName), false)
	if err != nil {
		return err
	}

	if !deleteConfirm {
		f.Log.Infof("Delete domain %s canceled", opts.domainName)
		return nil
	}

	return runDeleteDomainNonInteractive(f, opts)
}

func runDeleteDomainNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.environmentID == "" && opts.id != "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.id)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Deleting selected domain: %s ...", opts.domainName)))
	s.Start()
	deleteResult, err := f.ApiClient.RemoveDomain(context.Background(), opts.domainName)
	if err != nil {
		return err
	}
	s.Stop()

	if !deleteResult {
		f.Log.Warnf("Delete domain %s failed", opts.domainName)
		return nil
	}

	f.Log.Infof("Delete domain %s success", opts.domainName)

	return nil
}
