package binfo

import (
	"runtime/debug"
	"time"
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
