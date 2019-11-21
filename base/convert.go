package base

import (
	"fmt"
	"math/big"
)

func Convert(input string, from, to int) (out string, err error) {
	var inputSet bool

	defer func() {
		if r := recover(); r != nil {
			if !inputSet {
				err = fmt.Errorf("from base: %v", r)
			} else {
				err = fmt.Errorf("to base: %v", r)
			}
		}
	}()

	i := new(big.Int)

	i, _ = i.SetString(input, from)
	inputSet = true

	out = i.Text(to)
	return out, nil
}
