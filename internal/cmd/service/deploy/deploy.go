package deploy

import (
	"context"
	"fmt"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	projectID       string
	template        string
	marketplaceCode string
	branchName      string
	name            string
	keyword         string
	repoID          int
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

	cmd.Flags().StringVar(&opts.name, "name", "", "Service Name")
	cmd.Flags().StringVar(&opts.template, "template", "", "Service template")
	cmd.Flags().StringVar(&opts.marketplaceCode, "marketplace-code", "", "Marketplace item code")
	cmd.Flags().IntVar(&opts.repoID, "repo-id", 0, "Git repository ID")
	cmd.Flags().StringVar(&opts.branchName, "branch-name", "", "Git branch name")
	cmd.Flags().StringVar(&opts.keyword, "keyword", "", "Git repository keyword")

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

	template := strings.ToUpper(opts.template)
	switch template {
	case "PREBUILT":
		opts.name = opts.marketplaceCode
		service, err := f.ApiClient.CreatePrebuiltService(ctx, opts.projectID, opts.marketplaceCode)
		if err != nil {
			return fmt.Errorf("create prebuilt service failed: %w", err)
		}

		f.Log.Infof("Service %s created", service.Name)
		return nil
	case "GIT":
		_, err = f.ApiClient.CreateService(context.Background(), f.Config.GetContext().GetProject().GetID(), opts.name, opts.repoID, opts.branchName)
		if err != nil {
			return fmt.Errorf("create service failed: %w", err)
		}

		f.Log.Infof("Service %s created", opts.name)
		return nil
	default:
		return fmt.Errorf("unsupported service template %s", opts.template)
	}
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	// fill project id if not set by asking user
	if _, err := f.ParamFiller.Project(&opts.projectID); err != nil {
		return err
	}

	if opts.template == "" {
		options := []string{"PREBUILT", "GIT"}
		serviceTemplate, err := f.Prompter.Select("Select service template", options[0], options)
		if err != nil {
			return err
		}
		opts.template = options[serviceTemplate]
	}

	ctx := context.Background()

	switch strings.ToUpper(opts.template) {

	case "PREBUILT":

		if opts.marketplaceCode == "" {

			s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
				spinner.WithColor(cmdutil.SpinnerColor),
				spinner.WithSuffix(" Fetching prebuilt marketplace..."),
				spinner.WithFinalMSG(cmdutil.SuccessIcon+" Prebuilt marketplace fetched 🌇\n"),
			)
			s.Start()
			prebuiltItems, err := f.ApiClient.GetPrebuiltItems(ctx)
			if err != nil {
				return fmt.Errorf("get prebuilt marketplace failed: %w", err)
			}
			s.Stop()

			prebuiltItemsList := make([]string, len(prebuiltItems))
			for i, item := range prebuiltItems {
				prebuiltItemsList[i] = item.Name + " (" + item.Description + ")"
			}

			index, err := f.Prompter.Select("Select prebuilt item", "", prebuiltItemsList)
			if err != nil {
				return fmt.Errorf("select prebuilt item failed: %w", err)
			}

			opts.marketplaceCode = prebuiltItems[index].ID
		}

		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Creating service..."),
		)
		s.Start()

		service, err := f.ApiClient.CreatePrebuiltService(ctx, opts.projectID, opts.marketplaceCode)
		if err != nil {
			return fmt.Errorf("create prebuilt service failed: %w", err)
		}

		s.Stop()

		fmt.Printf("%s Service %s created 🚀\n", cmdutil.SuccessIcon, service.Name)
		fmt.Printf("https://dash.zeabur.com/projects/%s/services/%s", opts.projectID, service.ID)

		return nil

	case "GIT":
		var s *spinner.Spinner

		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching Git Repositories..."),
			spinner.WithFinalMSG(cmdutil.SuccessIcon+" Repositories fetched 🌇\n"),
		)
		s.Start()
		gitRepositories, err := f.ApiClient.SearchGitRepositories(ctx, &opts.keyword)
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
		opts.name = gitRepositories[index].Name

		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching Git Repository Branches..."),
		)
		s.Start()
		branches, err := f.ApiClient.GetRepoBranchesByRepoID(opts.repoID)
		if err != nil {
			return fmt.Errorf("get git repository branches failed: %w", err)
		}
		s.Stop()

		if len(branches) == 1 {
			opts.branchName = branches[0]
		} else {
			_, err = f.Prompter.Select("Select branch", opts.branchName, branches)
			if err != nil {
				return err
			}
		}

		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Creating service..."),
			spinner.WithFinalMSG(cmdutil.SuccessIcon+" Service created 🥂\n"),
		)
		s.Start()

		_, err = f.ApiClient.CreateService(context.Background(), f.Config.GetContext().GetProject().GetID(), opts.name, opts.repoID, opts.branchName)
		if err != nil {
			return err
		}
		s.Stop()

		return nil

	default:
		return fmt.Errorf("unsupported service template %s", opts.template)
	}
}

func paramCheck(opts *Options) error {
	if opts.template == "" {
		return fmt.Errorf("please specify service template with --template")
	}

	if strings.ToUpper(opts.template) != "PREBUILT" && strings.ToUpper(opts.template) != "GIT" {
		return fmt.Errorf("unsupported service template %s, only support PREBUILT and GIT", opts.template)
	}

	if opts.template == "PREBUILT" && opts.marketplaceCode == "" {
		return fmt.Errorf("please specify marketplace item code with --marketplace-code")
	}

	if opts.template == "GIT" && opts.repoID == 0 {
		return fmt.Errorf("please specify git repository ID with --repo-id")
	}

	if opts.template == "GIT" && opts.branchName == "" {
		return fmt.Errorf("please specify git branch name with --branch-name")
	}

	return nil
}
