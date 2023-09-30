package builddefinition

import (
	"github.com/ispringtech/brewkit/internal/common/either"
	"github.com/ispringtech/brewkit/internal/common/maybe"
)

type Config struct {
	APIVersion string                                     `json:"apiVersion"`
	Targets    map[string]either.Either[[]string, Target] `json:"targets"`
	Vars       map[string]Var                             `json:"vars"`
}

type Target struct {
	DependsOn []string `json:"dependsOn"`
	*Stage    `json:",inline"`
}

type Stage struct {
	From     string                          `json:"from"`
	Env      map[string]string               `json:"env"`
	SSH      maybe.Maybe[SSH]                `json:"ssh"`
	Cache    []Cache                         `json:"cache"`
	Copy     either.Either[[]Copy, Copy]     `json:"copy"`
	Secrets  either.Either[[]Secret, Secret] `json:"secret"`
	Platform maybe.Maybe[string]             `json:"platform"`
	WorkDir  string                          `json:"workdir"`
	Network  maybe.Maybe[string]             `json:"network"`
	Command  maybe.Maybe[string]             `json:"command"`
	Output   maybe.Maybe[Output]             `json:"output"`
}

type Var struct {
	From     string                          `json:"from"`
	Platform maybe.Maybe[string]             `json:"platform"`
	WorkDir  string                          `json:"workdir"`
	Env      map[string]string               `json:"env"`
	Cache    []Cache                         `json:"cache"`
	Copy     either.Either[[]Copy, Copy]     `json:"copy"`
	Secrets  either.Either[[]Secret, Secret] `json:"secrets"`
	Network  maybe.Maybe[string]             `json:"network"`
	SSH      maybe.Maybe[SSH]                `json:"ssh"`
	Command  string                          `json:"command"`
}

type Cache struct {
	ID   string `yaml:"id"`
	Path string `yaml:"path"`
}

type Copy struct {
	From maybe.Maybe[string] `json:"from"`
	Src  string              `json:"src"`
	Dst  string              `json:"dst"`
}

type SSH struct{}

type Secret struct {
	ID   string `yaml:"id"`
	Path string `yaml:"path"`
}

type Output struct {
	Artifact string `json:"artifact"`
	Local    string `json:"local"`
}
