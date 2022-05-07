package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var Version = "1.0.0"

var (
	// flags
	total     float64
	fee       float64
	precision int
)

func main() {
	cmd := cobra.Command{
		Use:     "fee {--total|--fee} charges...",
		Short:   "Calculate proportion of fee for each charge in set",
		Version: Version,
		RunE:    run,
	}

	cmd.Flags().Float64VarP(&total, "total", "t", 0.0, "total of charges plus fee")
	cmd.Flags().Float64VarP(&fee, "fee", "f", 0.0, "fee applied to sum of charges. If both this flag and --total are provided, this flag is used")
	cmd.Flags().IntVarP(&precision, "precision", "p", 2, "numeric precision of output")

	cmd.ParseFlags(os.Args[1:])

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no charges provided")
	}

	var sum float64
	var charges []float64
	for _, arg := range args {
		charge, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return fmt.Errorf("parse float: %w", err)
		}

		sum += charge
		charges = append(charges, charge)
	}

	if fee == 0.0 {
		if total == 0.0 {
			return fmt.Errorf("either --total or --fee must be provided")
		}

		fee = total - sum
	}

	if fee < 0 {
		fmt.Fprint(os.Stderr, "warn: fee is negative")
	}

	for _, charge := range charges {
		proportion := fee * (charge / sum)
		outputFormat := fmt.Sprintf("%%.%df %%.%df\n", precision, precision)
		fmt.Fprintf(os.Stdout, outputFormat, charge, proportion)
	}

	return nil
}
