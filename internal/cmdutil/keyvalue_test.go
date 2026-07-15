package cmdutil

import (
	"reflect"
	"testing"
)

func TestParseKeyValuePairs(t *testing.T) {
	tests := []struct {
		name    string
		values  []string
		want    map[string]string
		wantErr bool
	}{
		{
			name:   "preserves variable reference",
			values: []string{"DATABASE_URL=${POSTGRESQL.POSTGRES_CONNECTION_STRING}"},
			want:   map[string]string{"DATABASE_URL": "${POSTGRESQL.POSTGRES_CONNECTION_STRING}"},
		},
		{
			name:   "preserves commas and equals signs",
			values: []string{"OPTIONS=one,two", "TOKEN=header=payload=signature"},
			want:   map[string]string{"OPTIONS": "one,two", "TOKEN": "header=payload=signature"},
		},
		{
			name:   "allows empty value",
			values: []string{"EMPTY="},
			want:   map[string]string{"EMPTY": ""},
		},
		{
			name:   "last duplicate wins",
			values: []string{"KEY=first", "KEY=second"},
			want:   map[string]string{"KEY": "second"},
		},
		{
			name:    "rejects missing separator",
			values:  []string{"INVALID"},
			wantErr: true,
		},
		{
			name:    "rejects empty key",
			values:  []string{"=value"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseKeyValuePairs(tt.values)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ParseKeyValuePairs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseKeyValuePairs() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
