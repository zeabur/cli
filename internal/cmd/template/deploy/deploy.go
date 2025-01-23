package deploy

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hasura/go-graphql-client"
	"github.com/spf13/cobra"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/util"
	"gopkg.in/yaml.v3"
)

type Options struct {
	file           string
	projectID      string
	skipValidation bool
}

func NewCmdDeploy(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Validate and deploy a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDeploy(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.file, "file", "f", "", "Template file")
	cmd.Flags().StringVar(&opts.projectID, "project-id", "", "Project ID to deploy on")
	cmd.Flags().BoolVar(&opts.skipValidation, "skip-validation", false, "Skip template validation")

	return cmd
}

func runDeploy(f *cmdutil.Factory, opts *Options) error {
	var err error

	if err := paramCheck(opts); err != nil {
		return err
	}

	var file []byte

	if strings.HasPrefix(opts.file, "https://") || strings.HasPrefix(opts.file, "http://") {
		s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Fetching remote template file ..."),
		)

		s.Start()
		get, err := http.Get(opts.file)
		if err != nil {
			f.Log.Errorf("fetch file failed: %v", err)
			return err
		}

		file, err = io.ReadAll(get.Body)
		if err != nil {
			f.Log.Errorf("read file failed: %v", err)
			return err
		}

		s.Stop()
	} else {
		file, err = os.ReadFile(opts.file)
		if err != nil {
			return fmt.Errorf("read file failed: %w", err)
		}
	}

	if !opts.skipValidation {
		if err := util.ValidateTemplate(file); err != nil {
			return fmt.Errorf("validate template: %w", err)
		}
	}

	if _, err := f.ParamFiller.ProjectCreatePreferred(&opts.projectID); err != nil {
		return err
	}

	type RawTemplate struct {
		APIVersion string `yaml:"apiVersion"`
		Kind       string `yaml:"kind"`
		Metadata   struct {
			Name string `yaml:"name"`
		} `yaml:"metadata"`
		Spec struct {
			Variables []struct {
				Key         string `yaml:"key"`
				Type        string `yaml:"type"`
				Name        string `yaml:"name"`
				Description string `yaml:"description"`
			} `yaml:"variables"`
		}
	}

	var raw RawTemplate
	err = yaml.Unmarshal(file, &raw)
	if err != nil {
		return fmt.Errorf("parse yaml failed: %w", err)
	}

	project, err := f.ApiClient.GetProject(context.Background(), opts.projectID, "", "")
	if err != nil {
		return fmt.Errorf("get project info failed: %w", err)
	}

	vars := model.Map{}
	for _, v := range raw.Spec.Variables {
		switch v.Type {
		case "DOMAIN":
			// Notice: flex shared cluster in China mainland (sha1) does not support generated domain
			if project.Region.ID == "sha1" {
				f.Log.Warnf("Selected region does not support generated domain, please bind a custom domain after template deployed.\n")
				continue
			}

			for {
				val, err := f.Prompter.InputWithHelp(v.Description, "For example, if you enter \"myapp\", the domain will be \"myapp.zeabur.app\"", "")
				if err != nil {
					return err
				}

				s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
					spinner.WithColor(cmdutil.SpinnerColor),
					spinner.WithSuffix(" Checking if domain "+val+".zeabur.app is available ..."),
				)

				s.Start()
				available, _, err := f.ApiClient.CheckDomainAvailable(context.Background(), val, true, project.Region.ID)
				if err != nil {
					return err
				}
				s.Stop()

				if !available {
					f.Log.Warnf("Domain %s.zeabur.app is not available, please try another one.\n", val)
					continue
				}

				f.Log.Infof("Domain %s.zeabur.app is available!\n", val)
				vars[v.Key] = val
				break
			}
		default:
			val, err := f.Prompter.Input(v.Description, "")
			if err != nil {
				return err
			}

			vars[v.Key] = val
		}
	}

	s := spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
		spinner.WithColor(cmdutil.SpinnerColor),
		spinner.WithSuffix(" Deploying template ..."),
	)

	s.Start()
	res, err := f.ApiClient.DeployTemplate(
		context.Background(),
		string(file),
		vars,
		model.RepoConfigs{},
		opts.projectID,
	)
	s.Stop()

	if err != nil {
		var graphqlErrors graphql.Errors
		if errors.As(err, &graphqlErrors) && len(graphqlErrors) > 0 {
			f.Log.Errorf("%s (code: %s)", graphqlErrors[0].Message, graphqlErrors[0].Extensions["code"])
			return nil
		}
		return err
	}

	f.Log.Infof("Template successfully deployed into project %q (https://dash.zeabur.com/projects/%s).", res.Name, res.ID)

	if d, ok := vars["PUBLIC_DOMAIN"]; ok && project.Region.ID != "sha1" {
		s = spinner.New(cmdutil.SpinnerCharSet, cmdutil.SpinnerInterval,
			spinner.WithColor(cmdutil.SpinnerColor),
			spinner.WithSuffix(" Waiting service status ..."),
		)

		s.Start()

		start := time.Now()
		for {
			if time.Since(start) > 2*time.Minute {
				s.Stop()
				return fmt.Errorf("failed to wait service ready, check logs in https://dash.zeabur.com/projects/%s", res.ID)
			}

			time.Sleep(2 * time.Second)
			get, err := http.Get(fmt.Sprintf("https://%s.zeabur.app/", d))
			if err != nil {
				continue
			}

			if get.StatusCode%100 != 5 {
				s.Stop()
				f.Log.Infof("Service ready, you can now visit via https://%s.zeabur.app/", d)
				break
			}

			continue
		}
	}

	return nil
}

func paramCheck(opts *Options) error {
	if opts.file == "" {
		return fmt.Errorf("please specify template file by -f or --file, you can use remote file by http(s)://... or local file path")
	}

	return nil
}
