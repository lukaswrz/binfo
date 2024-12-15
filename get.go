package binfo

import (
	"errors"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/hashicorp/go-multierror"
)

func Get() (Binfo, error) {
	o, ok := debug.ReadBuildInfo()
	if !ok {
		return Binfo{}, errors.New("unable to read build info")
	}

	var merr *multierror.Error

	b := Binfo{}

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

	return b, merr.ErrorOrNil()
}

func MustGet() Binfo {
	b, err := Get()
	if err != nil {
		panic(err)
	}
	return b
}
