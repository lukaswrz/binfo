package binfo

import (
	_ "embed"
	"strings"
	"text/template"
)

type SummaryMode uint

const (
	Module SummaryMode = 1 << iota
	Build
	CGO
	VCS
	Multiline
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

var (
	//go:embed summary.tmpl
	st   string
	t, _ = template.New("").Parse(st)
)

func (b Binfo) Summarize(name string, version string, mode SummaryMode) string {
	wants := func(test SummaryMode) bool {
		return mode&test == test
	}

	var (
		brk string
		sep string
	)

	if wants(Multiline) {
		brk = "\n"
		sep = "\n"
	} else {
		brk = " "
		sep = ", "
	}

	sb := new(strings.Builder)
	err := t.Execute(sb, params{
		Module: wants(Module),
		Build:  wants(Build),
		CGO:    wants(CGO),
		VCS:    wants(VCS),
		Brk:    brk,
		Sep:    sep,
		I:      b,
	})
	if err != nil {
		return ""
	}

	return sb.String()
}
