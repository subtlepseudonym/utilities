package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/spf13/cobra"
)

var Version = "1.3.0"

var (
	// flags
	totalFlag string
	feeFlag   string
	precision int
)

func main() {
	cmd := cobra.Command{
		Use:     "fee {--total|--fee} charges...",
		Short:   "Calculate proportion of fee for each charge in set",
		Version: Version,
		RunE:    run,
	}

	cmd.Flags().StringVarP(&totalFlag, "total", "t", "", "total of charges plus fee")
	cmd.Flags().StringVarP(&feeFlag, "fee", "f", "", "fee applied to sum of charges. If both this flag and --total are provided, a missing charge is inferred if the sum of charges does not equal (total + fee)")
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

	if feeFlag == "" && totalFlag == "" {
		return fmt.Errorf("either --fee or --total must be provided")
	}

	total, ok := new(big.Rat).SetString(totalFlag)
	if !ok {
		return fmt.Errorf("parse rat: %q", totalFlag)
	}

	fee, ok := new(big.Rat).SetString(feeFlag)
	if !ok {
		return fmt.Errorf("parse rat: %q", totalFlag)
	}

	sum := new(big.Rat)
	var charges []*big.Rat
	for _, arg := range args {
		charge, ok := new(big.Rat).SetString(arg)
		if !ok {
			return fmt.Errorf("parse rat: %q", arg)
		}
		sum.Add(sum, charge)
		charges = append(charges, charge)
	}

	var warnings []string
	if total.Cmp(new(big.Rat)) != 0 {
		inferred := total.Sub(total, new(big.Rat).Add(fee, sum))
		cmp := inferred.Cmp(new(big.Rat))
		if cmp != 0 {
			if fee.Cmp(new(big.Rat)) != 0 {
				sum.Add(sum, inferred)
				charges = append(charges, inferred)
			} else {
				fee = inferred
			}
		}

		if cmp < 0 {
			warnings = append(warnings, "inferred charge is negative")
		}
	}
	total = new(big.Rat).Add(sum, fee)

	if fee.Sign() < 0 {
		warnings = append(warnings, "fee is negative")
	}

	lineSum := new(big.Rat)
	for _, charge := range charges {
		proportion := new(big.Rat).Quo(new(big.Rat).Mul(fee, charge), sum)
		lineTotal := new(big.Rat).Add(charge, proportion).FloatString(precision)
		fmt.Fprintf(
			os.Stdout,
			"%s + %s = %s\n",
			charge.FloatString(precision),
			proportion.FloatString(precision),
			lineTotal,
		)

		rounded, _ := new(big.Rat).SetString(lineTotal)
		lineSum.Add(lineSum, rounded)
	}

	if total.Cmp(lineSum) != 0 {
		remainder := new(big.Rat).Sub(total, lineSum)
		n, exact := remainder.FloatPrec()
		var format string
		if exact {
			format = fmt.Sprintf("remainder %%.%df", n)
		} else {
			format = fmt.Sprintf("remainder %%.%df", n+2)
		}

		remainderFloat, _ := remainder.Float64()
		warnings = append(warnings, fmt.Sprintf(format, remainderFloat))
	}

	for _, warn := range warnings {
		fmt.Fprintf(os.Stderr, "warn: %s\n", warn)
	}

	return nil
}
