package api

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/google/go-github/github"
)

func (c *client) GetRepoBranches(ctx context.Context, repoOwner string, repoName string) ([]string, error) {
	client := github.NewClient(nil)

	branches, _, err := client.Repositories.ListBranches(context.Background(), repoOwner, repoName, nil)
	if err != nil {
		return nil, err
	}

	branchNames := make([]string, 0, len(branches))
	for _, branch := range branches {
		branchNames = append(branchNames, *branch.Name)
	}

	return branchNames, nil
}

func (c *client) GetRepoID(repoOwner string, repoName string) (int, error) {
	//TODO: Deal with GitHub Auth, reading token env and set HTTP client header
	client := github.NewClient(nil)

	repo, _, err := client.Repositories.Get(context.Background(), repoOwner, repoName)
	if err != nil {
		return 0, err
	}

	return int(*repo.ID), nil
}

func (c *client) GetRepoInfo() (string, string, error) {
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

func (c *client) GetRepoBranchesByRepoID(repoID int) ([]string, error) {
	var query struct {
		GitRepoBranches []string `graphql:"gitRepoBranches(repoID: $repoID)"`
	}

	err := c.Query(context.Background(), &query, V{
		"repoID": repoID,
	})
	if err != nil {
		return nil, err
	}

	return query.GitRepoBranches, nil
}
