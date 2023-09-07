package deploy

import (
	"context"
	"fmt"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/internal/util"
)

type Options struct {
	name string
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:     "deploy",
		Short:   "Deploy a local Git Service",
		PreRunE: util.NeedProjectContextWhenNonInteractive(f),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVar(&opts.name, "name", "", "Service name")

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
	repoOwner, repoName, err := f.ApiClient.GetRepoInfo()
	if err != nil {
		return err
	}

	repoID, err := f.ApiClient.GetRepoID(repoOwner, repoName)
	if err != nil {
		return err
	}

	f.Log.Debugf("repoID: %d", repoID)

	// TODO: Deploy Local Git Service NonInteractive

	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	var repoOwner string
	var repoName string
	var err error

	repoOwner, repoName, err = f.ApiClient.GetRepoInfo()
	if err != nil {
		return fmt.Errorf("failed to get repository info: %w", err)
	}

	user, err := f.ApiClient.GetUserInfo(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching repository information..."),
	)
	s.Start()

	if repoOwner != user.Username {
		s.Stop()
		confirm, err := f.Prompter.Confirm("You are not the owner of this repository, would you like to wipe the repository information and continue? (y/n)", true)
		if err != nil {
			return err
		}

		if confirm {
			f.ApiClient.WipeRepo()
			f.Log.Info("Repository information wiped")

			s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
				spinner.WithColor(cmdutil.SpinnerColor),
				spinner.WithSuffix(" Fetching organizations..."),
			)
			s.Start()
			orgs, err := f.ApiClient.GetOrganizationList(user.Username)
			if err != nil {
				return err
			}
			s.Stop()

			index, err := f.Prompter.Select("Select organization", repoOwner, orgs)
			if err != nil {
				return err
			}
			repoOwner = orgs[index]

			repo, err := f.ApiClient.InitRepo(repoName, repoOwner)
			if err != nil {
				return err
			}

			fmt.Println(repo)

			s.Stop()
		} else {
			return nil
		}
	}

	// Use repo name as default service name
	if opts.name == "" {
		opts.name = repoName
	}

	var eg errgroup.Group
	var repoID int
	var branches []string

	eg.Go(func() error {
		repoID, err = f.ApiClient.GetRepoID(repoOwner, repoName)
		return err
	})

	eg.Go(func() error {
		branches, err = f.ApiClient.GetRepoBranches(context.Background(), repoOwner, repoName)
		return err
	})

	if err = eg.Wait(); err != nil {
		return err
	}

	s.Stop()

	// If repo has only one branch, use it as default branch
	// Otherwise, ask user to select a branch
	var branch string

	if len(branches) == 1 {
		branch = branches[0]
	} else {
		_, err = f.Prompter.Select("Select branch", branch, branches)
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

	_, err = f.ApiClient.CreateService(context.Background(), f.Config.GetContext().GetProject().GetID(), opts.name, repoID, branch)
	if err != nil {
		return err
	}
	s.Stop()

	return nil
}
