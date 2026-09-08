package env_test

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.uber.org/zap"

	variableenv "github.com/zeabur/cli/internal/cmd/variable/env"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
)

type stubVariableClient struct {
	api.Client

	updated map[string]string
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

func TestEnvDoesNotExposeImportedVariableValues(t *testing.T) {
	for _, jsonOutput := range []bool{false, true} {
		t.Run(map[bool]string{false: "table", true: "json"}[jsonOutput], func(t *testing.T) {
			envFile := filepath.Join(t.TempDir(), ".env")
			if err := os.WriteFile(envFile, []byte("API_TOKEN=imported-secret-value\n"), 0o600); err != nil {
				t.Fatalf("write env fixture: %v", err)
			}

			client := &stubVariableClient{}
			printer := &capturePrinter{}
			f := &cmdutil.Factory{
				ApiClient: client,
				Printer:   printer,
				Log:       zap.NewNop().Sugar(),
			}
			f.Interactive = false
			f.JSON = jsonOutput

			cmd := variableenv.NewCmdEnvVariable(f)
			cmd.SetArgs([]string{
				"--id", "service-id",
				"--env-id", "environment-id",
				"--file", envFile,
			})
			cmd.SilenceErrors = true
			cmd.SilenceUsage = true
			if err := cmd.Execute(); err != nil {
				t.Fatalf("import variables: %v", err)
			}

			output := printer.output()
			if strings.Contains(output, "imported-secret-value") {
				t.Fatalf("env output exposed imported variable value: %s", output)
			}
			if !strings.Contains(output, "API_TOKEN") {
				t.Fatalf("env output did not identify imported key: %s", output)
			}
			if got := client.updated["API_TOKEN"]; got != "imported-secret-value" {
				t.Fatalf("imported value missing from update payload: %q", got)
			}
		})
	}
}
