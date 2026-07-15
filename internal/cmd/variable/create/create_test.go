package create

import (
	"reflect"
	"testing"
)

func TestKeyFlagPreservesValues(t *testing.T) {
	cmd := NewCmdCreateVariable(nil)
	want := []string{
		"DATABASE_URL=${POSTGRESQL.POSTGRES_CONNECTION_STRING}",
		"OPTIONS=one,two",
	}

	if err := cmd.ParseFlags([]string{"--key", want[0], "--key", want[1]}); err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	got, err := cmd.Flags().GetStringArray("key")
	if err != nil {
		t.Fatalf("GetStringArray() error = %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("key values = %#v, want %#v", got, want)
	}
}
