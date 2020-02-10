package main

import (
	"fmt"
	"os"

	"github.com/subtlepseudonym/utilities/vtag"

	"github.com/urfave/cli"
)

const (
	gitBinary = "binary"
	goGit     = "go-git"
	libGit2   = "libgit2"

	branchRegex  = `.*-rc`
	tagNameRegex = `tags.*\^0`
)

var Version = "1.0.0"

func main() {
	app := cli.NewApp()

	app.Name = "vtag"
	app.Version = Version
	app.Usage = "generate version tag for code releases"

	app.Writer = os.Stdout
	app.ErrWriter = os.Stderr

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "ignore-rc",
			Usage: "don't use \"*-rc\" branch names as version",
		},
		cli.StringFlag{
			Name:  "git-lib",
			Usage: "method for retrieving git repo info",
			Value: gitBinary,
		},
	}

	app.Action = versionTag

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
}

func versionTag(ctx *cli.Context) error {
	lib := ctx.String("git-lib")
	switch lib {
	case gitBinary:
		return cliAction(ctx)
	case goGit:
		return goGitAction(ctx)
	case libGit2:
		return libGit2Action(ctx)
	default:
		_, err := fmt.Fprintf(ctx.App.ErrWriter, "unrecognized git lib %q: use %q, %q, %q", lib, gitBinary, goGit, libGit2)
		return err
	}
}

func cliAction(ctx *cli.Context) error {
	version, err := vtag.GitBinaryVersion(!ctx.Bool("ignore-rc"))
	if err != nil {
		return fmt.Errorf("version: %v", err)
	}

	buildTag, err := vtag.GitBinaryBuildTag()
	if err != nil {
		return fmt.Errorf("build tag: %v", err)
	}

	_, err = fmt.Fprintf(ctx.App.Writer, "%s%s", version, buildTag)
	return err
}

func goGitAction(ctx *cli.Context) error {
	_, err := fmt.Fprintf(ctx.App.Writer, "not implemented")
	return err
}

func libGit2Action(ctx *cli.Context) error {
	_, err := fmt.Fprintf(ctx.App.Writer, "not implemented")
	return err
}
