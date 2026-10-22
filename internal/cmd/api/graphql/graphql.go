// Package graphql provides `zeabur api graphql`, which sends an arbitrary
// GraphQL request to the Zeabur API using the CLI's stored credentials.
package graphql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"

	"github.com/zeabur/cli/internal/cmdutil"
)

type Options struct {
	rawFields   []string
	typedFields []string
	input       string

	stdin io.Reader
}

func NewCmdGraphQL(f *cmdutil.Factory) *cobra.Command {
	opts := &Options{}

	cmd := &cobra.Command{
		Use:   "graphql",
		Short: "Send a GraphQL request to the Zeabur API",
		Long: heredoc.Doc(`
			Send an arbitrary GraphQL query or mutation to the Zeabur API with the
			credentials of the logged-in user, and print the raw JSON response.

			The query is passed as the "query" field; every other field becomes a
			GraphQL variable. "operationName" selects the operation to run when the
			document contains several.

			  -f/--raw-field key=value   adds a string value.
			  -F/--field key=value       adds a typed value: true, false, null and
			                             integers are converted to JSON types, and a
			                             value starting with "@" is read from a file
			                             ("@-" reads standard input).

			--input reads a complete JSON request body ({"query": ..., "variables": ...})
			from a file, or from standard input with "-". Fields given with -f/-F are
			merged on top of it.

			The command exits with a non-zero status when the request fails or the
			response contains GraphQL errors; the response body is still printed.

			The request is sent as-is: --workspace and the saved context do not rewrite
			it, so pass team or project IDs explicitly as variables.
		`),
		Example: heredoc.Doc(`
			$ zeabur api graphql -f query='query { me { _id username } }'
			$ zeabur api graphql -f query='query($id: ObjectID!) { project(_id: $id) { name } }' -F id=<project-id>
			$ zeabur api graphql -F query=@query.graphql -F limit=10
			$ echo '{"query": "query { me { _id } }"}' | zeabur api graphql --input -
		`),
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.stdin = cmd.InOrStdin()
			return runGraphQL(cmd.Context(), f, opts, cmd.OutOrStdout())
		},
	}

	cmd.Flags().StringArrayVarP(&opts.rawFields, "raw-field", "f", nil, "Add a string field in key=value format")
	cmd.Flags().StringArrayVarP(&opts.typedFields, "field", "F", nil, "Add a typed field in key=value format (use @file or @- to read a value)")
	cmd.Flags().StringVar(&opts.input, "input", "", "Read the JSON request body from a file (use \"-\" for standard input)")

	return cmd
}

func runGraphQL(ctx context.Context, f *cmdutil.Factory, opts *Options, out io.Writer) error {
	if ctx == nil {
		ctx = context.Background()
	}

	body, err := buildRequestBody(opts)
	if err != nil {
		return err
	}

	resp, reqErr := f.ApiClient.GraphQLRaw(ctx, body)
	if len(resp) > 0 {
		if _, err := out.Write(resp); err != nil {
			return fmt.Errorf("write response: %w", err)
		}
		if resp[len(resp)-1] != '\n' {
			_, _ = io.WriteString(out, "\n")
		}
	}

	var parsed struct {
		Errors []json.RawMessage `json:"errors"`
	}
	parseErr := json.Unmarshal(resp, &parsed)

	if reqErr != nil {
		if n := len(parsed.Errors); parseErr == nil && n > 0 {
			return fmt.Errorf("%w: response contains %d GraphQL error(s)", reqErr, n)
		}
		return reqErr
	}
	if parseErr != nil {
		return fmt.Errorf("parse response: %w", parseErr)
	}
	if n := len(parsed.Errors); n > 0 {
		return fmt.Errorf("GraphQL response contains %d error(s)", n)
	}

	return nil
}

// buildRequestBody assembles the JSON request body from --input and the
// -f/-F fields. "query" and "operationName" are top-level keys; every other
// field is a variable.
func buildRequestBody(opts *Options) ([]byte, error) {
	stdinUsed := false
	readStdin := func(what string) ([]byte, error) {
		if stdinUsed {
			return nil, fmt.Errorf("%s: standard input can only be read once", what)
		}
		stdinUsed = true
		return io.ReadAll(opts.stdin)
	}

	req := map[string]any{}
	variables := map[string]any{}

	if opts.input != "" {
		var (
			raw []byte
			err error
		)
		if opts.input == "-" {
			raw, err = readStdin("--input")
		} else {
			raw, err = os.ReadFile(opts.input)
		}
		if err != nil {
			return nil, fmt.Errorf("read --input: %w", err)
		}
		if err := json.Unmarshal(raw, &req); err != nil {
			return nil, fmt.Errorf("--input must be a JSON object: %w", err)
		}
		if v, ok := req["variables"]; ok && v != nil {
			vars, ok := v.(map[string]any)
			if !ok {
				return nil, errors.New("--input: \"variables\" must be a JSON object")
			}
			variables = vars
		}
	}

	set := func(key string, value any) {
		switch key {
		case "query", "operationName":
			req[key] = value
		default:
			variables[key] = value
		}
	}

	for _, field := range opts.rawFields {
		key, value, err := splitField(field)
		if err != nil {
			return nil, err
		}
		set(key, value)
	}

	for _, field := range opts.typedFields {
		key, value, err := splitField(field)
		if err != nil {
			return nil, err
		}
		typed, err := parseTypedValue(value, readStdin)
		if err != nil {
			return nil, fmt.Errorf("field %q: %w", key, err)
		}
		set(key, typed)
	}

	query, ok := req["query"].(string)
	if !ok || strings.TrimSpace(query) == "" {
		return nil, errors.New("a query is required: pass -f query='...', -F query=@file, or --input")
	}

	if len(variables) > 0 {
		req["variables"] = variables
	}

	return json.Marshal(req)
}

func splitField(field string) (string, string, error) {
	key, value, ok := strings.Cut(field, "=")
	if !ok || key == "" {
		return "", "", fmt.Errorf("field %q must be in key=value format", field)
	}
	return key, value, nil
}

func parseTypedValue(value string, readStdin func(string) ([]byte, error)) (any, error) {
	switch value {
	case "true":
		return true, nil
	case "false":
		return false, nil
	case "null":
		return nil, nil
	}

	if path, ok := strings.CutPrefix(value, "@"); ok {
		var (
			content []byte
			err     error
		)
		if path == "-" {
			content, err = readStdin("@-")
		} else {
			content, err = os.ReadFile(path)
		}
		if err != nil {
			return nil, fmt.Errorf("read %s: %w", value, err)
		}
		return string(content), nil
	}

	if n, err := strconv.ParseInt(value, 10, 64); err == nil {
		return n, nil
	}

	return value, nil
}
