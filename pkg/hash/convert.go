package hash

import (
	"encoding/hex"
	"fmt"
)

func FromPrefix(prefix string) ([]byte, error) {
	if len(prefix) != 5 {
		return nil, fmt.Errorf("prefix must be 5 characters long")
	}
	return hex.DecodeString("0" + prefix)
}
