package bloom

import (
	"math/bits"
)

const (
	wordSize     = 64
	log2WordSize = 6
)

//todo should we just wrap big.Int instead?

type BitSet struct {
	length uint
	set    []uint64
}

func wordNeededFor(i uint) uint {
	//todo check for max cap
	return (i + wordSize - 1) >> log2WordSize
}

func wordIndex(i uint) uint {

}

func maxCap() uint {
	return bits.UintSize
}

func (b *BitSet) Set(p uint) {
	// todo
}

func (b *BitSet) Get(i uint) bool {
	if i > b.length {
		return false
	}
	//todo
	return b.set[i>>log2WordSize]&(1<<)
}
