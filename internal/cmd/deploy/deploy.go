package deploy

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

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
		PreRunE: util.NeedProjectContext(f),
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
	repoID, err := getRepoID()
	if err != nil {
		return err
	}

	fmt.Println(repoID)

	return nil
}

func runDeployInteractive(f *cmdutil.Factory, opts *Options) error {
	s := spinner.New(spinner.CharSets[1], 100*time.Millisecond)
	s.Start()
	repoID, err := getRepoID()
	if err != nil {
		return err
	}

	service, err := f.ApiClient.CreateService(context.Background(), f.Config.GetContext().GetProject().GetID(), opts.name, repoID, "")
	if err != nil {
		return err
	}

	s.Stop()

	f.Log.Infof("Service %s created", service.Name)

	return nil
}

func getRepoID() (int, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	repoURL := strings.TrimSpace(string(out))
	parts := strings.Split(repoURL, "/")
	if len(parts) < 2 {
		return 0, fmt.Errorf("invalid repository URL")
	}

	repoOwner := strings.TrimPrefix(parts[len(parts)-2], "git@github.com:")
	repoName := strings.TrimSuffix(parts[len(parts)-1], ".git")

	client := github.NewClient(nil)

	repo, _, err := client.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		return 0, err
	}

	return int(*repo.ID), nil
}
