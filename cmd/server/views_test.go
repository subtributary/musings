package main

import (
	"testing"
	"time"
)

func TestView_SetData(t *testing.T) {
	t.Parallel()

	older := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	templateTime := time.Date(2026, 1, 2, 12, 0, 0, 0, time.UTC)
	newer := time.Date(2026, 1, 3, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		dataModified time.Time
		want         time.Time
	}{
		{
			name: "no data set",
			want: templateTime,
		},
		{
			name:         "data older than template",
			dataModified: older,
			want:         templateTime,
		},
		{
			name:         "data newer than template",
			dataModified: newer,
			want:         newer,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			view := View{modified: templateTime}

			if !tt.dataModified.IsZero() {
				_ = view.SetData(nil, tt.dataModified)
			}

			if view.modified != tt.want {
				t.Errorf("modified = %v, want %v", view.modified, tt.want)
			}
		})
	}
}

// ViewFactory is too complex for unit tests.
// We need to test it with an integration test.
