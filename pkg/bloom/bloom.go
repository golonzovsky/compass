package bloom

type Filter struct {
	b *BitSet

	hf []hashGen
}

func NewFilter(bits, hashes uint) *Filter {
	return &Filter{
		b:  NewBitSet(bits),
		hf: initHashFuncs(hashes, bits),
	}
}

func (f *Filter) Add([]byte) {

}

func (f *Filter) Test([]byte) bool {
	return true
}

type hashGen func(h1, h2 uint) uint

//func (f *Filter) keyToPos(key []byte) uint {
//	hM := murmur3.New128().Sum(key)
//	hF := fnv.New128().Sum(key)
//	for _, f := range f.hf {
//		f(hM, hF)
//	}
//	return 0
//}

func initHashFuncs(num, bits uint) []hashGen {
	var funcs []hashGen
	for i := uint(0); i < num; i++ {
		funcs = append(funcs, func(h1, h2 uint) uint {
			return (h1 + i*h2 + i*i) % bits
		})
	}
	return funcs
}
