package util_test

import (
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeabur/cli/pkg/util"
)

//go:embed testdata/01-valid.yaml
var validTemplateSpec []byte

//go:embed testdata/02-invalid-chinese.yaml
var invalidChineseTemplateSpec []byte

//go:embed testdata/03-invalid-yaml.yaml
var invalidYAMLTemplateSpec []byte

//go:embed testdata/04-invalid-missing-field.yaml
var invalidMissingFieldSpec []byte

func TestValidateTemplate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		templateSpec []byte
		expectError  string
	}{
		{
			name:         "valid template",
			templateSpec: validTemplateSpec,
			expectError:  "",
		},
		{
			name:         "invalid chinese",
			templateSpec: invalidChineseTemplateSpec,
			expectError:  "'網站' does not match pattern '^[a-zA-Z][ -~]*$'",
		},
		{
			name:         "invalid yaml",
			templateSpec: invalidYAMLTemplateSpec,
			expectError:  "mapping values are not allowed in this context",
		},
		{
			name:         "invalid missing field",
			templateSpec: invalidMissingFieldSpec,
			expectError:  "got null, want object", // fixme: friendly error message
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := util.ValidateTemplate(tt.templateSpec)
			if tt.expectError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.ErrorContains(t, err, tt.expectError)
			}
		})
	}
}
