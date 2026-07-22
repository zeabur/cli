package model_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeabur/cli/pkg/model"
)

func ptr[T any](v T) *T { return &v }

// indexOf returns the column index of name in header, failing the test when the
// column is missing — keeps the assertions below readable and order-independent.
func indexOf(t *testing.T, header []string, name string) int {
	t.Helper()
	for i, h := range header {
		if h == name {
			return i
		}
	}
	t.Fatalf("column %q not found in header %v", name, header)
	return -1
}

func TestServerListItemsOSColumn(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		hasK3s *bool
		want   string
	}{
		{name: "k3s installed", hasK3s: ptr(true), want: "ZeaburOS"},
		{name: "explicitly no k3s", hasK3s: ptr(false), want: "Ubuntu"},
		// A legacy server predates the field and does have Zeabur services; the
		// backend infers that from certificate data but exposes the raw column,
		// so null must never be read as "no k3s".
		{name: "legacy server reports null", hasK3s: nil, want: "ZeaburOS"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			items := model.ServerListItems{{ID: "abc", Name: "s", HasK3s: tt.hasK3s}}
			rows := items.Rows()
			require.Len(t, rows, 1)
			assert.Equal(t, tt.want, rows[0][indexOf(t, items.Header(), "OS")])
		})
	}
}

func TestServerDetailOSColumn(t *testing.T) {
	t.Parallel()

	server := &model.ServerDetail{ID: "abc", Name: "s", HasK3s: ptr(false)}
	rows := server.Rows()
	require.Len(t, rows, 1)
	assert.Equal(t, "Ubuntu", rows[0][indexOf(t, server.Header(), "OS")])
}

func TestServerDetailResourceColumns(t *testing.T) {
	t.Parallel()

	t.Run("renders usage when measured", func(t *testing.T) {
		t.Parallel()

		server := &model.ServerDetail{
			Status: model.ServerStatus{
				UsedCPU: 1000, TotalCPU: 4000,
				UsedMemory: 512, TotalMemory: 2048,
				UsedDisk: 3000, TotalDisk: 40000,
			},
		}
		row := server.Rows()[0]
		header := server.Header()
		// CPU is millicores, so the unit must be spelled out — a 4-core machine
		// reports 4000 and would otherwise read as 4000 cores.
		assert.Equal(t, "1000/4000 m", row[indexOf(t, header, "CPU")])
		assert.Equal(t, "512/2048 MB", row[indexOf(t, header, "Memory")])
		assert.Equal(t, "3000/40000 MB", row[indexOf(t, header, "Disk")])
	})

	// A real machine cannot have zero cores, so a zero total means the metrics
	// were never collected. Printing "0/0" would read as a machine in trouble.
	t.Run("renders a dash when not measured", func(t *testing.T) {
		t.Parallel()

		server := &model.ServerDetail{}
		row := server.Rows()[0]
		header := server.Header()
		assert.Equal(t, "—", row[indexOf(t, header, "CPU")])
		assert.Equal(t, "—", row[indexOf(t, header, "Memory")])
		assert.Equal(t, "—", row[indexOf(t, header, "Disk")])
	})
}

func TestServerJSONCarriesOS(t *testing.T) {
	t.Parallel()

	t.Run("detail", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(&model.ServerDetail{ID: "abc", HasK3s: ptr(false)})
		require.NoError(t, err)

		var got map[string]any
		require.NoError(t, json.Unmarshal(data, &got))
		assert.Equal(t, "Ubuntu", got["os"])
		assert.Equal(t, "abc", got["ID"])
		// The raw three-state field stays available for callers that need it.
		assert.Equal(t, false, got["HasK3s"])
	})

	t.Run("list", func(t *testing.T) {
		t.Parallel()

		data, err := json.Marshal(model.ServerListItems{{ID: "abc", HasK3s: ptr(true)}})
		require.NoError(t, err)

		var got []map[string]any
		require.NoError(t, json.Unmarshal(data, &got))
		require.Len(t, got, 1)
		assert.Equal(t, "ZeaburOS", got[0]["os"])
	})
}
