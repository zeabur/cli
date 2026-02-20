package create

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
	"github.com/zeabur/cli/pkg/fill"
)

type Options struct {
	id            string
	name          string
	domainName    string
	environmentID string
	skipConfirm   bool
	IsGenerated   bool
	RedirectTo    string
}

func NewCmdCreateDomain(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "create",
		Short: "create a domain",
		Long:  `Create a domain for a service`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runCreateDomain(f, opts)
		},
	}

	util.AddServiceParam(cmd, &opts.id, &opts.name)
	util.AddEnvOfServiceParam(cmd, &opts.environmentID)
	cmd.Flags().BoolVarP(&opts.skipConfirm, "yes", "y", false, "Skip confirmation")
	cmd.Flags().StringVar(&opts.domainName, "domain", "", "Domain name")
	cmd.Flags().BoolVarP(&opts.IsGenerated, "generated", "g", false, "Is this a generated domain")
	cmd.Flags().StringVar(&opts.RedirectTo, "redirect", "", "Redirect to an existing domain")

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

	if _, err := f.ParamFiller.ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions{
		ProjectCtx:    zctx,
		ServiceID:     &opts.id,
		ServiceName:   &opts.name,
		EnvironmentID: &opts.environmentID,
		CreateNew:     false,
	}); err != nil {
		return err
	}

	isGeneratedSelect, err := f.Prompter.Select("Is this a generated domain?", "Yes", []string{"Yes", "No"})
	if err != nil {
		return err
	}
	opts.IsGenerated = false
	if isGeneratedSelect == 0 {
		opts.IsGenerated = true
	}

	if opts.IsGenerated {
		subDomainInput, err := f.Prompter.Input("The subdomain part of zeabur.app: ", "")
		if err != nil {
			return err
		}
		opts.domainName = subDomainInput
	} else {
		domainInput, err := f.Prompter.Input("Custom Domain: ", "")
		if err != nil {
			return err
		}
		opts.domainName = domainInput
	}

	// Get project from the service to check domain availability
	service, err := f.ApiClient.GetService(context.Background(), opts.id, "", "", "")
	if err != nil {
		return fmt.Errorf("get service failed: %w", err)
	}
	projectID := ""
	if service.Project != nil {
		projectID = service.Project.ID
	}

	project, err := f.ApiClient.GetProject(context.Background(), projectID, "", "")
	if err != nil {
		return err
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Checking domain availability ..."),
	)
	s.Start()
	available, _, err := f.ApiClient.CheckDomainAvailable(context.Background(), opts.domainName, opts.IsGenerated, project.Region.ID)
	if err != nil {
		return err
	}
	s.Stop()

	if !available {
		f.Log.Warnf("Domain %s is not available", opts.domainName)
		return nil
	}

	s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching existing domains ..."),
	)
	s.Start()
	existedDomains, err := f.ApiClient.ListDomains(context.Background(), opts.id, opts.environmentID)
	if err != nil {
		return err
	}
	s.Stop()

	if len(existedDomains) != 0 {
		existedDomainNames := make([]string, len(existedDomains))
		for i, existedDomain := range existedDomains {
			existedDomainNames[i] = existedDomain.Domain
		}
		existedDomainNames = append(existedDomainNames, "None")
		redirectDomainSelect, err := f.Prompter.Select("Redirect to", "None", existedDomainNames)
		if err != nil {
			return err
		}
		if existedDomainNames[redirectDomainSelect] != "None" {
			opts.RedirectTo = existedDomainNames[redirectDomainSelect]
		}
	} else {
		f.Log.Infof("No domains found, skipped set redirect domain")
	}

	return runCreateDomainNonInteractive(f, opts)
}

func runCreateDomainNonInteractive(f *cmdutil.Factory, opts *Options) error {
	if opts.id == "" && opts.name != "" {
		service, err := util.GetServiceByName(f.Config, f.ApiClient, opts.name)
		if err != nil {
			return err
		}
		opts.id = service.ID
	}

	if opts.id == "" {
		return fmt.Errorf("--id or --name is required")
	}

	if opts.environmentID == "" {
		envID, err := util.ResolveEnvironmentIDByServiceID(f.ApiClient, opts.id)
		if err != nil {
			return err
		}
		opts.environmentID = envID
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(fmt.Sprintf(" Creating domain %s ...", opts.domainName)),
	)
	s.Start()
	var domain *string
	var err error
	if opts.RedirectTo != "" {
		domain, err = f.ApiClient.AddDomain(context.Background(), opts.id, opts.environmentID, opts.IsGenerated, opts.domainName, opts.RedirectTo)
	} else {
		domain, err = f.ApiClient.AddDomain(context.Background(), opts.id, opts.environmentID, opts.IsGenerated, opts.domainName)
	}

	if err != nil {
		return fmt.Errorf("add domain failed: %w", err)
	}
	s.Stop()

	f.Log.Infof("Domain %s added", *domain)

	return nil
}
