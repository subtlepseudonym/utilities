package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/satori/go.uuid"
	"github.com/urfave/cli"
)

var Version = "1.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "uuid"
	app.Version = Version
	app.Usage = "generates v4 UUIDs"
	app.UsageText = fmt.Sprintf("%s [flags] [count]", app.Name)
	app.Description = `Generates UUIDs. If you provide an integer argument, it will treat it like usage of the --count flag. Explicit usage of the --count flag will override this behavior.`

	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "count, n",
			Usage: "number of UUIDs to generate",
			Value: 1,
		},
	}

	app.Before = processArgs
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

func processArgs(ctx *cli.Context) error {
	if ctx.NArg() == 0 || ctx.IsSet("count") {
		return nil
	}

	countArg := ctx.Args()[0]
	_, err := strconv.Atoi(countArg)
	if err != nil {
		return nil
	}

	ctx.Set("count", countArg)

	return nil
}

func action(ctx *cli.Context) error {
	for i := 0; i < ctx.Int("count"); i++ {
		id, err := uuid.NewV4()
		if err != nil {
			return fmt.Errorf("new uuid: %v", err)
		}

		fmt.Fprintln(ctx.App.Writer, id.String())
	}

	return nil
}
