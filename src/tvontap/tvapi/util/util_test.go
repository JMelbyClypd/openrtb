package util

import (
	"testing"
)

var KEYS = []string{"orders", "orders", "buyer123", "reqabcd"}

func TestMungeUnmunge(t *testing.T) {
	key := Munge(KEYS)
	t.Logf("Munge output: %s", key)
	ukeys := Unmunge(key)
	l := len(KEYS)
	if len(ukeys) != l {
		t.Errorf("Lengths don't match - got %d, expected %d", len(ukeys), l)
		for _, q := range ukeys {
			t.Logf("%s", q)
		}
		return
	}
	for i := 0; i < l; i++ {
		if ukeys[i] != KEYS[i] {
			t.Errorf("Key %s doesn't match expected %s", ukeys[i], KEYS[i])
			return
		}
	}
}
