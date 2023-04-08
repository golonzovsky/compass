package bloom

import "testing"

func Test_wordsNeededFor(t *testing.T) {
	cases := []struct {
		num uint
		exp uint
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
		got := wordNeededFor(c.num)

		if got != c.exp {
			t.Errorf("expected %d, got %d", c.exp, got)
		}
	}
}
