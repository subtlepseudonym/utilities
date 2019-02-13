package main

import (
	"fmt"
	"math/big"
	"os"

	"github.com/urfave/cli"
)

var Version = "v0.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "base"
	app.Version = Version
	app.Usage = "utility for converting numbers between different bases"

	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "from",
			Usage: "base of input",
			Value: 10,
		},

		cli.IntFlag{
			Name:  "to",
			Usage: "base of output",
			Value: 2,
		},

		cli.StringFlag{
			Name:  "pad",
			Usage: "what to pad output with",
			Value: "0",
		},
	}

	app.Action = convertBases

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func convertBases(ctx *cli.Context) error {
	var nums []*big.Int
	for _, arg := range ctx.Args() {
		argNum, ok := big.NewInt(0).SetString(arg, ctx.Int("from"))
		if !ok {
			return cli.NewExitError(fmt.Errorf("unable to parse %q", arg), 1)
		}

		nums = append(nums, argNum)
	}

	var out []string
	for _, num := range nums {
		out = append(out, num.Text(ctx.Int("to")))
	}

	var longest int
	for _, str := range out {
		if len(str) > longest {
			longest = len(str)
		}
	}

	for _, str := range out {
		for i := 0; i < longest-len(str); i++ {
			fmt.Print(ctx.String("pad"))
		}

		fmt.Fprintln(ctx.App.Writer, str)
	}

	return nil
}
