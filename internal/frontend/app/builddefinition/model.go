package builddefinition

import (
	stdslices "golang.org/x/exp/slices"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/common/maybe"
)

type Definition struct {
	Vertexes []api.Vertex
	Vars     []api.Var
}

func (d Definition) Vertex(name string) maybe.Maybe[api.Vertex] {
	i := stdslices.IndexFunc(d.Vertexes, func(vertex api.Vertex) bool {
		return vertex.Name == name
	})
	if i == -1 {
		return maybe.Maybe[api.Vertex]{}
	}

	return maybe.NewJust(d.Vertexes[i])
}
