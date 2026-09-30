// Package prompt provides a prompter interface for prompting the user for input
// and a survey implementation of that interface.
// based on github.com/cli/cli/internal/prompter/prompter.go
package prompt

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"golang.org/x/term"
)

// ErrNonInteractive is returned when a prompt is needed but stdin is not a
// terminal. survey blocks forever on an open pipe (CI, agent shells, cron),
// so fail fast instead.
var ErrNonInteractive = errors.New("stdin is not a terminal")

func ensureTerminal(message string) error {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		return nil
	}
	return fmt.Errorf("cannot prompt %q: %w; pass the required flags (see --help) or use -i=false", message, ErrNonInteractive)
}

type prompter struct{}

// New creates a new prompter
func New() Prompter {
	return &prompter{}
}

const defaultPageSize = 10

func (p *prompter) Select(message string, defaultValue string, options []string) (result int, err error) {
	if err := ensureTerminal(message); err != nil {
		return 0, err
	}

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
	if err := ensureTerminal(message); err != nil {
		return nil, err
	}

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
	if err := ensureTerminal(prompt); err != nil {
		return "", err
	}

	err = survey.AskOne(&survey.Input{
		Message: prompt,
		Default: defaultValue,
	}, &result)

	return
}

func (p *prompter) InputWithHelp(prompt, defaultValue, help string) (result string, err error) {
	if err := ensureTerminal(prompt); err != nil {
		return "", err
	}

	err = survey.AskOne(&survey.Input{
		Message: prompt,
		Default: defaultValue,
		Help:    help,
	}, &result)

	return
}

func (p *prompter) Confirm(prompt string, defaultValue bool) (bool, error) {
	if err := ensureTerminal(prompt); err != nil {
		return false, err
	}

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

	if err := ensureTerminal(input.Message); err != nil {
		return err
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
