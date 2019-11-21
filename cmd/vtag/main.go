package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

const (
	cliGit  = "cli"
	goGit   = "go-git"
	libGit2 = "libgit2"

	branchRegex  = `.*-rc`
	tagNameRegex = `tags.*\^0`
)

var Version = "0.0.0"

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
			Value: "cli",
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
	case cliGit:
		return cliAction(ctx)
	case goGit:
		return goGitAction(ctx)
	case libGit2:
		return libGit2Action(ctx)
	default:
		_, err := fmt.Fprintf(ctx.App.ErrWriter, "unrecognized git lib %q: use %q, %q, %q", lib, cliGit, goGit, libGit2)
		return err
	}
}

func buildTag(shortRevision []byte, added, deleted, updated int) string {
	builder := new(strings.Builder)
	var buildAdded bool
	var changesAdded bool

	if len(shortRevision) > 0 {
		fmt.Fprintf(builder, "+%s", shortRevision)
		buildAdded = true
	}

	builder, buildAdded, changesAdded = addChange(builder, 'a', added, buildAdded, changesAdded)
	builder, buildAdded, changesAdded = addChange(builder, 'd', deleted, buildAdded, changesAdded)
	builder, buildAdded, changesAdded = addChange(builder, 'u', updated, buildAdded, changesAdded)

	return builder.String()
}

func addChange(builder *strings.Builder, tag rune, count int, build, changes bool) (*strings.Builder, bool, bool) {
	if count == 0 {
		return builder, build, changes
	}

	if !build {
		builder.WriteRune('+')
	} else if !changes {
		builder.WriteRune('.')
	}
	fmt.Fprintf(builder, "%c%d", tag, count)

	return builder, true, true
}

func countLines(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func cliAction(ctx *cli.Context) error {
	versionTagCmd := exec.Command("git", "describe", "--abbrev=0")
	tag, err := versionTagCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("get version: %w: %s", err, bytes.Trim(tag, "\n"))
	}
	version := tag

	if !ctx.Bool("ignore-rc") {
		branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
		branch, err := branchCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("get branch: %w: %s", err, bytes.Trim(branch, "\n"))
		}

		if regexp.MustCompile(branchRegex).Match(branch) {
			version = branch
		}
	}
	version = bytes.Trim(version, "\n")

	revNameCmd := exec.Command("git", "name-rev", "--name-only", "HEAD")
	revName, err := revNameCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("get rev name: %w: %s", err, bytes.Trim(revName, "\n"))
	}

	var revision []byte
	if !regexp.MustCompile(tagNameRegex).Match(revName) {
		shortRevCmd := exec.Command("git", "rev-list", "-n1", "--abbrev-commit", "HEAD")
		shortRev, err := shortRevCmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("get short rev: %w: %s", err, bytes.Trim(shortRev, "\n"))
		}
		revision = bytes.Trim(shortRev, "\n")
	}

	var added, deleted int
	diffStatCmd := exec.Command("git", "diff-files", "--numstat")
	diffStatOutput, err := diffStatCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("get diff stat: %w: %s", err, bytes.Trim(diffStatOutput, "\n"))
	}

	for _, line := range strings.Split(string(diffStatOutput), "\n") {
		if line == "" {
			continue
		}
		stats := strings.Split(line, "\t")

		addCount, err := strconv.Atoi(stats[0])
		if err != nil {
			return fmt.Errorf("parse int %q: %w", stats[0], err)
		}
		added += addCount

		delCount, err := strconv.Atoi(stats[1])
		if err != nil {
			return fmt.Errorf("parse int %q: %w", stats[1], err)
		}
		deleted += delCount
	}

	var updated int
	lsFilesCmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	lsFilesOutput, err := lsFilesCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("get untracked: %w: %s", err, bytes.Trim(lsFilesOutput, "\n"))
	}
	for _, filename := range strings.Split(string(lsFilesOutput), "\n") {
		if filename == "" {
			continue
		}

		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("open file: %w", err)
		}

		lineCount, err := countLines(f)
		if err != nil {
			return fmt.Errorf("count lines: %w", err)
		}
		updated += lineCount
	}

	_, err = fmt.Fprintf(ctx.App.Writer, "%s%s", version, buildTag(revision, added, deleted, updated))
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
