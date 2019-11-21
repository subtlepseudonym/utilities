package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/satori/go.uuid"
)

func main() {
	n := 1
	if len(os.Args) > 0 {
		n64, err := strconv.Atoi(os.Args[1])
		if err != nil {
			_, e := os.Stderr.WriteString(fmt.Sprintf("unable to parse %q", os.Args[1]))
			if e != nil {
				panic(e)
			}
			os.Exit(1)
		}

		n = int(n64)
	}

	for i := 0; i < n; i++ {
		id, err := uuid.NewV4()
		if err != nil {
			_, e := os.Stderr.WriteString("generate uuid failed: " + err.Error())
			if e != nil {
				panic(e)
			}
			os.Exit(1)
		}

		os.Stdout.WriteString(id.String() + "\n")
	}
}
