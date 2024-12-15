package binfo

import (
	"fmt"
	"strings"
)

type SummaryMode uint

const (
	ModeModule SummaryMode = 1 << iota
	ModeBuild
	ModeCGO
	ModeVCS
	ModeMultiline
)

func (b Binfo) Summarize(name string, version string, mode SummaryMode) string {
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

	lines := make([]string, 0, 4)

	if wants(ModeModule) {
		lines = append(
			lines,
			fmt.Sprintf("module %s (%s) (sum %s)", b.Module.Path, b.Module.Version, b.Module.Sum),
		)
	}

	if wants(ModeBuild) {
		lines = append(
			lines,
			fmt.Sprintf("built with %s (%s) (mode %s)", b.Build.Compiler, b.Build.GoVersion, b.Build.Mode),
		)
	}

	if wants(ModeCGO) {
		if b.CGO.Enabled {
			lines = append(
				lines,
				fmt.Sprintf("with cgo (c %q) (cpp %q) (cxx %q) (ld %q)", b.CGO.Flags.C, b.CGO.Flags.CPP, b.CGO.Flags.CXX, b.CGO.Flags.LD),
			)
		} else {
			lines = append(
				lines,
				"without cgo",
			)
		}
	}

	if wants(ModeVCS) {
		var m string
		if b.VCS.Modified {
			m = " (modified)"
		} else {
			m = ""
		}

		lines = append(
			lines,
			fmt.Sprintf("via %s (rev %s) (at %s)%s", b.VCS.Name, b.VCS.Revision, b.VCS.Time.Format("2006-01-02 15:04:05"), m),
		)
	}

	j := strings.Join(lines, sep)

	if name == "" {
		return j
	} else {
		return fmt.Sprintf("%s %s:%s%s", name, version, brk, j)
	}
}
