// +build all_tests,localtests

package internal_test

import (
	"atamedomain.name/promdex/internal"
	"testing"
)

func TestDiscoveryAgainstLocalProm(t *testing.T) {
	testProm, e := internal.NewPromdexTarget("http://localhost:9090")
	if e != nil {
		t.Error(e)
	}
	if testProm.GetTarget() != "http://localhost:9090" {
		t.Errorf("expected http://localhost:9090 for target, received %s", testProm.GetTarget())
	}
}

func TestDiscoveryReceivesDataFromLocalProm(t *testing.T) {
	testProm, e := internal.NewPromdexTarget("http://localhost:9090")
	if e != nil {
		t.Error(e)
	}

	d, e := testProm.QueryTargetMetadata()

	if e != nil {
		t.Error(e)
	}
	
	if len(d) == 0 {
		t.Errorf("expected >0 results querying local Prometheus, received %d", len(d))
	}
}