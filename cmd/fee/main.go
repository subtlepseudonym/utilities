package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var Version = "1.0.0"

var (
	total float64
	fee float64
)

func main() {
	cmd := cobra.Command{
		Use:     "fee {--total|--fee} charges...",
		Short:   "Calculate proportion of fee for each charge in total",
		Long:    "Calculate the difference between the sum of all charges and the sum of all charges plus a fee and display the proportion of the fee for which each charge is responsible. For meaningful results, either the --total or --fee flag must be used.",
		Version: Version,
		RunE:    run,
	}

	cmd.Flags().Float64VarP(&total, "total", "t", 0.0, "Total of charges plus fee")
	cmd.Flags().Float64VarP(&fee, "fee", "f", 0.0, "Fee to use in calculating total charge. If both this flag and --total are provided, this flag is used")

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
		fmt.Fprintf(os.Stdout, "%f %f\n", charge, proportion)
	}

	return nil
}
