package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/gofrs/uuid/v5"
	"github.com/urfave/cli"
)

var Version = "1.1.0"

var (
	uuidType  byte // uuid.V1, uuid.V3, etc
	namespace uuid.UUID
)

func main() {
	app := cli.NewApp()

	app.Name = "uuid"
	app.Version = Version
	app.Usage = "generates v4 UUIDs"
	app.UsageText = fmt.Sprintf("%s [flags] [count]", app.Name)
	app.Description = `Generates UUIDs.

	Supported UUID types:
		- [1|v1]:            Version 1 (date-time and MAC address)
		- [3|v3|md5]:        Version 3 (namespace name-based MD5)
		- [4|v4|random]:     Version 4 (random)
		- [5|v5|sha1|sha-1]: Version 5 (namespace name-based SHA-1)
		- [6|v6]:            Version 6 (k-sortable and random, field-compatible with v1)
		- [7|v7]:            Version 7 (k-sortable and random)
`

	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	app.Flags = []cli.Flag{
		cli.IntFlag{
			Name:  "count, n",
			Usage: "number of UUIDs to generate",
			Value: 1,
		},
		cli.StringFlag{
			Name:  "type, t",
			Usage: "type of UUID",
			Value: "v4",
		},
		cli.StringFlag{
			Name:  "namespace",
			Usage: "namespace for UUID v3 or v5",
		},
		cli.StringFlag{
			Name:  "name",
			Usage: "name for UUID v3 or v5",
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

func processArgs(ctx *cli.Context) (err error) {
	if !ctx.IsSet("count") && ctx.NArg() > 0 {
		countArg := ctx.Args()[0]
		_, err = strconv.Atoi(countArg)
		if err != nil {
			return fmt.Errorf("parse count arg: %w", err)
		}
		ctx.Set("count", countArg)
	}

	switch strings.ToLower(ctx.String("type")) {
	case "1", "v1":
		uuidType = uuid.V1
	case "3", "v3", "md5":
		uuidType = uuid.V3
	case "4", "v4", "random":
		uuidType = uuid.V4
	case "5", "v5", "sha1", "sha-1":
		uuidType = uuid.V5
	case "6", "v6":
		uuidType = uuid.V6
	case "7", "v7":
		uuidType = uuid.V7
	default:
		return fmt.Errorf("unknown UUID type: %q", ctx.String("type"))
	}

	if uuidType == uuid.V3 || uuidType == uuid.V5 {
		namespace, err = uuid.FromString(ctx.String("namespace"))
		if err != nil {
			return fmt.Errorf("parse namespace UUID: %w", err)
		}

		if !ctx.IsSet("name") {
			return fmt.Errorf("name is required for UUID v3 and v5")
		}
	}

	return nil
}

func action(ctx *cli.Context) error {
	for i := 0; i < ctx.Int("count"); i++ {
		var id uuid.UUID
		var err error
		switch uuidType {
		case uuid.V1:
			id, err = uuid.NewV1()
		case uuid.V3:
			id = uuid.NewV3(namespace, ctx.String("name"))
		case uuid.V4:
			id, err = uuid.NewV4()
		case uuid.V5:
			id = uuid.NewV5(namespace, ctx.String("name"))
		case uuid.V6:
			id, err = uuid.NewV6()
		case uuid.V7:
			id, err = uuid.NewV7()
		default:
			// the old "this should never happen"
			return fmt.Errorf("internal: unknown uuid type: %#v", uuidType)
		}

		if err != nil {
			return fmt.Errorf("new uuid: %v", err)
		}
		fmt.Fprintln(ctx.App.Writer, id.String())
	}

	return nil
}
