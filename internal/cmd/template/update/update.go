package update

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/util"
)

type Options struct {
	code string
	file string
}

func NewCmdUpdate(f *cmdutil.Factory) *cobra.Command {
	opts := Options{}

	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update template from file",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUpdate(f, opts)
		},
	}

	cmd.Flags().StringVarP(&opts.code, "code", "c", "", "Template code")
	cmd.Flags().StringVarP(&opts.file, "file", "f", "", "Template file")

	return cmd
}

func runUpdate(f *cmdutil.Factory, opts Options) error {
	if opts.code == "" {
		return fmt.Errorf("code is required, use -c or --code to specify the template code")
	}

	if opts.file == "" {
		return fmt.Errorf("file is required, use -f or --file to specify the file path")
	}

	var file []byte
	var err error

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

	if err := util.ValidateTemplate(file); err != nil {
		return fmt.Errorf("validate template: %w", err)
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

	_, err = f.ApiClient.UpdateTemplateFromFile(context.Background(), opts.code, string(file))
	if err != nil {
		return err
	}

	f.Log.Infof("Template updated.")
	return nil
}
