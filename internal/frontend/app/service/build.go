package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
	"github.com/ispringtech/brewkit/internal/frontend/app/buildconfig"
	"github.com/ispringtech/brewkit/internal/frontend/app/builddefinition"
	appconfig "github.com/ispringtech/brewkit/internal/frontend/app/config"
)

const (
	allTargetKeyword = "all"
)

type BuildService interface {
	Build(ctx context.Context, p BuildParams) error

	DumpBuildDefinition(ctx context.Context, configPath string) (string, error)
	DumpCompiledBuildDefinition(ctx context.Context, configPath string) (string, error)
}

type BuildParams struct {
	Targets         []string // Target names to run
	BuildDefinition string

	ForcePull bool
}

func NewBuildService(
	configParser buildconfig.Parser,
	definitionBuilder builddefinition.Builder,
	builder api.BuilderAPI,
	config appconfig.Config,
) BuildService {
	return &buildService{
		configParser:      configParser,
		definitionBuilder: definitionBuilder,
		builder:           builder,
		config:            config,
	}
}

type buildService struct {
	configParser      buildconfig.Parser
	definitionBuilder builddefinition.Builder
	builder           api.BuilderAPI
	config            appconfig.Config
}

func (service *buildService) Build(ctx context.Context, p BuildParams) error {
	c, err := service.configParser.Parse(p.BuildDefinition)
	if err != nil {
		return err
	}

	definition, err := service.definitionBuilder.Build(c, service.config.Secrets)
	if err != nil {
		return err
	}

	vertex, err := service.buildVertex(p.Targets, definition)
	if err != nil {
		return err
	}

	secrets := slices.Map(service.config.Secrets, func(s appconfig.Secret) api.SecretSrc {
		return api.SecretSrc{
			ID:         s.ID,
			SourcePath: s.Path,
		}
	})

	return service.builder.Build(
		ctx,
		vertex,
		definition.Vars,
		secrets,
		api.BuildParams{
			ForcePull: p.ForcePull,
		},
	)
}

func (service *buildService) DumpBuildDefinition(_ context.Context, configPath string) (string, error) {
	c, err := service.configParser.Parse(configPath)
	if err != nil {
		return "", err
	}

	definition, err := service.definitionBuilder.Build(c, service.config.Secrets)
	if err != nil {
		return "", err
	}

	d, err := json.Marshal(definition)
	if err != nil {
		return "", errors.WithStack(err)
	}

	return string(d), nil
}

func (service *buildService) DumpCompiledBuildDefinition(_ context.Context, configPath string) (string, error) {
	return service.configParser.CompileConfig(configPath)
}

func (service *buildService) buildVertex(targets []string, definition builddefinition.Definition) (api.Vertex, error) {
	if len(targets) == 0 {
		v, err := service.findTarget(allTargetKeyword, definition)
		return v, errors.Wrap(err, "failed to find default target")
	}

	if len(targets) == 1 {
		return service.findTarget(targets[0], definition)
	}

	vertexes, err := slices.MapErr(targets, func(t string) (api.Vertex, error) {
		return service.findTarget(t, definition)
	})
	if err != nil {
		return api.Vertex{}, err
	}

	return api.Vertex{
		Name:      allTargetKeyword,
		DependsOn: vertexes,
	}, nil
}

func (service *buildService) findTarget(target string, definition builddefinition.Definition) (api.Vertex, error) {
	vertex := definition.Vertex(target)
	if !maybe.Valid(vertex) {
		return api.Vertex{}, errors.Errorf("target %s not found", target)
	}
	return maybe.Just(vertex), nil
}
