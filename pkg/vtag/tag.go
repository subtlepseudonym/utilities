package vtag

import (
	"fmt"
	"strings"
)

// BuildTag generates a build tag using a shortened git revision and the state of
// the current work tree
func BuildTag(shortRevision []byte, added, deleted, updated int) string {
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
