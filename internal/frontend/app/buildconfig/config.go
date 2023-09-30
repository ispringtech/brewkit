package buildconfig

import (
	"github.com/ispringtech/brewkit/internal/common/maybe"
)

type Config struct {
	APIVersion string
	Vars       []VarData
	Targets    []TargetData
}

type VarData struct {
	Name     string
	From     string
	Platform maybe.Maybe[string]
	WorkDir  string
	Env      map[string]string
	Cache    []Cache
	Copy     []Copy
	Secrets  []Secret
	Network  maybe.Maybe[string]
	SSH      maybe.Maybe[SSH]
	Command  string
}

type TargetData struct {
	Name      string
	DependsOn []string
	Stage     maybe.Maybe[StageData]
}

type StageData struct {
	From     string
	Env      map[string]string
	Command  maybe.Maybe[string]
	SSH      maybe.Maybe[SSH]
	Cache    []Cache
	Copy     []Copy
	Secrets  []Secret
	Platform maybe.Maybe[string]
	WorkDir  string
	Network  maybe.Maybe[string]
	Output   maybe.Maybe[Output]
}

type SSH struct{}

type Cache struct {
	ID   string
	Path string
}

type Copy struct {
	From maybe.Maybe[string]
	Src  string
	Dst  string
}

type Secret struct {
	ID   string
	Path string
}

type Output struct {
	Artifact string
	Local    string
}
