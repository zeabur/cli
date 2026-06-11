package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// UploadTimeout returns how long an upload of contentLength bytes is allowed
// to take: a 2-minute floor plus one second per 256 KiB of payload.
//
// http.Client.Timeout covers the entire request including body transfer, so a
// fixed timeout makes large uploads fail on slow links no matter how patient
// the user is. The 256 KiB/s floor is deliberately conservative (~2 Mbps);
// a 50 MiB zip gets ~5.3 minutes.
func UploadTimeout(contentLength int64) time.Duration {
	return 2*time.Minute + time.Duration(contentLength/(256*1024))*time.Second
}

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
