package deploy

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	projectID  string
	template   string
	itemCode   string
	branchName string
	name       string
	repoID     int
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	zctx := f.Config.GetContext()

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy a service",
		PreRunE: util.RunEChain(
			util.NeedProjectContextWhenNonInteractive(f),
			util.DefaultIDNameByContext(zctx.GetProject(), &opts.projectID, new(string)),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.template, "template", "", "Service template")
	cmd.Flags().StringVar(&opts.itemCode, "item-code", "", "Marketplace item code")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	if f.Interactive {
		return runDeployInteractive(f, opts)
	} else {
		return runDeployNonInteractive(f, opts)
	}
}

func runDeployNonInteractive(f *cmdutil.Factory, opts *Options) error {
	err := paramCheck(opts)
	if err != nil {
		return err
	}

	ctx := context.Background()

	if opts.template == "MARKETPLACE" {
		opts.name = opts.itemCode
		service, err := f.ApiClient.CreateServiceFromMarketplace(ctx, opts.projectID, opts.name, opts.itemCode)
		if err != nil {
			return fmt.Errorf("create service failed: %w", err)
		}

		f.Log.Infof("Service %s created", service.Name)
	}

	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	// fill project id if not set by asking user
	if _, err := f.ParamFiller.Project(&opts.projectID); err != nil {
		return err
	}

	serviceTemplate, err := f.Prompter.Select("Select service template", "MARKETPLACE", []string{"MARKETPLACE", "GIT"})
	if err != nil {
		return err
	}

	ctx := context.Background()

	if serviceTemplate == 0 {
		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching marketplace items..."),
			spinner.WithFinalMSG(cmdutil.SuccessIcon+" Marketplace fetched 🌇\n"),
		)
		s.Start()
		marketplaceItems, err := f.ApiClient.GetMarketplaceItems(ctx)
		if err != nil {
			return fmt.Errorf("get marketplace items failed: %w", err)
		}
		s.Stop()

		marketplaceItemsList := make([]string, len(marketplaceItems))
		for i, item := range marketplaceItems {
			marketplaceItemsList[i] = item.Name + " (" + item.Description + ")"
		}

		index, err := f.Prompter.Select("Select marketplace item", "", marketplaceItemsList)
		if err != nil {
			return fmt.Errorf("select marketplace item failed: %w", err)
		}

		opts.itemCode = marketplaceItems[index].Code
		opts.name = opts.itemCode

		// use a closure to get the service name after creation
		serviceName := ""
		getServiceName := func() string {
			return serviceName
		}

		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Creating service..."),
		)
		// use a closure to update the spinner's final message(especially the service name)
		s.PreUpdate = func(s *spinner.Spinner) {
			s.FinalMSG = fmt.Sprintf("%s Service %s created 🚀\n", cmdutil.SuccessIcon, getServiceName())
		}
		s.Start()

		service, err := f.ApiClient.CreateServiceFromMarketplace(ctx, opts.projectID, opts.name, opts.itemCode)
		if err != nil {
			return fmt.Errorf("create service failed: %w", err)
		}
		serviceName = service.Name

		s.Stop()
	} else if serviceTemplate == 1 {
		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching Git Repositories..."),
			spinner.WithFinalMSG(cmdutil.SuccessIcon+" Repositories fetched 🌇\n"),
		)
		s.Start()
		gitRepositories, err := f.ApiClient.SearchGitRepositories(ctx, nil)
		if err != nil {
			return fmt.Errorf("search git repositories failed: %w", err)
		}
		s.Stop()

		gitRepositoriesList := make([]string, len(gitRepositories))
		for i, repo := range gitRepositories {
			gitRepositoriesList[i] = repo.Owner + "/" + repo.Name
		}

		index, err := f.Prompter.Select("Select git repository", "", gitRepositoriesList)
		if err != nil {
			return fmt.Errorf("select git repository failed: %w", err)
		}

		opts.repoID = gitRepositories[index].ID

	}

	return nil
}

func paramCheck(opts *Options) error {
	if opts.template == "" {
		return fmt.Errorf("please specify service template with --template")
	}

	if opts.template == "MARKETPLACE" && opts.itemCode == "" {
		return fmt.Errorf("please specify marketplace item code with --item-code")
	}

	return nil
}
