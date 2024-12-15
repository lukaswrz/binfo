package binfo

import (
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

type SummaryMode uint

const (
	ModeModule SummaryMode = 1 << iota
	ModeBuild
	ModeCGO
	ModeVCS
	ModeMultiline
)

type params struct {
	Name    string
	Version string

	Module bool
	Build  bool
	CGO    bool
	VCS    bool

	Brk string
	Sep string

	I Binfo
}

//go:embed summary.tmpl
var st string

func (b Binfo) Summarize(name string, version string, mode SummaryMode) (string, error) {
	wants := func(test SummaryMode) bool {
		return mode&test == test
	}

	var (
		brk string
		sep string
	)

	if wants(ModeMultiline) {
		brk = "\n"
		sep = "\n"
	} else {
		brk = " "
		sep = ", "
	}

	t, err := template.New("").Parse(st)
	if err != nil {
		return "", fmt.Errorf("cannot parse summary template: %w", err)
	}
	sb := new(strings.Builder)
	err = t.Execute(sb, params{
		Module: wants(ModeModule),
		Build:  wants(ModeBuild),
		CGO:    wants(ModeCGO),
		VCS:    wants(ModeVCS),
		Brk:    brk,
		Sep:    sep,
		I:      b,
	})
	if err != nil {
		return "", fmt.Errorf("cannot execute summary template: %w", err)
	}

	return sb.String(), nil
}

func (b Binfo) MustSummarize(name string, version string, mode SummaryMode) string {
	s, err := b.Summarize(name, version, mode)
	if err != nil {
		panic(err)
	}
	return s
}
