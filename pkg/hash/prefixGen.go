package hash

import (
	"fmt"
)

func NewPrefixGen() <-chan string {
	ch := make(chan string)
	go func() {
		for i := 0; i < 0xFFFFF; i++ {
			ch <- fmt.Sprintf("%05x", i)
		}
		close(ch)
	}()
	return ch
}
