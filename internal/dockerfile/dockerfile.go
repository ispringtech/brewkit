package dockerfile

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ispringtech/brewkit/internal/common/maybe"
	"github.com/ispringtech/brewkit/internal/common/slices"
)

type Syntax string

const (
	Scratch = "scratch"

	Dockerfile14 Syntax = "docker/dockerfile:1.4"
)

type Dockerfile struct {
	SyntaxHeader Syntax
	Stages       []Stage
}

func (d Dockerfile) Format() string {
	s := []string{
		fmt.Sprintf("# syntax=%s", d.SyntaxHeader),
		strings.Join(slices.Map(d.Stages, func(s Stage) string {
			return s.Format()
		}), "\n"),
	}

	return strings.Join(s, "\n")
}

type Stage struct {
	From         string
	As           maybe.Maybe[string]
	Instructions []Instruction
}

func (s Stage) Format() string {
	instructions := strings.Join(slices.Map(s.Instructions, func(i Instruction) string {
		return i.FormatInstruction()
	}), "\n")

	var asBlock string
	if maybe.Valid(s.As) {
		asBlock = fmt.Sprintf("as %s", maybe.Just(s.As))
	}

	return fmt.Sprintf("FROM %s %s\n%s", s.From, asBlock, instructions)
}

type Instruction interface {
	FormatInstruction() string
}

type Workdir string

func (w Workdir) FormatInstruction() string {
	return fmt.Sprintf("WORKDIR %s", w)
}

type Env struct {
	K, V string
}

func (e Env) FormatInstruction() string {
	return fmt.Sprintf("ENV %s=%s", e.K, e.V)
}

type Copy struct {
	Src  string
	Dst  string
	From maybe.Maybe[string]
}

func (c Copy) FormatInstruction() string {
	var from string
	if maybe.Valid(c.From) {
		from = fmt.Sprintf("--from=%s", maybe.Just(c.From))
	}

	return fmt.Sprintf("COPY %s %s %s", from, c.Src, c.Dst)
}

type Run struct {
	Mounts  []Mount
	Network string
	Command string
}

func (r Run) FormatInstruction() string {
	instructions := make([]string, 0, len(r.Mounts))
	for _, mount := range r.Mounts {
		m := mount.FormatMount()
		instructions = append(instructions, fmt.Sprintf("--mount=%s", m))
	}

	if r.Network != "" {
		instructions = append(instructions, fmt.Sprintf("--network=%s", r.Network))
	}

	return fmt.Sprintf(
		"RUN %s \\\n %s",
		strings.Join(instructions, " \\\n"),
		r.Command,
	)
}

type Mount interface {
	FormatMount() string
}

type MountBind struct {
	Target    string
	Source    maybe.Maybe[string]
	From      maybe.Maybe[string]
	ReadWrite maybe.Maybe[bool]
}

func (m MountBind) FormatMount() string {
	s := settings{}

	s.addKV("type", "bind")
	s.addKV("target", m.Target)

	if maybe.Valid(m.Source) {
		s.addKV("source", maybe.Just(m.Source))
	}

	if maybe.Valid(m.From) {
		s.addKV("from", maybe.Just(m.From))
	}

	if maybe.Valid(m.ReadWrite) {
		s.addKV("rw", strconv.FormatBool(maybe.Just(m.ReadWrite)))
	}

	return s.formatSettings()
}

type MountCache struct {
	ID       maybe.Maybe[string]
	Target   string
	ReadOnly maybe.Maybe[bool]
	From     maybe.Maybe[string]
	Source   maybe.Maybe[string]
	Mode     maybe.Maybe[string]
	UID      maybe.Maybe[string]
	GID      maybe.Maybe[string]
}

func (m MountCache) FormatMount() string {
	s := settings{}

	s.addKV("type", "cache")
	s.addKV("target", m.Target)

	if maybe.Valid(m.ReadOnly) {
		s.addKV("source", strconv.FormatBool(maybe.Just(m.ReadOnly)))
	}

	if maybe.Valid(m.From) {
		s.addKV("from", maybe.Just(m.From))
	}

	if maybe.Valid(m.Source) {
		s.addKV("source", maybe.Just(m.Source))
	}

	if maybe.Valid(m.Mode) {
		s.addKV("mode", maybe.Just(m.Mode))
	}

	if maybe.Valid(m.UID) {
		s.addKV("uid", maybe.Just(m.UID))
	}

	if maybe.Valid(m.GID) {
		s.addKV("gid", maybe.Just(m.GID))
	}

	return s.formatSettings()
}

type MountSSH struct {
	ID       maybe.Maybe[string]
	Target   maybe.Maybe[string]
	Required maybe.Maybe[bool]
	Mode     maybe.Maybe[string]
	UID      maybe.Maybe[string]
	GID      maybe.Maybe[string]
}

func (m MountSSH) FormatMount() string {
	s := settings{}

	s.addKV("type", "ssh")

	if maybe.Valid(m.Target) {
		s.addKV("target", maybe.Just(m.Target))
	}

	if maybe.Valid(m.Required) {
		s.addKV("required", strconv.FormatBool(maybe.Just(m.Required)))
	}

	if maybe.Valid(m.Mode) {
		s.addKV("mode", maybe.Just(m.Mode))
	}

	if maybe.Valid(m.UID) {
		s.addKV("uid", maybe.Just(m.UID))
	}

	if maybe.Valid(m.GID) {
		s.addKV("gid", maybe.Just(m.GID))
	}

	return s.formatSettings()
}

type MountSecret struct {
	ID       maybe.Maybe[string]
	Target   maybe.Maybe[string]
	Required maybe.Maybe[bool]
	Mode     maybe.Maybe[string]
	UID      maybe.Maybe[string]
	GID      maybe.Maybe[string]
}

func (m MountSecret) FormatMount() string {
	s := settings{}

	s.addKV("type", "secret")

	if maybe.Valid(m.ID) {
		s.addKV("id", maybe.Just(m.ID))
	}

	if maybe.Valid(m.Target) {
		s.addKV("target", maybe.Just(m.Target))
	}

	if maybe.Valid(m.Required) {
		s.addKV("required", strconv.FormatBool(maybe.Just(m.Required)))
	}

	if maybe.Valid(m.Mode) {
		s.addKV("mode", maybe.Just(m.Mode))
	}

	if maybe.Valid(m.UID) {
		s.addKV("uid", maybe.Just(m.UID))
	}

	if maybe.Valid(m.GID) {
		s.addKV("gid", maybe.Just(m.GID))
	}

	return s.formatSettings()
}

type settings []string

func (s *settings) addKV(k, v string) {
	*s = append(*s, fmt.Sprintf("%s=%s", k, v))
}

func (s *settings) formatSettings() string {
	return strings.Join(*s, ",")
}
