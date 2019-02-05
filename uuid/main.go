package main

import (
	"flag"
	"os"

	"github.com/satori/go.uuid"
)

func main() {
	n := flag.Int("n", 1, "number of UUID to generate")
	flag.Parse()

	for i := 0; i < *n; i++ {
		id, err := uuid.NewV4()
		if err != nil {
			_, e := os.Stderr.WriteString("generate uuid failed: " + err.Error())
			if err != nil {
				panic(e)
			}
			os.Exit(1)
		}

		os.Stdout.WriteString(id.String() + "\n")
	}
}
