package bloom

import "math/bits"

const (
	wordSize     = 64
	log2WordSize = 6
)

// todo should we just wrap big.Int instead?

type BitSet struct {
	length uint64
	set    []uint64
}

func NewBitSet(length uint64) *BitSet {
	return &BitSet{
		length: length,
		set:    make([]uint64, wordsNeededFor(length)),
	}
}

func wordsNeededFor(i uint64) int {
	if i > maxCap() {
		panic("bitset max cap exceeded")
	}
	return int((i + wordSize - 1) >> log2WordSize)
}

func indexInWord(i uint64) uint64 {
	return i & (wordSize - 1)
}

func maxCap() uint64 {
	return ^uint64(0)
}

func (b *BitSet) Set(i uint64) {
	if i >= b.length {
		b.extendToFit(i)
	}
	b.set[i>>log2WordSize] |= 1 << indexInWord(i)
}

func (b *BitSet) Clear(i uint64) {
	b.set[i>>log2WordSize] &^= 1 << indexInWord(i)
}

func (b *BitSet) SetTo(i uint64, value bool) {
	if value {
		b.Set(i)
	}
	b.Clear(i)
}

func (b *BitSet) Get(i uint64) bool {
	if i > b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<indexInWord(i)) != 0
}

func (b *BitSet) extendToFit(i uint64) {
	if i >= maxCap() {
		panic("bitset max cap exceeded")
	}
	if i < b.length {
		return
	}
	nWords := wordsNeededFor(i + 1)
	if b.set == nil {
		b.set = make([]uint64, nWords)
	} else if cap(b.set) >= nWords {
		b.set = b.set[:nWords]
	} else if len(b.set) < nWords {
		newBits := make([]uint64, nWords, 2*nWords) // increase capacity 2x
		copy(newBits, b.set)
		b.set = newBits
	}
	b.length = i + 1
}

func (b *BitSet) CountOnes() uint {
	if b == nil && b.set == nil {
		return 0
	}

	var cnt int
	for _, x := range b.set {
		cnt += bits.OnesCount64(x)
	}
	return uint(cnt)
}
