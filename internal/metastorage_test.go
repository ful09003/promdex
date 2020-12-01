//+build all_tests

package internal_test

import (
	"atamedomain.name/promdex/internal"
	"testing"
)

func TestStringingStorageTypesIsLogical(t *testing.T) {
	var table = []struct {
		in       internal.PromdexStorageType
		expected string
	}{
		{internal.NilStore, "null"},
		{internal.SQLiteStore, "sqlite"},
	}
	for _, tt := range table {
		if tt.in.String() != tt.expected {
			t.Errorf("got %s, expected %s", tt.in.String(), tt.expected)
		}
	}
}
