package vtag

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/Masterminds/semver/v3"
	"github.com/subtlepseudonym/version"
)

const tagRevision = `tags.*\^0`

// GitBinaryVersion retrieves a semantic version from the local git repo's tags,
// preferring to use the branch name if useBranch is true and the branch is
// a parseable semantic version with prerelease tags
func GitBinaryVersion(useBranch bool) (string, error) {
	if useBranch {
		prerelease := checkPreReleaseBranch()
		if prerelease != "" {
			return prerelease, nil
		}
	}

	ver, err := version.Latest(version.GitBinary, ".")
	if err != nil {
		return "", fmt.Errorf("get latest version: %v", err)
	}

	return ver, nil
}

// GitBinaryBuildTag generates a build tag based upon the current git revision and
// the state of the current work tree
func GitBinaryBuildTag() (string, error) {
	nameRev := exec.Command("git", "name-rev", "--name-only", "HEAD")
	nameRevOut, err := nameRev.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git name-rev: %v", err)
	}
	revisionName := bytes.Trim(nameRevOut, "\n")

	var revision []byte
	if !regexp.MustCompile(tagRevision).Match(revisionName) {
		revList := exec.Command("git", "rev-list", "-n1", "--abbrev-commit", "HEAD")
		revListOut, err := revList.CombinedOutput()
		if err != nil {
			return "", fmt.Errorf("git rev-list: %v", err)
		}
		revision = bytes.Trim(revListOut, "\n")
	}

	var added, deleted int
	diffStat := exec.Command("git", "diff-files", "--numstat")
	diffStatOut, err := diffStat.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff-files: %v: %s", err, bytes.Trim(diffStatOut, "\n"))
	}

	for _, line := range strings.Split(string(diffStatOut), "\n") {
		if line == "" {
			continue
		}
		stats := strings.Split(line, "\t")

		addCount, err := strconv.Atoi(stats[0])
		if err != nil {
			return "", fmt.Errorf("parse int %q: %w", stats[0], err)
		}
		added += addCount

		delCount, err := strconv.Atoi(stats[1])
		if err != nil {
			return "", fmt.Errorf("parse int %q: %w", stats[1], err)
		}
		deleted += delCount
	}

	var updated int
	lsFiles := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	lsFilesOut, err := lsFiles.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git ls-files: %v: %s", err, bytes.Trim(lsFilesOut, "\n"))
	}
	for _, filename := range strings.Split(string(lsFilesOut), "\n") {
		if filename == "" {
			continue
		}

		f, err := os.Open(filename)
		if err != nil {
			return "", fmt.Errorf("open file: %w", err)
		}

		lineCount, err := countLines(f)
		if err != nil {
			return "", fmt.Errorf("count lines: %w", err)
		}
		updated += lineCount
	}

	return BuildTag(revision, added, deleted, updated), nil
}

// checkPreReleaseBranch parses the current branch name to check if it matches
// the semantic version spec and has a prerelease version
//
// The intention here is to generate version strings for bleeding edge prerelease
// versions directly from the prerelease development branch.
func checkPreReleaseBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// FIXME: do something with this error (even if it's just logging on debug mode)
		return ""
	}
	branch := strings.Trim(string(out), "\n")

	v, err := semver.NewVersion(branch)
	if err == nil && v.Prerelease() != "" {
		v.SetMetadata("") // zero out build data
		return v.String()
	}

	// FIXME: do something with this error
	return ""
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
