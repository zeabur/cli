package graphql_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/zeabur/cli/internal/cmd/api/graphql"
	"github.com/zeabur/cli/internal/cmdutil"
	"github.com/zeabur/cli/pkg/api"
)

// stubClient records the request body and returns a canned response. Other
// Client methods inherit the embedded nil interface and panic if called.
type stubClient struct {
	api.Client

	gotBody []byte
	resp    []byte
	err     error
}

func (c *stubClient) GraphQLRaw(_ context.Context, body []byte) ([]byte, error) {
	c.gotBody = body
	return c.resp, c.err
}

type request struct {
	Query         string         `json:"query"`
	OperationName string         `json:"operationName"`
	Variables     map[string]any `json:"variables"`
}

func run(t *testing.T, client *stubClient, stdin string, args ...string) (string, error) {
	t.Helper()

	f := &cmdutil.Factory{ApiClient: client}
	cmd := graphql.NewCmdGraphQL(f)
	out := &bytes.Buffer{}
	cmd.SilenceUsage = true // mirrors the root command
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})
	cmd.SetIn(strings.NewReader(stdin))
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), err
}

func decode(t *testing.T, body []byte) request {
	t.Helper()

	var req request
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("request body is not JSON: %v (%s)", err, body)
	}
	return req
}

func TestGraphQL_FieldsBuildRequest(t *testing.T) {
	t.Parallel()

	client := &stubClient{resp: []byte(`{"data":{"me":{"_id":"1"}}}`)}
	out, err := run(t, client, "",
		"-f", "query=query($id: ObjectID!) { project(_id: $id) { name } }",
		"-f", "operationName=Q",
		"-f", "id=123",
		"-F", "limit=10",
		"-F", "flag=true",
		"-F", "none=null",
		"-F", "name=abc",
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out != "{\"data\":{\"me\":{\"_id\":\"1\"}}}\n" {
		t.Fatalf("unexpected output: %q", out)
	}

	req := decode(t, client.gotBody)
	if req.Query != "query($id: ObjectID!) { project(_id: $id) { name } }" || req.OperationName != "Q" {
		t.Fatalf("unexpected query/operationName: %+v", req)
	}
	want := map[string]any{"id": "123", "limit": float64(10), "flag": true, "none": nil, "name": "abc"}
	if len(req.Variables) != len(want) {
		t.Fatalf("unexpected variables: %#v", req.Variables)
	}
	for k, v := range want {
		if got, ok := req.Variables[k]; !ok || got != v {
			t.Fatalf("variable %q = %#v, want %#v", k, got, v)
		}
	}
}

func TestGraphQL_QueryFromFileAndNoVariables(t *testing.T) {
	t.Parallel()

	path := filepath.Join(t.TempDir(), "q.graphql")
	if err := os.WriteFile(path, []byte("query { me { _id } }"), 0o600); err != nil {
		t.Fatal(err)
	}

	client := &stubClient{resp: []byte(`{"data":{}}`)}
	if _, err := run(t, client, "", "-F", "query=@"+path); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(client.gotBody) != `{"query":"query { me { _id } }"}` {
		t.Fatalf("unexpected body: %s", client.gotBody)
	}
}

func TestGraphQL_InputFromStdinMergesFields(t *testing.T) {
	t.Parallel()

	client := &stubClient{resp: []byte(`{"data":{}}`)}
	stdin := `{"query":"query { me { _id } }","variables":{"a":1,"b":"x"}}`
	if _, err := run(t, client, stdin, "--input", "-", "-f", "b=y"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	req := decode(t, client.gotBody)
	if req.Query != "query { me { _id } }" || req.Variables["a"] != float64(1) || req.Variables["b"] != "y" {
		t.Fatalf("unexpected request: %+v", req)
	}
}

func TestGraphQL_GraphQLErrorsExitNonZero(t *testing.T) {
	t.Parallel()

	resp := `{"errors":[{"message":"boom"}],"data":null}`
	client := &stubClient{resp: []byte(resp)}
	out, err := run(t, client, "", "-f", "query=query { me { _id } }")
	if err == nil || !strings.Contains(err.Error(), "1 error(s)") {
		t.Fatalf("expected GraphQL error, got %v", err)
	}
	if out != resp+"\n" {
		t.Fatalf("response body should still be printed, got %q", out)
	}
}

func TestGraphQL_HTTPErrorPrintsBody(t *testing.T) {
	t.Parallel()

	client := &stubClient{resp: []byte(`{"errors":[{"message":"unauthorized"}]}`), err: errors.New("graphql request failed: HTTP 401")}
	out, err := run(t, client, "", "-f", "query=query { me { _id } }")
	if err == nil || !strings.Contains(err.Error(), "HTTP 401") || !strings.Contains(err.Error(), "1 GraphQL error(s)") {
		t.Fatalf("expected HTTP error, got %v", err)
	}
	if !strings.Contains(out, "unauthorized") {
		t.Fatalf("response body should still be printed, got %q", out)
	}
}

func TestGraphQL_InvalidInput(t *testing.T) {
	t.Parallel()

	cases := map[string][]string{
		"missing query":     {"-f", "id=1"},
		"malformed field":   {"-f", "query"},
		"empty key":         {"-f", "=x"},
		"stdin read twice":  {"--input", "-", "-F", "query=@-"},
		"positional args":   {"extra"},
		"missing file":      {"-F", "query=@/nonexistent/q.graphql"},
		"variables not obj": {"--input", "-"},
	}
	for name, args := range cases {
		stdin := `{"query":"query { me { _id } }"}`
		if name == "variables not obj" {
			stdin = `{"query":"query { me { _id } }","variables":[1]}`
		}
		client := &stubClient{}
		if _, err := run(t, client, stdin, args...); err == nil {
			t.Errorf("%s: expected an error", name)
		}
		if client.gotBody != nil {
			t.Errorf("%s: request must not be sent", name)
		}
	}
}
