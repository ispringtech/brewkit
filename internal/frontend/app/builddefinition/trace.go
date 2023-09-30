package builddefinition

import (
	"fmt"
	"strings"

	stdslices "golang.org/x/exp/slices"

	"github.com/ispringtech/brewkit/internal/common/slices"
)

const (
	from          = "from"
	deps          = "deps"
	copyDirective = "copy"
)

type traceEntry struct {
	name      string
	directive string
}

type trace []traceEntry

func (t *trace) push(s traceEntry) {
	*t = append(*t, s)
}

func (t *trace) pop() {
	*t = (*t)[:len(*t)-1]
}

func (t *trace) has(s string) bool {
	return stdslices.ContainsFunc(*t, func(entry traceEntry) bool {
		return entry.name == s
	})
}

func (t *trace) String() string {
	if len(*t) == 0 {
		return ""
	}

	res := make([]traceEntry, 0, len(*t))
	for i := len(*t) - 1; i >= 0; i-- {
		res = append(res, (*t)[i])
	}

	return strings.Join(slices.Map(res, func(e traceEntry) string {
		return fmt.Sprintf("%s(%s)", e.name, e.directive)
	}), "->")
}
