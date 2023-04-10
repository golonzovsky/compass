package bloom

import (
	"fmt"
	"strconv"
	"testing"
)

func Test_bloom(t *testing.T) {
	bloom := NewFilter(50, 0.001)

	keys := []string{"quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
	for _, k := range keys {
		bloom.Add([]byte(k))
	}

	for _, k := range keys {
		if !bloom.Contains([]byte(k)) {
			t.Errorf("expected %s to be in the filter", k)
		}
	}

	fmt.Println("bloom filter ApproximatedSize:", bloom.ApproximatedSize())

	for i := 0; i < 10000; i++ {
		if bloom.Contains([]byte(strconv.Itoa(i))) {
			t.Errorf("expected %d to be in the filter", i)
		}
	}
}

func Test_bloomBig(t *testing.T) {
	bloom := NewFilter(600000000, 0.001)
	keys := []string{"quick", "brown", "fox", "jumps", "over", "the", "lazy", "dog"}
	for _, k := range keys {
		bloom.Add([]byte(k))
	}

	for _, k := range keys {
		if !bloom.Contains([]byte(k)) {
			t.Errorf("expected %s to be in the filter", k)
		}
	}
}
