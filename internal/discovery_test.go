// +build all_tests

package internal_test

import (
	"atamedomain.name/promdex/internal"
	"testing"
)

func TestDiscoveryConfiguratorHandlesBadURLs(t *testing.T) {
	var urls = []string{"!htt1://localhost:9100", "##:!", "h_t_t_p_s:&", "nil"}

	for _, u := range urls {
		_, e := internal.NewPromdexTarget(u)
		if e == nil {
			t.Errorf("%s slipped through validation: %s", u, e)
		}
	}
}