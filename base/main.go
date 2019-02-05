package main

import (
	"math/big"
	"os"
	"fmt"
)

func main() {
	var nums []*big.Int
	for i := 1; i < len(os.Args); i++ {
		argNum, ok := big.NewInt(0).SetString(os.Args[i], 10)
		if !ok {
			panic(fmt.Errorf("could not parse arg %q", os.Args[i]))
		}
		nums = append(nums, argNum)
	}

	var longest int
	for _, num := range nums {
		if num.BitLen() > longest {
			longest = num.BitLen()
		}
	}

	for _, num := range nums {
		for i := 0; i < longest - num.BitLen(); i++ {
			fmt.Print("0")
		}
		for i := num.BitLen() - 1; i >= 0; i-- {
			fmt.Printf("%d", num.Bit(i))
		}
		fmt.Println()
	}
}
