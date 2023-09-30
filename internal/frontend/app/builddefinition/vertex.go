package builddefinition

import (
	"github.com/pkg/errors"
	stdslices "golang.org/x/exp/slices"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/common/either"
	"github.com/ispringtech/brewkit/internal/common/maps"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
	"github.com/ispringtech/brewkit/internal/frontend/app/buildconfig"
	"github.com/ispringtech/brewkit/internal/frontend/app/config"
)

func newVertexGraphBuilder(secrets []config.Secret, targets []buildconfig.TargetData) *vertexGraphBuilder {
	return &vertexGraphBuilder{
		visitedVertexes: map[string]api.Vertex{},
		vertexesSet: maps.SetFromSlice(targets, func(t buildconfig.TargetData) string {
			return t.Name
		}),
		targetsMap: maps.FromSlice(targets, func(t buildconfig.TargetData) (string, buildconfig.TargetData) {
			return t.Name, t
		}),
		trace:   trace{},
		secrets: secrets,
	}
}

type vertexGraphBuilder struct {
	visitedVertexes map[string]api.Vertex // Set of visited vertexes
	vertexesSet     maps.Set[string]
	targetsMap      map[string]buildconfig.TargetData

	trace   trace // Trace to detect cyclic graphs
	secrets []config.Secret
}

func (builder *vertexGraphBuilder) graphVertexes() ([]api.Vertex, error) {
	vertexes := make([]api.Vertex, 0, len(builder.vertexesSet))
	for vertex := range builder.vertexesSet {
		v, err := builder.recursiveGraph(vertex)
		if err != nil {
			return nil, errors.Wrap(err, "graph solve error")
		}

		vertexes = append(vertexes, v)
	}

	return vertexes, nil
}

func (builder *vertexGraphBuilder) recursiveGraph(vertex string) (api.Vertex, error) {
	if v, found := builder.visitedVertexes[vertex]; found {
		return v, nil
	}

	if builder.trace.has(vertex) {
		return api.Vertex{}, errors.Errorf("recursive graph detected by '%s' target, trace: %s", vertex, builder.trace.String())
	}

	t, ok := builder.targetsMap[vertex]
	if !ok {
		return api.Vertex{}, errors.Errorf("logic error: TargetData for Vertex %s not found", vertex)
	}

	var (
		fromV     maybe.Maybe[*api.Vertex]
		copyDirs  []api.Copy
		dependsOn []api.Vertex
	)

	//nolint:nestif
	if maybe.Valid(t.Stage) {
		stage := maybe.Just(t.Stage)

		// found means target 'From' set as another target
		if builder.vertexesSet.Has(stage.From) {
			v, err := builder.walkFrom(vertex, stage.From)
			if err != nil {
				return api.Vertex{}, err
			}

			fromV = maybe.NewJust(v)
		}

		if len(stage.Copy) != 0 {
			var err error
			copyDirs, err = builder.walkCopy(vertex, stage.Copy)
			if err != nil {
				return api.Vertex{}, err
			}
		}
	}

	if len(t.DependsOn) != 0 {
		var err error
		dependsOn, err = builder.walkDependsOn(vertex, t)
		if err != nil {
			return api.Vertex{}, err
		}
	}

	var stage maybe.Maybe[api.Stage]

	stage, err := maybe.MapErr(t.Stage, func(s buildconfig.StageData) (api.Stage, error) {
		return mapStage(t.Name, maybe.Just(t.Stage), copyDirs, builder.secrets)
	})
	if err != nil {
		return api.Vertex{}, err
	}

	return api.Vertex{
		Name:      vertex,
		Stage:     stage,
		From:      fromV,
		DependsOn: dependsOn,
	}, nil
}

// solves 'from' dependencies
func (builder *vertexGraphBuilder) walkFrom(vertexName, fromVName string) (*api.Vertex, error) {
	builder.trace.push(traceEntry{
		name:      vertexName,
		directive: from,
	})
	defer builder.trace.pop()

	vertex, err := builder.recursiveGraph(fromVName)
	if err != nil {
		return nil, err
	}

	return &vertex, nil
}

