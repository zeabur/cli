package deploy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/google/go-github/github"
	"github.com/spf13/cobra"
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
		Short:   "Deploy a service",
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
	repoOwner, repoName, err := getRepoInfo()
	if err != nil {
		return err
	}

	repoID, err := getRepoID(repoOwner, repoName)
	if err != nil {
		return err
	}

	fmt.Println(repoID)

	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Fetching repository informations..."),
	)
	s.Start()
	repoOwner, repoName, err := getRepoInfo()
	if err != nil {
		return err
	}

	// Use repo name as default service name
	if opts.name == "" {
		opts.name = repoName
	}

	// Get repo ID and branches concurrently
	repoIDChan := make(chan int)
	branchesChan := make(chan []string)
	errChan := make(chan error)

	go func() {
		repoID, err := getRepoID(repoOwner, repoName)
		if err != nil {
			errChan <- err
			return
		}
		repoIDChan <- repoID
	}()

	go func() {
		branches, err := getRepoBranches(repoOwner, repoName)
		if err != nil {
			errChan <- err
			return
		}
		branchesChan <- branches
	}()

	repoID := <-repoIDChan
	branches := <-branchesChan

	// Check for errors
	select {
	case err := <-errChan:
		return err
	default:
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

func getRepoInfo() (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err != nil {
		return "", "", err
	}

	repoURL := strings.TrimSpace(string(out))
	parts := strings.Split(repoURL, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repository URL")
	}

	repoOwner := strings.TrimPrefix(parts[len(parts)-2], "git@github.com:")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")

	return repoOwner, repoName, nil
}

func getRepoID(repoOwner string, repoName string) (int, error) {
	client := github.NewClient(nil)

	repo, _, err := client.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		return 0, err
	}

	return int(*repo.ID), nil
}

// Get repo branches by repo owner and repo name
func getRepoBranches(repoOwner string, repoName string) ([]string, error) {
	client := github.NewClient(nil)

	branches, _, err := client.Repositories.ListBranches(context.Background(), repoOwner, repoName, nil)
	if err != nil {
		return nil, err
	}

	var branchNames []string
	for _, branch := range branches {
		branchNames = append(branchNames, *branch.Name)
	}

	return branchNames, nil
}
