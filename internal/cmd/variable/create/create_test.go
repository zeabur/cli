package create_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"

	"github.com/zeabur/cli/internal/cmd/variable/create"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
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

func TestCreateDoesNotExposeVariableValues(t *testing.T) {
	for _, jsonOutput := range []bool{false, true} {
		t.Run(map[bool]string{false: "table", true: "json"}[jsonOutput], func(t *testing.T) {
			client := &stubVariableClient{variables: model.Variables{
				{Key: "DATABASE_PASSWORD", Value: "existing-database-password"},
			}}
			printer := &capturePrinter{}
			f := &cmdutil.Factory{
				ApiClient: client,
				Printer:   printer,
				Log:       zap.NewNop().Sugar(),
			}
			f.Interactive = false
			f.JSON = jsonOutput

			cmd := create.NewCmdCreateVariable(f)
			cmd.SetArgs([]string{
				"--id", "service-id",
				"--env-id", "environment-id",
				"--key", "NEW_SECRET=new-secret-value",
			})
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			if err := cmd.Execute(); err != nil {
				t.Fatalf("create variable: %v", err)
			}

			output := printer.output()
			for _, sensitive := range []string{
				"existing-database-password",
				"new-secret-value",
				"DATABASE_PASSWORD",
			} {
				if strings.Contains(output, sensitive) {
					t.Fatalf("create output exposed existing or sensitive data %q: %s", sensitive, output)
				}
			}
			if !strings.Contains(output, "NEW_SECRET") {
				t.Fatalf("create output did not identify the created key: %s", output)
			}
			if got := client.updated["DATABASE_PASSWORD"]; got != "existing-database-password" {
				t.Fatalf("existing variable was not preserved in update payload: %q", got)
			}
			if got := client.updated["NEW_SECRET"]; got != "new-secret-value" {
				t.Fatalf("new variable missing from update payload: %q", got)
			}
		})
	}
}
