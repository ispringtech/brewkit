package build

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/backend/app/docker"
	"github.com/ispringtech/brewkit/internal/backend/app/dockerfile"
	"github.com/ispringtech/brewkit/internal/backend/app/reporter"
	"github.com/ispringtech/brewkit/internal/backend/app/ssh"
	"github.com/ispringtech/brewkit/internal/common/maps"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
	df "github.com/ispringtech/brewkit/internal/dockerfile"
)

type Service interface {
	api.BuilderAPI
}

func NewBuildService(
	dockerClient docker.Client,
	dockerfileImage string,
	sshAgentProvider ssh.AgentProvider,
	backendReporter reporter.Reporter,
) Service {
	return &buildService{
		dockerClient:     dockerClient,
		dockerfileImage:  dockerfileImage,
		sshAgentProvider: sshAgentProvider,
		reporter:         backendReporter,
	}
}

type buildService struct {
	dockerClient     docker.Client
	dockerfileImage  string
	sshAgentProvider ssh.AgentProvider
	reporter         reporter.Reporter
}

func (service *buildService) Build(
	ctx context.Context,
	v api.Vertex,
	vars []api.Var,
	secretsSrc []api.SecretSrc,
	params api.BuildParams,
) error {
	err := service.prePullImages(ctx, v, vars, params.ForcePull)
	if err != nil {
		return err
	}

	varsMap, err := service.calculateVars(ctx, vars)
	if err != nil {
		return err
	}

	return service.buildVertex(ctx, v, varsMap, secretsSrc)
}

func (service *buildService) calculateVars(ctx context.Context, vars []api.Var) (dockerfile.Vars, error) {
	if len(vars) == 0 {
		return nil, nil
	}

	d, err := dockerfile.NewVarGenerator(service.dockerfileImage).GenerateDockerfile(vars)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate dockerfile for variables")
	}

	res := map[string]string{}

	for _, v := range vars {
		// Check if context closed before running Value
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		data, err2 := service.dockerClient.Value(ctx, d, docker.ValueParams{
			Var:      v.Name,
			SSHAgent: maybe.NewJust(service.sshAgentProvider.Default()),
			UseCache: false, // Disable cache for retrieving variable value
		})
		if err2 != nil {
			return nil, errors.Wrapf(err2, "failed to calculate %s var", v.Name)
		}

		res[v.Name] = string(data)
	}

	return res, nil
}

func (service *buildService) buildVertex(
	ctx context.Context,
	v api.Vertex,
	vars dockerfile.Vars,
	secretsSrc []api.SecretSrc,
) error {
	d, err := dockerfile.NewTargetGenerator(v, vars, service.dockerfileImage).GenerateDockerfile()
	if err != nil {
		return err
	}

	service.reporter.Debugf("dockerfile:\n%s\n", d.Format())

	executedVertexes := maps.Set[string]{}

	secrets := slices.Map(secretsSrc, func(s api.SecretSrc) docker.SecretData {
		return docker.SecretData{
			ID:   s.ID,
			Path: s.SourcePath,
		}
	})

	var recursiveBuild func(ctx context.Context, v api.Vertex) error
	recursiveBuild = func(ctx context.Context, v api.Vertex) error {
		if executedVertexes.Has(v.Name) {
			// Skip already executed stages
			return nil
		}

		if maybe.Valid(v.From) && shouldExplicitRunFrom(*maybe.Just(v.From)) {
			err2 := recursiveBuild(ctx, *maybe.Just(v.From))
			if err2 != nil {
				return err2
			}
		}

		for _, childVertex := range v.DependsOn {
			err2 := recursiveBuild(ctx, childVertex)
			if err2 != nil {
				return err2
			}
		}

		if !maybe.Valid(v.Stage) {
			return nil
		}

		executedVertexes.Add(v.Name)

		targetName := v.Name
		var output maybe.Maybe[string]

		stage := maybe.Just(v.Stage)
		if maybe.Valid(stage.Output) {
			o := maybe.Just(stage.Output)

			// Execute output stage to save artifacts
			targetName = fmt.Sprintf("%s-out", v.Name)
			output = maybe.NewJust(o.Local)
		}

		// Check if context closed before running build
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		return service.dockerClient.Build(ctx, d, docker.BuildParams{
			Target:   targetName,
			SSHAgent: maybe.NewJust(service.sshAgentProvider.Default()),
			Output:   output,
			Secrets:  secrets,
		})
	}

	return recursiveBuild(ctx, v)
}

func (service *buildService) prePullImages(
	ctx context.Context,
	v api.Vertex,
	vars []api.Var,
	forcePull bool,
) error {
	images := maps.Set[string]{}
	images.Add(service.dockerfileImage)

	images = service.listVertexImages(v, images)
	images = service.listVarsImages(vars, images)

	if forcePull {
		service.reporter.Logf("Force pull images\n")
		for image := range images {
			err2 := service.dockerClient.PullImage(ctx, image)
			if err2 != nil {
				return err2
			}
		}
		return nil
	}

	imagesSlice := maps.ToSlice(images, func(image string, _ struct{}) string {
		return image
	})
	existingImages, err := service.dockerClient.ListImages(ctx, imagesSlice)
	if err != nil {
		return errors.Wrap(err, "failed to filter existing images")
	}

	imagesToPull := slices.Diff(imagesSlice, slices.Map(existingImages, func(img docker.Image) string {
		return img.Img
	}))

	if len(imagesToPull) == 0 {
		return nil
	}

	service.reporter.Logf("Absent images: %s\n", strings.Join(imagesToPull, " "))
	for _, image := range imagesToPull {
		err2 := service.dockerClient.PullImage(ctx, image)
		if err2 != nil {
			return err2
		}
	}

	return nil
}

func (service *buildService) listVertexImages(v api.Vertex, images maps.Set[string]) maps.Set[string] {
	// Recursive walk to From stage
	if maybe.Valid(v.From) {
		images = service.listVertexImages(*maybe.Just(v.From), images)
	}

	// Pull 'From' image
	if !maybe.Valid(v.From) && maybe.Valid(v.Stage) {
		image := maybe.Just(v.Stage).From
		// There is no need to pull scratch image
		if image != df.Scratch && !images.Has(image) {
			images.Add(image)
		}
	}

	if maybe.Valid(v.Stage) {
		copyDirs := maybe.Just(v.Stage).Copy

		for _, c := range copyDirs {
			if !maybe.Valid(c.From) {
				// Skip vertexes with empty from
				continue
			}

			maybe.Just(c.From).
				MapLeft(func(copyV *api.Vertex) {
					images = service.listVertexImages(*copyV, images)
				}).
				MapRight(func(image string) {
					if image != df.Scratch && !images.Has(image) {
						images.Add(image)
					}
				})
		}
	}

	for _, childVertex := range v.DependsOn {
		images = service.listVertexImages(childVertex, images)
	}

	return images
}

func (service *buildService) listVarsImages(vars []api.Var, images maps.Set[string]) maps.Set[string] {
	for _, v := range vars {
		image := v.From
		if !images.Has(image) {
			images.Add(image)
		}
	}

	return images
}

func shouldExplicitRunFrom(v api.Vertex) bool {
	var hasOutput bool
	if maybe.Valid(v.Stage) {
		hasOutput = maybe.Valid(maybe.Just(v.Stage).Output)
	}

	hasDependsOn := len(v.DependsOn) > 0

	return hasOutput || hasDependsOn
}
