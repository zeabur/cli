package util_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeabur/cli/pkg/util"
)

func TestFormatHTTPError(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		statusCode int
		body       string
		wantSubstr []string
	}{
		{
			name:       "structured error with code (REQUIRE_PAID_PLAN)",
			statusCode: 400,
			body:       `{"code":"REQUIRE_PAID_PLAN","error":"Upload size 53048320 bytes exceeds your plan's limit of 52428800 bytes.","limit":52428800,"content_length":53048320}`,
			wantSubstr: []string{
				"do thing",
				"Upload size 53048320 bytes exceeds",
				"code: REQUIRE_PAID_PLAN",
				"status: 400",
			},
		},
		{
			name:       "structured error without code",
			statusCode: 403,
			body:       `{"error":"forbidden"}`,
			wantSubstr: []string{"do thing", "forbidden", "status: 403"},
		},
		{
			name:       "non-JSON body falls back to raw body",
			statusCode: 502,
			body:       "<html>Bad Gateway</html>",
			wantSubstr: []string{"do thing", "status code 502", "<html>Bad Gateway</html>"},
		},
		{
			name:       "empty body falls back to status code only",
			statusCode: 504,
			body:       "",
			wantSubstr: []string{"do thing", "status code 504"},
		},
		{
			name:       "JSON without error field falls back to raw body",
			statusCode: 400,
			body:       `{"foo":"bar"}`,
			wantSubstr: []string{"status code 400", `{"foo":"bar"}`},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			resp := &http.Response{
				StatusCode: tc.statusCode,
				Body:       io.NopCloser(strings.NewReader(tc.body)),
			}

			err := util.FormatHTTPError("do thing", resp)
			assert.Error(t, err)
			for _, s := range tc.wantSubstr {
				assert.Contains(t, err.Error(), s)
			}
		})
	}
}
