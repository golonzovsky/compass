package hash

import (
	"testing"
)

func Test_conversions(t *testing.T) {
	prefix := "a0b0c"
	hash, err := FromPrefix(prefix)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if ToPrefix(hash) != prefix {
		t.Errorf("expected %s, got %s", prefix, ToPrefix(hash))
	}
}
