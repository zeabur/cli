package prompt

import "github.com/AlecAivazis/survey/v2"

// todo: complete this implementation
type prompter struct {
}

func New() Prompter {
	return &prompter{}
}

func (p *prompter) Select(message string, defaultValue string, options []string) (int, error) {
	//TODO implement me
	panic("implement me")
}

func (p *prompter) MultiSelect(message string, defaultValues, options []string) ([]int, error) {
	//TODO implement me
	panic("implement me")
}

func (p *prompter) Input(prompt, defaultValue string) (string, error) {
	//TODO implement me
	panic("implement me")
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
	//TODO implement me
	panic("implement me")
}

var _ Prompter = (*prompter)(nil)
