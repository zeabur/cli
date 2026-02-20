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

	type scoredTemplate struct {
		template *model.Template
		score    int // higher = more relevant
	}

	var matched []scoredTemplate
	for _, t := range allTemplates {
		name := strings.ToLower(t.Name)
		desc := strings.ToLower(t.Description)

		score := 0
		if strings.Contains(name, keyword) {
			score = 3
		} else if strings.Contains(desc, keyword) {
			score = 2
		} else {
			for _, svc := range t.Services {
				if strings.Contains(strings.ToLower(svc.Name), keyword) {
					score = 1
					break
				}
			}
		}

		if score > 0 {
			matched = append(matched, scoredTemplate{template: t, score: score})
		}
	}

	sort.Slice(matched, func(i, j int) bool {
		if matched[i].score != matched[j].score {
			return matched[i].score > matched[j].score
		}
		return matched[i].template.DeploymentCnt > matched[j].template.DeploymentCnt
	})

	if len(matched) == 0 {
		fmt.Println("No templates found")
		return nil
	}

	header := []string{"Code", "Name", "Description", "Deployments"}
	rows := make([][]string, 0, len(matched))
	for _, m := range matched {
		t := m.template
		rows = append(rows, []string{t.Code, t.Name, t.Description, strconv.Itoa(t.DeploymentCnt)})
	}
	f.Printer.Table(header, rows)

	return nil
}
