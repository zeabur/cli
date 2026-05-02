package searchrepo

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
)

// NewCmdSearchRepo creates the service search-repo command.
func NewCmdSearchRepo(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search-repo [keyword]",
		Short: "Search Git repositories",
		Args:  cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var keyword string
			if len(args) > 0 {
				keyword = args[0]
			}
			return runSearchRepo(cmd, f, keyword)
		},
	}

	return cmd
}

func runSearchRepo(cmd *cobra.Command, f *cmdutil.Factory, keyword string) error {
	if keyword == "" {
		if !f.Interactive {
			return fmt.Errorf("keyword is required")
		}
		k, err := f.Prompter.Input("Enter search keyword:", "")
		if err != nil {
			return err
		}
		keyword = strings.TrimSpace(k)
		if keyword == "" {
			return fmt.Errorf("keyword is required")
		}
	}

	repos, err := f.ApiClient.SearchGitRepositories(cmd.Context(), &keyword)
	if err != nil {
		return fmt.Errorf("search repositories failed: %w", err)
	}

	if len(repos) == 0 {
		if f.JSON {
			return f.Printer.JSON([]any{})
		}
		f.Log.Infof("No repositories found for %q", keyword)
		return nil
	}

	if f.JSON {
		return f.Printer.JSON(repos)
	}

	for _, repo := range repos {
		fmt.Printf("%s/%s (ID: %d)\n", repo.Owner, repo.Name, repo.ID)
	}

	return nil
}
