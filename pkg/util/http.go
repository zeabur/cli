package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FormatHTTPError reads a non-2xx HTTP response and returns an error that
// preserves whatever the server actually said.
//
// The Zeabur API returns JSON bodies of the form
//
//	{"code": "REQUIRE_PAID_PLAN", "error": "Upload size ... exceeds ..."}
//
// for client-visible failures. When that shape is present, the message and
// (if available) code are surfaced. Otherwise the raw body is included so
// the user still has something to act on instead of a bare status code.
//
// action is a short verb phrase describing what failed, e.g. "create upload session".
func FormatHTTPError(action string, resp *http.Response) error {
	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return fmt.Errorf("%s: failed to read response body (status: %d): %w", action, resp.StatusCode, err)
	}

	var errResp struct {
		Code  string `json:"code"`
		Error string `json:"error"`
	}
	if err := json.Unmarshal(body, &errResp); err == nil && errResp.Error != "" {
		if errResp.Code != "" {
			return fmt.Errorf("%s: %s (code: %s, status: %d)", action, errResp.Error, errResp.Code, resp.StatusCode)
		}
		return fmt.Errorf("%s: %s (status: %d)", action, errResp.Error, resp.StatusCode)
	}

	if len(body) > 0 {
		return fmt.Errorf("%s: status code %d, body: %s", action, resp.StatusCode, string(body))
	}
	return fmt.Errorf("%s: status code %d", action, resp.StatusCode)
}
