package bloom

import (
	"fmt"
	"hash/fnv"
	"testing"
)

func Test_bloom(t *testing.T) {
	fnv128 := fnv.New128a()
	fnv128.Write([]byte("hello hi hello hello hello"))

	hash := fnv128.Sum(nil)
	fmt.Println(hash)
}
