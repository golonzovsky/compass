package bloom

import (
	"hash/fnv"
	"math"

	"github.com/twmb/murmur3"
)

// todo eventually switch to github.com/bits-and-blooms/bloom/v3
// but implement here for study purposes

type Filter struct {
	m uint64
	k uint64
	b *BitSet
}

func NewFilter(n uint, p float64) *Filter {
	m := uint64(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
	k := uint64(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return &Filter{
		m: m,
		k: k,
		b: NewBitSet(m),
	}
}

func (f *Filter) keyToPos(key []byte) []uint64 {
	hM := murmur3.Sum64(key)

	fn := fnv.New64()
	fn.Write(key)
	hF := fn.Sum64()

	res := make([]uint64, f.k)
	for i := uint64(0); i < f.k; i++ {
		pos := (hM + i*hF + i*i) % f.m
		res = append(res, pos)
	}
	return res
}

func (f *Filter) Add(key []byte) {
	pos := f.keyToPos(key)
	for _, p := range pos {
		f.b.Set(p)
	}
}

func (f *Filter) Contains(key []byte) bool {
	pos := f.keyToPos(key)
	for _, p := range pos {
		if !f.b.Get(p) {
			return false
		}
	}
	return true
}
