package search

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
)

type Options struct {
	keyword string
}

func NewCmdSearch(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "search [keyword]",
		Short: "Search templates by keyword",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.keyword = args[0]
			}
			return runSearch(f, opts)
		},
	}

	return cmd
}

func runSearch(f *cmdutil.Factory, opts Options) error {
	if opts.keyword == "" {
		if f.Interactive {
			keyword, err := f.Prompter.Input("Search keyword: ", "")
			if err != nil {
				return err
			}
			opts.keyword = keyword
		} else {
			return fmt.Errorf("keyword is required")
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Searching templates..."),
	)
	s.Start()
	allTemplates, err := f.ApiClient.ListAllTemplates(context.Background())
	if err != nil {
		s.Stop()
		return err
	}
	s.Stop()

	keyword := strings.ToLower(opts.keyword)
	var matched model.Templates
	for _, t := range allTemplates {
		name := strings.ToLower(t.Name)
		desc := strings.ToLower(t.Description)
		if strings.Contains(name, keyword) || strings.Contains(desc, keyword) {
			matched = append(matched, t)
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		return matched[i].DeploymentCnt > matched[j].DeploymentCnt
	})

	if len(matched) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	header := []string{"Code", "Name", "Description", "Deployments"}
	rows := make([][]string, 0, len(matched))
	for _, t := range matched {
		rows = append(rows, []string{t.Code, t.Name, t.Description, strconv.Itoa(t.DeploymentCnt)})
	}
	f.Printer.Table(header, rows)

	return nil
}
