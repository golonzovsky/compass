package bloom

import (
	"testing"
)

func Test_wordsNeededFor(t *testing.T) {
	cases := []struct {
		num uint64
		exp int
	}{
		{
			1,
			1,
		},
		{
			70, //64+
			2,
		},
	}

	for _, c := range cases {
		got := wordsNeededFor(c.num)

		if got != c.exp {
			t.Errorf("expected %d, got %d", c.exp, got)
		}
	}
}

func Test_bitSet(t *testing.T) {
	bs := BitSet{}
	for i := uint64(0); i < 10000; i += 2 {
		bs.Set(i)
	}

	for i := uint64(0); i < 10000; i++ {
		bit := bs.Get(i)
		if i%2 == 0 && !bit {
			t.Errorf("expected %d to be set", i)
		}
		if i%2 == 1 && bit {
			t.Errorf("expected %d to be unset", i)
		}
	}
}
