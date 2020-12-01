package internal

import "testing"

func TestStringingStorageTypesIsLogical(t *testing.T) {
	var table = []struct {
		in       PromdexStorageType
		expected string
	}{
		{NilStore, "null"},
		{SQLiteStore, "sqlite"},
	}
	for _, tt := range table {
		if tt.in.String() != tt.expected {
			t.Errorf("got %s, expected %s", tt.in.String(), tt.expected)
		}
	}
}
