package builddefinition

import (
	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
	"github.com/ispringtech/brewkit/internal/frontend/app/buildconfig"
	"github.com/ispringtech/brewkit/internal/frontend/app/config"
	"github.com/ispringtech/brewkit/internal/frontend/app/version"
)

type Builder interface {
	Build(config buildconfig.Config, secrets []config.Secret) (Definition, error)
}

func NewBuilder() Builder {
	return &builder{}
}

type builder struct{}

func (builder builder) Build(c buildconfig.Config, secrets []config.Secret) (Definition, error) {
	if c.APIVersion != version.APIVersionV1 {
		return Definition{}, errors.Wrapf(ErrUnsupportedAPIVersion, "version: %s", c.APIVersion)
	}

	vertexes, err := newVertexGraphBuilder(secrets, c.Targets).graphVertexes()
	if err != nil {
		return Definition{}, err
	}

	vars, err := builder.variables(c.Vars, secrets)
	if err != nil {
		return Definition{}, err
	}

	return Definition{
		Vertexes: vertexes,
		Vars:     vars,
	}, err
}

func (builder builder) variables(vars []buildconfig.VarData, secrets []config.Secret) ([]api.Var, error) {
	return slices.MapErr(vars, func(v buildconfig.VarData) (api.Var, error) {
		mappedSecrets, err := mapSecrets(v.Secrets, secrets)
		if err != nil {
			return api.Var{}, errors.Wrapf(err, "failed to map secrets in %s variable", v.Name)
		}

		return api.Var{
			Name: v.Name,
			From: v.From,
			Platform: maybe.Map(v.Platform, func(p string) string {
				return p
			}),
			WorkDir: v.WorkDir,
			Env:     v.Env,
			Cache:   slices.Map(v.Cache, mapCache),
			Copy:    slices.Map(v.Copy, mapCopy),
			Network: maybe.Map(v.Network, func(n string) api.Network {
				return api.Network{
					Network: n,
				}
			}),
			SSH: maybe.Map(v.SSH, func(s buildconfig.SSH) api.SSH {
				return api.SSH{}
			}),
			Secrets: mappedSecrets,
			Command: v.Command,
		}, nil
	})
}
