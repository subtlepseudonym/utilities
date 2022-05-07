package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/scru128/go-scru128"
	"github.com/spf13/cobra"
)

var Version = "1.0.0"

func main() {
	cmd := cobra.Command{
		Use:       "scru [count]",
		Short:     "Generate SCRU128 IDs",
		ValidArgs: []string{"count"},
		Version:   Version,
		RunE:      run,
	}

	err := cmd.Execute()
	if err != nil {
		fmt.Println("ERR:", err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	var err error
	count := 1
	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("parse count arg: %w", err)
		}
	}

	generator := scru128.NewGenerator()
	for i := 0; i < count; i++ {
		id, err := generator.Generate()
		if err != nil {
			panic(err)
		}

		fmt.Println(id)
	}

	return nil
}
