package delete_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"

	variabledelete "github.com/zeabur/cli/internal/cmd/variable/delete"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/fill"
	"github.com/zeabur/cli/pkg/model"
	"github.com/zeabur/cli/pkg/prompt"
)

type stubVariableClient struct {
	api.Client

	variables model.Variables
	updated   map[string]string
}

func (c *stubVariableClient) ListVariables(context.Context, string, string) (model.Variables, model.Variables, error) {
	return c.variables, nil, nil
}

func (c *stubVariableClient) UpdateVariables(_ context.Context, _, _ string, data map[string]string) (bool, error) {
	c.updated = data
	return true, nil
}

type capturePrinter struct {
	header    []string
	rows      [][]string
	jsonValue any
}

type stubParamFiller struct {
	fill.ParamFiller
}

func (*stubParamFiller) ServiceByNameWithEnvironment(fill.ServiceByNameWithEnvironmentOptions) (bool, error) {
	return false, nil
}

type capturePrompter struct {
	prompt.Prompter

	options []string
}

func (p *capturePrompter) Select(_, _ string, options []string) (int, error) {
	p.options = options
	return 0, nil
}

func (*capturePrompter) Confirm(string, bool) (bool, error) {
	return true, nil
}

func (p *capturePrinter) Table(header []string, rows [][]string) {
	p.header = header
	p.rows = rows
}

func (p *capturePrinter) JSON(v any) error {
	p.jsonValue = v
	return nil
}

func (p *capturePrinter) output() string {
	encoded, _ := json.Marshal(struct {
		Header []string
		Rows   [][]string
		JSON   any
	}{p.header, p.rows, p.jsonValue})
	return string(encoded)
}

func TestDeleteDoesNotExposeRemainingVariableValues(t *testing.T) {
	client := &stubVariableClient{variables: model.Variables{
		{Key: "DELETE_ME", Value: "deleted-secret-value"},
		{Key: "KEEP_ME", Value: "remaining-secret-value"},
	}}
	printer := &capturePrinter{}
	f := &cmdutil.Factory{
		ApiClient: client,
		Printer:   printer,
		Log:       zap.NewNop().Sugar(),
	}
	f.Interactive = false

	cmd := variabledelete.NewCmdDeleteVariable(f)
	cmd.SetArgs([]string{
		"--id", "service-id",
		"--env-id", "environment-id",
		"--delete-keys", "DELETE_ME",
	})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delete variable: %v", err)
	}

	output := printer.output()
	for _, sensitive := range []string{"deleted-secret-value", "remaining-secret-value", "KEEP_ME"} {
		if strings.Contains(output, sensitive) {
			t.Fatalf("delete output exposed remaining or sensitive data %q: %s", sensitive, output)
		}
	}
	if !strings.Contains(output, "DELETE_ME") {
		t.Fatalf("delete output did not identify the deleted key: %s", output)
	}
	if _, ok := client.updated["DELETE_ME"]; ok {
		t.Fatal("deleted variable remained in update payload")
	}
	if got := client.updated["KEEP_ME"]; got != "remaining-secret-value" {
		t.Fatalf("remaining variable was not preserved in update payload: %q", got)
	}
}

func TestDeleteInteractiveMasksValuesInSelection(t *testing.T) {
	client := &stubVariableClient{variables: model.Variables{
		{Key: "DELETE_ME", Value: "deleted-secret-value"},
	}}
	prompter := &capturePrompter{}
	f := &cmdutil.Factory{
		ApiClient:   client,
		Printer:     &capturePrinter{},
		Log:         zap.NewNop().Sugar(),
		ParamFiller: &stubParamFiller{},
		Prompter:    prompter,
	}
	f.Interactive = true

	cmd := variabledelete.NewCmdDeleteVariable(f)
	cmd.SetArgs([]string{
		"--id", "service-id",
		"--env-id", "environment-id",
	})
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	if err := cmd.Execute(); err != nil {
		t.Fatalf("delete variable interactively: %v", err)
	}

	selection := strings.Join(prompter.options, "\n")
	if strings.Contains(selection, "deleted-secret-value") {
		t.Fatalf("interactive delete selection exposed variable value: %s", selection)
	}
	if !strings.Contains(selection, "del***") {
		t.Fatalf("interactive delete selection did not contain masked value: %s", selection)
	}
}
