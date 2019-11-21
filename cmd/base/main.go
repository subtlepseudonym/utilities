package main

import (
	"fmt"
	"os"

	"github.com/subtlepseudonym/utilities/base"

	"github.com/urfave/cli"
)

var Version = "1.1.0"

func main() {
	app := cli.NewApp()

	app.Name = "base"
	app.Version = Version
	app.Usage = "convert numbers between different bases"

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

	app.Before = validateFlags
	app.Action = action

	app.Authors = []cli.Author{
		{
			Name:  "Connor Demille",
			Email: "subtlepseudonym@gmail.com",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintln(app.ErrWriter, err)
	}
}

func validateFlags(ctx *cli.Context) error {
	pad := ctx.String("pad")
	if len(pad) > 1 {
		return fmt.Errorf("pad flag must be a single character: %q has length %d", pad, len(pad))
	}

	return nil
}

func action(ctx *cli.Context) error {
	var out []string
	for _, arg := range ctx.Args() {
		output, err := base.Convert(arg, ctx.Int("from"), ctx.Int("to"))
		if err != nil {
			return fmt.Errorf("convert: %v", err)
		}

		out = append(out, output)
	}

	padded := padOutput(out, rune(ctx.String("pad")[0]))
	for _, line := range padded {
		fmt.Fprintln(ctx.App.Writer, line)
	}

	return nil
}

func padOutput(out []string, pad rune) []string {
	var longest int
	for _, str := range out {
		if len(str) > longest {
			longest = len(str)
		}
	}

	format := fmt.Sprintf("%%%c%ds", pad, longest)
	padded := make([]string, 0, len(out))
	for _, str := range out {
		padded = append(padded, fmt.Sprintf(format, str))
	}

	return padded
}
