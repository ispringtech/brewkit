package dockerfile

import (
	"fmt"

	"github.com/ispringtech/brewkit/internal/backend/api"
	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
	"github.com/ispringtech/brewkit/internal/dockerfile"
)

type VarGenerator interface {
	GenerateDockerfile(vars []api.Var) (dockerfile.Dockerfile, error)
}

func NewVarGenerator(dockerfileImage string) VarGenerator {
	return &varGenerator{dockerfileImage: dockerfileImage}
}

type varGenerator struct {
	dockerfileImage string
}

func (generator varGenerator) GenerateDockerfile(vars []api.Var) (dockerfile.Dockerfile, error) {
	return dockerfile.Dockerfile{
		SyntaxHeader: dockerfile.Syntax(generator.dockerfileImage),
		Stages:       slices.Map(vars, generator.stageForVar),
	}, nil
}

func (generator varGenerator) stageForVar(v api.Var) dockerfile.Stage {
	return dockerfile.Stage{
		From:         v.From,
		As:           maybe.NewJust(v.Name),
		Instructions: generator.instructionsForVar(v),
	}
}

func (generator varGenerator) instructionsForVar(v api.Var) []dockerfile.Instruction {
	//nolint:prealloc
	var instructions []dockerfile.Instruction

	instructions = append(instructions, dockerfile.Workdir(v.WorkDir))

	for k, v := range v.Env {
		instructions = append(instructions, dockerfile.Env{
			K: k,
			V: v,
		})
	}

	for _, c := range v.Copy {
		instructions = append(instructions, dockerfile.Copy{
			Src:  c.Src,
			Dst:  c.Dst,
			From: c.From,
		})
	}

	//nolint:prealloc
	var mounts []dockerfile.Mount

	for _, cache := range v.Cache {
		mounts = append(mounts, dockerfile.MountCache{
			ID:     maybe.NewJust(cache.ID),
			Target: cache.Path,
		})
	}

	for _, secret := range v.Secrets {
		mounts = append(mounts, dockerfile.MountSecret{
			ID:       maybe.NewJust(secret.ID),
			Target:   maybe.NewJust(secret.MountPath),
			Required: maybe.NewJust(true), // make error if secret unavailable
		})
	}

	if maybe.Valid(v.SSH) {
		mounts = append(mounts, dockerfile.MountSSH{
			Required: maybe.NewJust(true), // make error if ssh key unavailable
		})
	}

	var network string
	if maybe.Valid(v.Network) {
		network = maybe.Just(v.Network).Network
	}

	command := generator.transformToHeredoc(v.Command)

	instructions = append(instructions, dockerfile.Run{
		Mounts:  mounts,
		Network: network,
		Command: command,
	})

	return instructions
}

func (generator varGenerator) transformToHeredoc(s string) string {
	const heredocHeader = "EOF"

	return fmt.Sprintf("<<%s\n%s\n%s", heredocHeader, s, heredocHeader)
}