// solves 'copy' dependencies
func (builder *vertexGraphBuilder) walkCopy(vertexName string, copyDirs []buildconfig.Copy) ([]api.Copy, error) {
	builder.trace.push(traceEntry{
		name:      vertexName,
		directive: copyDirective,
	})
	defer builder.trace.pop()

	return slices.MapErr(copyDirs, func(c buildconfig.Copy) (api.Copy, error) {
		if !maybe.Valid(c.From) {
			return api.Copy{
				Src: c.Src,
				Dst: c.Dst,
			}, nil
		}

		copyFrom := maybe.Just(c.From)

		if !builder.vertexesSet.Has(copyFrom) {
			return api.Copy{
				Src:  c.Src,
				Dst:  c.Dst,
				From: maybe.NewJust(either.NewRight[*api.Vertex, string](copyFrom)),
			}, nil
		}

		vertex, err := builder.recursiveGraph(copyFrom)
		if err != nil {
			return api.Copy{}, err
		}

		return api.Copy{
			Src:  c.Src,
			Dst:  c.Dst,
			From: maybe.NewJust(either.NewLeft[*api.Vertex, string](&vertex)),
		}, nil
	})
}

// solves 'dependsOn' dependencies
func (builder *vertexGraphBuilder) walkDependsOn(vertexName string, t buildconfig.TargetData) ([]api.Vertex, error) {
	builder.trace.push(traceEntry{
		name:      vertexName,
		directive: deps,
	})
	defer builder.trace.pop()

	return slices.MapErr(t.DependsOn, func(dependencyName string) (api.Vertex, error) {
		exists := builder.vertexesSet.Has(dependencyName)
		if !exists {
			return api.Vertex{}, errors.Errorf("%s depends on unknown target %s", vertexName, dependencyName)
		}

		return builder.recursiveGraph(dependencyName)
	})
}

func mapStage(
	stageName string,
	s buildconfig.StageData,
	copyDirs []api.Copy,
	secrets []config.Secret,
) (api.Stage, error) {
	mappedSecrets, err := mapSecrets(s.Secrets, secrets)
	if err != nil {
		return api.Stage{}, errors.Wrapf(err, "failed to map secrets in %s stage", stageName)
	}

	return api.Stage{
		From: s.From,
		Platform: maybe.Map(s.Platform, func(p string) string {
			return p
		}),
		WorkDir: s.WorkDir,
		Env:     s.Env,
		Cache:   slices.Map(s.Cache, mapCache),
		Copy:    copyDirs,
		Network: maybe.Map(s.Network, func(n string) api.Network {
			return api.Network{
				Network: n,
			}
		}),
		SSH: maybe.Map(s.SSH, func(s buildconfig.SSH) api.SSH {
			return api.SSH{}
		}),
		Secrets: mappedSecrets,
		Command: s.Command,
		Output: maybe.Map(s.Output, func(o buildconfig.Output) api.Output {
			return api.Output{
				Artifact: o.Artifact,
				Local:    o.Local,
			}
		}),
	}, nil
}

func mapCache(cache buildconfig.Cache) api.Cache {
	return api.Cache{
		ID:   cache.ID,
		Path: cache.Path,
	}
}

func mapCopy(c buildconfig.Copy) api.CopyVar {
	return api.CopyVar{
		Src:  c.Src,
		Dst:  c.Dst,
		From: c.From,
	}
}

func mapSecrets(secrets []buildconfig.Secret, secretSrc []config.Secret) ([]api.Secret, error) {
	return slices.MapErr(secrets, func(s buildconfig.Secret) (api.Secret, error) {
		return mapSecret(s, secretSrc)
	})
}

func mapSecret(secret buildconfig.Secret, secrets []config.Secret) (api.Secret, error) {
	found := stdslices.ContainsFunc(secrets, func(s config.Secret) bool {
		return s.ID == secret.ID
	})
	if !found {
		return api.Secret{}, errors.Errorf("reference to unknown secret %s", secret.ID)
	}
	return api.Secret{
		ID:        secret.ID,
		MountPath: secret.Path,
	}, nil
}
