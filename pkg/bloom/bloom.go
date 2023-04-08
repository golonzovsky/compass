package bloom

type Filter struct {
	b *BitSet
}

func (f *Filter) Add([]byte) {

}

func (f *Filter) Test([]byte) bool {
	return true
}
