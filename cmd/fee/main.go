package main

import (
	"fmt"
	"math/big"
	"os"
	"strconv"

	"github.com/spf13/cobra"
)

var Version = "1.2.0"

var (
	// flags
	totalFloat float64
	feeFloat   float64
	precision  int
)

func main() {
	cmd := cobra.Command{
		Use:     "fee {--total|--fee} charges...",
		Short:   "Calculate proportion of fee for each charge in set",
		Version: Version,
		RunE:    run,
	}

	cmd.Flags().Float64VarP(&totalFloat, "total", "t", 0.0, "total of charges plus fee")
	cmd.Flags().Float64VarP(&feeFloat, "fee", "f", 0.0, "fee applied to sum of charges. If both this flag and --total are provided, a missing charge is inferred if the sum of charges does not equal (total + fee + charges)")
	cmd.Flags().IntVarP(&precision, "precision", "p", 2, "numeric precision of output")

	cmd.ParseFlags(os.Args[1:])

	err := cmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// pEq returns a boolean indicating whether two floats are equal at a given decimal precision
func pEq(x, y *big.Float, precision int) bool {
	return x.Text('f', precision) == y.Text('f', precision)
}

func run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("no charges provided")
	}

	if feeFloat == 0.0 && totalFloat == 0.0 {
		return fmt.Errorf("either --fee or --total must be provided")
	}

	total := big.NewFloat(totalFloat)
	fee := big.NewFloat(feeFloat)

	sum := new(big.Float)
	var charges []*big.Float
	for _, arg := range args {
		c, err := strconv.ParseFloat(arg, 64)
		if err != nil {
			return fmt.Errorf("parse float: %w", err)
		}

		charge := big.NewFloat(c)
		sum.Add(sum, charge)
		charges = append(charges, charge)
	}

	if total.Sign() > 0 {
		inferred := total.Sub(total, sum)
		if !pEq(inferred, new(big.Float), precision) {
			if fee.Sign() > 0 {
				inferred.Sub(inferred, fee)
				sum.Add(sum, inferred)
				charges = append(charges, inferred)
			} else {
				fee = inferred
			}
		}

		if inferred.Sign() < 0 {
			fmt.Fprint(os.Stderr, "warn: inferred charge is negative\n")
		}
	}

	if fee.Sign() < 0 {
		fmt.Fprint(os.Stderr, "warn: fee is negative\n")
	}

	proportionSum := new(big.Float)
	for _, charge := range charges {
		proportion := new(big.Float).Mul(fee, new(big.Float).Quo(charge, sum))
		proportionSum.Add(proportionSum, proportion)
		outputFormat := fmt.Sprintf("%%.%df + %%.%df = %%.%df\n", precision, precision, precision)
		fmt.Fprintf(os.Stdout, outputFormat, charge, proportion, new(big.Float).Add(charge, proportion))
	}

	if !pEq(fee, proportionSum, precision) {
		format := fmt.Sprintf("warn: remainder: %%.%df\n", precision)
		fmt.Fprintf(os.Stderr, format, fee.Sub(fee, proportionSum))
	}

	return nil
}
