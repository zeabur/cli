package api

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v54/github"
	"golang.org/x/oauth2"
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

// Get Repo Info
func (c *client) GetRepoInfo() (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = "."
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return "", "", nil
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

// Get Organizations
func (c *client) GetOrganizationList(username string) ([]string, error) {
	client := github.NewClient(nil)

	orgs, _, err := client.Organizations.List(context.Background(), username, nil)
	if err != nil {
		return nil, err
	}

	orgNames := make([]string, 0, len(orgs))
	orgNames = append(orgNames, username)
	for _, org := range orgs {
		orgNames = append(orgNames, *org.Login)
	}

	return orgNames, nil
}

func (c *client) WipeRepo() error {
	// wipe current repo info by init
	cmd := exec.Command("git", "remote", "remove", "origin")
	cmd.Dir = "."
	_, err := cmd.Output()
	if err != nil {
		return err
	}

	return nil
}

// func (c *client) PushRepo(repo *github.Repository) error {

// 	cmd = exec.Command("git", "remote", "add", "origin", *repo.CloneURL)
// 	cmd.Dir = "."
// 	_, err = cmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to add remote: %v", err)
// 	}

// 	cmd = exec.Command("git", "add", ".")
// 	cmd.Dir = "."
// 	_, err = cmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to add files: %v", err)
// 	}

// 	cmd = exec.Command("git", "commit", "-m", "Initial commit with Zeabur CLI")
// 	cmd.Dir = "."
// 	_, err = cmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to commit files: %v", err)
// 	}

// 	cmd = exec.Command("git", "push", "-u", "origin", "master")
// 	cmd.Dir = "."
// 	_, err = cmd.Output()
// 	if err != nil {
// 		return fmt.Errorf("failed to push files: %v", err)
// 	}

// 	return nil
// }

func (c *client) InitRepo(repoOwner string, repoName string) (*github.Repository, error) {
	cmd := exec.Command("git", "init")
	cmd.Dir = "."
	_, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ""},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	repo, _, err := client.Repositories.Create(context.Background(), repoOwner, &github.Repository{
		Name: github.String(repoName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create repository: %v", err)
	}

	for {
		repo, _, err = client.Repositories.Get(context.Background(), repoOwner, repoName)
		if err != nil {
			fmt.Println("failed to get repository: ", err)
		}
		if repo.GetCloneURL() != "" {
			break
		}

		fmt.Println("waiting for repository to be ready...")
		time.Sleep(5 * time.Second)
	}

	return repo, nil
}
