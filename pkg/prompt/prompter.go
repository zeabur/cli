// Package prompt provides a prompter interface for prompting the user for input
// and a survey implementation of that interface.
// based on github.com/cli/cli/internal/prompter/prompter.go
package prompt

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

type prompter struct {
}

func New() Prompter {
	return &prompter{}
}

const defaultPageSize = 10

func (p *prompter) Select(message string, defaultValue string, options []string) (result int, err error) {
	q := &survey.Select{
		Message:  message,
		Options:  options,
		PageSize: defaultPageSize,
	}

	if defaultValue != "" {
		// in some situations, defaultValue ends up not being a valid option; do
		// not set default in that case as it will make survey panic
		for _, o := range options {
			if o == defaultValue {
				q.Default = defaultValue
				break
			}
		}
	}

	err = survey.AskOne(q, &result)

	return
}

func (p *prompter) MultiSelect(message string, defaultValues, options []string) (results []int, err error) {
	q := &survey.MultiSelect{
		Message:  message,
		Options:  options,
		PageSize: defaultPageSize,
	}

	var defaults []string

	if len(defaultValues) > 0 {
		// in some situations, defaultValue ends up not being a valid option; do
		// not set default in that case as it will make survey panic
		for _, o := range options {
			for _, d := range defaultValues {
				if o == d {
					defaults = append(defaults, o)
				}
			}
		}
		if len(defaults) > 0 {
			q.Default = defaults
		}
	}

	err = survey.AskOne(q, &results)

	return
}

func (p *prompter) Input(prompt, defaultValue string) (result string, err error) {
	err = survey.AskOne(&survey.Input{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}

func (p *prompter) Confirm(prompt string, defaultValue bool) (bool, error) {
	res := defaultValue
	confirm := survey.Confirm{
		Message: prompt,
		Default: defaultValue,
	}
	err := survey.AskOne(&confirm, &res)
	if err != nil {
		return false, err
	}
	return res, nil
}

func (p *prompter) ConfirmDeletion(requiredValue string) error {
	var result string

	input := &survey.Input{
		Message: fmt.Sprintf("Type %s to confirm deletion:", requiredValue),
	}

	validator := func(val interface{}) error {
		if str := val.(string); !strings.EqualFold(str, requiredValue) {
			return fmt.Errorf("you entered %s", str)
		}
		return nil
	}

	return survey.AskOne(input, &result, survey.WithValidator(validator))
}

var _ Prompter = (*prompter)(nil)
