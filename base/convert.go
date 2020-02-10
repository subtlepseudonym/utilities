package base

import (
	"fmt"
	"math/big"
)

func Convert(input string, from, to int) (out string, err error) {
	var inputOk bool

	defer func() {
		if r := recover(); r != nil {
			if !inputOk {
				err = fmt.Errorf("from base: %v", r)
			} else {
				err = fmt.Errorf("to base: %v", r)
			}
		}
	}()

	i := new(big.Int)
	i, inputOk = i.SetString(input, from)

	return i.Text(to), nil
}
