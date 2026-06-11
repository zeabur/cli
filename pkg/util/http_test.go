package util_test

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

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

func TestUploadTimeout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size int64
		want time.Duration
	}{
		{"zero size gets the 2min floor", 0, 2 * time.Minute},
		{"tiny payload still gets the floor", 1024, 2 * time.Minute},
		{"1 MiB adds 4s over the floor", 1 << 20, 2*time.Minute + 4*time.Second},
		{"50 MiB scales to ~5.3 min", 50 << 20, 2*time.Minute + 200*time.Second},
		{"1 GiB scales to ~70 min", 1 << 30, 2*time.Minute + 4096*time.Second},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, tc.want, util.UploadTimeout(tc.size))
		})
	}
}
