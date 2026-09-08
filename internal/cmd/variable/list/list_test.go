package list_test

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"go.uber.org/zap"

	variablelist "github.com/zeabur/cli/internal/cmd/variable/list"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
	"github.com/zeabur/cli/pkg/model"
)

type stubVariableClient struct {
	api.Client

	variables         model.Variables
	readonlyVariables model.Variables
}

func (c *stubVariableClient) ListVariables(context.Context, string, string) (model.Variables, model.Variables, error) {
	return c.variables, c.readonlyVariables, nil
}

type capturePrinter struct {
	tables    [][][]string
	jsonValue any
}

func (p *capturePrinter) Table(header []string, rows [][]string) {
	table := make([][]string, 0, len(rows)+1)
	table = append(table, header)
	table = append(table, rows...)
	p.tables = append(p.tables, table)
}

func (p *capturePrinter) JSON(v any) error {
	p.jsonValue = v
	return nil
}

func (p *capturePrinter) output() string {
	encoded, _ := json.Marshal(struct {
		Tables [][][]string
		JSON   any
	}{p.tables, p.jsonValue})
	return string(encoded)
}

func TestListMasksVariableValuesByDefault(t *testing.T) {
	for _, jsonOutput := range []bool{false, true} {
		t.Run(map[bool]string{false: "table", true: "json"}[jsonOutput], func(t *testing.T) {
			printer := runList(t, jsonOutput)
			output := printer.output()
			for _, secret := range []string{"service-database-password", "shared-workspace-token"} {
				if strings.Contains(output, secret) {
					t.Fatalf("list output exposed variable value %q: %s", secret, output)
				}
			}
			if !strings.Contains(output, "ser***") || !strings.Contains(output, "sha***") {
				t.Fatalf("list output did not contain masked values: %s", output)
			}
		})
	}
}

func TestListShowValuesRequiresExplicitFlag(t *testing.T) {
	for _, jsonOutput := range []bool{false, true} {
		t.Run(map[bool]string{false: "table", true: "json"}[jsonOutput], func(t *testing.T) {
			output := runList(t, jsonOutput, "--show-values").output()
			for _, secret := range []string{"service-database-password", "shared-workspace-token"} {
				if !strings.Contains(output, secret) {
					t.Fatalf("--show-values did not reveal %q: %s", secret, output)
				}
			}
		})
	}
}

func runList(t *testing.T, jsonOutput bool, extraArgs ...string) *capturePrinter {
	t.Helper()
	client := &stubVariableClient{
		variables: model.Variables{
			{Key: "DATABASE_PASSWORD", Value: "service-database-password", ServiceID: "service-id"},
		},
		readonlyVariables: model.Variables{
			{Key: "SHARED_TOKEN", Value: "shared-workspace-token", ServiceID: "other-service-id"},
		},
	}
	printer := &capturePrinter{}
	f := &cmdutil.Factory{
		ApiClient: client,
		Printer:   printer,
		Log:       zap.NewNop().Sugar(),
	}
	f.Interactive = false
	f.JSON = jsonOutput

	cmd := variablelist.NewCmdListVariables(f)
	args := make([]string, 0, 4+len(extraArgs))
	args = append(args,
		"--id", "service-id",
		"--env-id", "environment-id",
	)
	cmd.SetArgs(append(args, extraArgs...))
	cmd.SilenceErrors = true
	cmd.SilenceUsage = true
	if err := cmd.Execute(); err != nil {
		t.Fatalf("list variables: %v", err)
	}
	return printer
}
