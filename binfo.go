package binfo

import (
	"fmt"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
)

type Binfo struct {
	// Module information.
	Module struct {
		// The module path, e.g. "golang.org/x/tools/cmd/stringer".
		Path string

		// The version of this module.
		Version string

		// The checksum of this module.
		Sum string
	}

	// Information about the build process and target.
	Build struct {
		// The version of the Go toolchain used to build the binary, e.g. "go1.19.2".
		GoVersion string

		// The buildmode flag value, e.g. "exe".
		Mode string

		// The compiler toolchain, e.g. "gc" or "gccgo".
		Compiler string

		// The architecture target, e.g. "amd64".
		Arch string

		// The operating system target, e.g. "linux".
		OS string
	}

	// Information about CGO.
	CGO struct {
		Enabled bool

		Flags struct {
			C   string
			CPP string
			CXX string
			LD  string
		}
	}

	// Information about version control at the time of building.
	VCS struct {
		// The name of the version control system, e.g. "git".
		Name string

		// The current revision, e.g. "6cb6d5fa113f26aa2bc139539eab8939632f0693".
		Revision string

		// The modification time for the current revision.
		Time time.Time

		// Whether or not the source tree had local modifications.
		Modified bool
	}

	// The original data source for build information.
	Orig *debug.BuildInfo
}

func Get() (Binfo, error) {
	var merr *multierror.Error

	b := Binfo{}

	if o, ok := debug.ReadBuildInfo(); ok {
		b.Orig = o

		b.Module.Version = o.Main.Version
		b.Module.Path = o.Main.Path
		b.Module.Sum = o.Main.Sum

		for _, setting := range o.Settings {
			switch setting.Key {
			case "-buildmode":
				b.Build.Mode = setting.Value
			case "-compiler":
				b.Build.Compiler = setting.Value
			case "GOARCH":
				b.Build.Arch = setting.Value
			case "GOOS":
				b.Build.OS = setting.Value

			case "CGO_ENABLED":
				switch setting.Value {
				case "1":
					b.CGO.Enabled = true
				case "0":
					b.CGO.Enabled = false
				default:
					merr = multierror.Append(merr, fmt.Errorf("failed to parse %s", setting.Key))
				}
			case "CGO_CFLAGS":
				b.CGO.Flags.C = setting.Value
			case "CGO_CPPFLAGS":
				b.CGO.Flags.CPP = setting.Value
			case "CGO_CXXFLAGS":
				b.CGO.Flags.CXX = setting.Value
			case "CGO_LDFLAGS":
				b.CGO.Flags.LD = setting.Value

			case "vcs":
				b.VCS.Name = setting.Value
			case "vcs.revision":
				b.VCS.Revision = setting.Value
			case "vcs.time":
				v, err := time.Parse(time.RFC3339, setting.Value)
				if err != nil {
					merr = multierror.Append(merr, fmt.Errorf("unable to parse VCS time: %w", err))
				}
				b.VCS.Time = v
			case "vcs.modified":
				v, err := strconv.ParseBool(setting.Value)
				if err != nil {
					merr = multierror.Append(merr, fmt.Errorf("unable to parse VCS modified: %w", err))
				}
				b.VCS.Modified = v
			}
		}
	}

	return b, merr.ErrorOrNil()
}

type SummaryMode uint

const (
	ModeModule SummaryMode = 1 << iota
	ModeBuild
	ModeCGO
	ModeVCS
	ModeMultiline
	ModeNamed
)

func (b Binfo) Summarize(name string, mode SummaryMode) string {
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

	lines := make([]string, 4)

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
		lines = append(
			lines,
			fmt.Sprintf("via %s (rev %s) (at %s)", b.VCS.Name, b.VCS.Revision, b.VCS.Time.Format("2006-01-02 15:04:05")),
		)
	}

	j := strings.Join(lines, sep)

	if name == "" {
		return j
	} else {
		return fmt.Sprintf("%s:%s%s", name, brk, j)
	}
}
