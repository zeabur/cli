package value_test

import (
	"testing"

	"github.com/zeabur/cli/internal/cmd/variable/value"
)

func TestMask(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want string
	}{
		{name: "empty", raw: "", want: "***"},
		{name: "short", raw: "ab", want: "***"},
		{name: "long", raw: "secret", want: "sec***"},
		{name: "unicode", raw: "密碼內容", want: "密碼內***"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := value.Mask(tt.raw); got != tt.want {
				t.Fatalf("Mask(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}
}
