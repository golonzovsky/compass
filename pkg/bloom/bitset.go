package bloom

const (
	wordSize     = 64
	log2WordSize = 6
)

//todo should we just wrap big.Int instead?

type BitSet struct {
	length uint
	set    []uint64
}

func NewBitSet(length uint) *BitSet {
	return &BitSet{
		length: length,
		set:    make([]uint64, wordsNeededFor(length)),
	}
}

func wordsNeededFor(i uint) int {
	//todo check for max cap
	return int((i + wordSize - 1) >> log2WordSize)
}

func indexInWord(i uint) uint {
	return i & (wordSize - 1)
}

func maxCap() uint {
	return ^uint(0)
}

func (b *BitSet) Set(i uint) {
	if i >= b.length {
		b.extendToFit(i)
	}
	b.set[i>>log2WordSize] |= 1 << indexInWord(i)
}

func (b *BitSet) Clear(i uint) {
	b.set[i>>log2WordSize] &^= 1 << indexInWord(i)
}

func (b *BitSet) SetTo(i uint, value bool) {
	if value {
		b.Set(i)
	}
	b.Clear(i)
}

func (b *BitSet) Get(i uint) bool {
	if i > b.length {
		return false
	}
	return b.set[i>>log2WordSize]&(1<<indexInWord(i)) != 0
}

func (b *BitSet) extendToFit(i uint) {
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
