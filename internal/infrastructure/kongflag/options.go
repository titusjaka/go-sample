package kongflag

import (
	"fmt"
	"os"
	"reflect"
	"runtime/debug"
	"sort"
	"strconv"
	"strings"

	"github.com/alecthomas/kong"
)

// WithDumpEnvs adds an optional flag '--dump-envs' to the service.
// If the application is run with this flag, it prints env-vars with default values and exits.
//
// Usage:
//
//	ctx := kong.Parse(&app kongflag.WithDumpEnvs())
//	ctx.FatalIfErrorf(ctx.Run())
//
// CLI usage:
//
//	$ go run main.go --dump-envs > example.env
//
// nolint:dupl // False-positive here
func WithDumpEnvs() kong.Option {
	return kong.PostBuild(func(k *kong.Kong) error {
		var dumpEnvsTarget dumpEnvs
		value := reflect.ValueOf(&dumpEnvsTarget).Elem()
		dumpEnvsFlag := &kong.Flag{
			Value: &kong.Value{
				Name:         "dump-envs",
				Help:         "Print env variables and their defaults to STDOUT and exit. Use redirection to write into file.",
				OrigHelp:     "Print env variables and their defaults to STDOUT and exit. Use redirection to write into file.",
				Mapper:       kong.MapperFunc(boolFlagMapper),
				Target:       value,
				Tag:          &kong.Tag{},
				DefaultValue: reflect.ValueOf(false),
				Passthrough:  true,
			},
		}

		k.Model.Flags = append(k.Model.Flags, dumpEnvsFlag)
		return nil
	})
}

// WithVersion adds an optional flag '--version' to the service.
// If the application is run with this flag, it prints version and exits.
//
// This option tries to extract version from build info. If it's empty,
// fallbackVersion is used.
//
// Usage:
//
//	ctx := kong.Parse(&app, kongflag.PrintVersionFlag("v1.0.0"))
//	ctx.FatalIfErrorf(ctx.Run())
//
// CLI usage:
//
//	$ go run main.go --version
//	example version: v1.0.0
//	 git branch: main
//	 git revision: 795859ac3599bca75bb417bd0f97303f79289b65
//	 go: go1.19.5
//
// nolint:dupl,forbidigo // False-positive here
func WithVersion(fallbackVersion string) kong.Option {
	return kong.PostBuild(func(k *kong.Kong) error {
		version := fetchVersionFromBuildInfo()
		if version == "" {
			version = fallbackVersion
		}

		if err := addVariableToKong(k, ServiceVersion, version); err != nil {
			return err
		}

		var versionTarget versionFlag
		value := reflect.ValueOf(&versionTarget).Elem()
		flag := &kong.Flag{
			Value: &kong.Value{
				Name:         "version",
				Help:         "Show the version of the service and exit.",
				OrigHelp:     "Show the version of the service and exit.",
				Mapper:       kong.MapperFunc(boolFlagMapper),
				Target:       value,
				Tag:          &kong.Tag{},
				DefaultValue: reflect.ValueOf(false),
				Passthrough:  true,
			},
		}

		k.Model.Flags = append(k.Model.Flags, flag)
		return nil
	})
}

// WithBuildInfo fetches data from debug.BuildInfo and saves it to the kong.Vars
//
// Saved variables:
//   - GoVersion ("go_version") – version of Go
//   - GitCommitSHA ("git_commit_sha") – SHA of the latest commit
func WithBuildInfo() kong.Option {
	return kong.PostBuild(func(k *kong.Kong) error {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return nil
		}

		vars := kong.Vars{}
		vars[GoVersion] = info.GoVersion

		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && setting.Value != "" {
				vars[GitCommitSHA] = setting.Value
			}
		}

		return vars.Apply(k)
	})
}

func boolFlagMapper(ctx *kong.DecodeContext, target reflect.Value) error {
	if ctx.Scan.Peek().Type == kong.FlagValueToken {
		token := ctx.Scan.Pop()
		switch v := token.Value.(type) {
		case string:
			b, err := strconv.ParseBool(v)
			if err != nil {
				return fmt.Errorf("parse bool flag: %w", err)
			}
			target.SetBool(b)
		case bool:
			target.SetBool(v)
		default:
			return fmt.Errorf("expected bool but got %q (%T)", token.Value, token.Value)
		}
	} else {
		target.SetBool(true)
	}
	return nil
}

// dumpEnvs is a hook flag. If called, kong will print envs to STDOUT and exit
type dumpEnvs bool

func (dumpEnvs) Decode(_ *kong.DecodeContext) error { return nil }

func (dumpEnvs) IsBool() bool { return true }

func (d dumpEnvs) BeforeApply(ctx *kong.Context) error {
	flags := ctx.Flags()

	// Sort flags by ENV var
	sort.Slice(flags, func(i, j int) bool {
		if len(flags[i].Envs) == 0 {
			return true
		}

		if len(flags[j].Envs) == 0 {
			return false
		}

		return flags[i].Envs[0] < flags[j].Envs[0]
	})

	groups := d.buildKongGroups(flags)

	flagsDefaults := d.buildKongFlagsDefaults(groups)

	_, _ = os.Stdout.WriteString(flagsDefaults)
	ctx.Kong.Exit(0)

	return nil
}

//nolint:gocognit // This function just prints flags to a string, it's fine
func (dumpEnvs) buildKongFlagsDefaults(groups map[string][]*kong.Flag) string {
	builder := strings.Builder{}

	titles := make([]string, 0, len(groups))

	for title := range groups {
		titles = append(titles, title)
	}

	sort.Strings(titles)

	for _, title := range titles {
		builder.WriteString("# -------------------------------------------------------\n")
		builder.WriteString("# " + title + "\n")
		builder.WriteString("# -------------------------------------------------------\n\n")

		for _, flag := range groups[title] {
			if len(flag.Envs) > 0 {
				if flag.Help != "" {
					builder.WriteString("# " + flag.Help + "\n")
				}

				builder.WriteString(flag.Envs[0] + "=")

				if flag.HasDefault {
					if flag.DefaultValue.Type().Name() == "string" {
						builder.WriteString(`"` + flag.Default + `"`)
					} else {
						builder.WriteString(flag.Default)
					}
				}

				builder.WriteString("\n")
			}
		}

		builder.WriteString("\n")
	}

	return builder.String()
}

func (dumpEnvs) buildKongGroups(flags []*kong.Flag) map[string][]*kong.Flag {
	groups := make(map[string][]*kong.Flag)

	for _, flag := range flags {
		if len(flag.Envs) == 0 {
			continue
		}

		var title string

		if flag.Group != nil {
			title = flag.Group.Title
		} else {
			title = "Default"
		}

		groups[title] = append(groups[title], flag)
	}
	return groups
}

// versionFlag is a hook flag. If called, kong will print version to STDOUT and exit
type versionFlag bool

func (versionFlag) Decode(_ *kong.DecodeContext) error { return nil }

func (versionFlag) IsBool() bool { return true }
func (versionFlag) BeforeApply(ctx *kong.Context) error {
	appName := ctx.Kong.Model.Name

	vars := ctx.Kong.Model.Vars()

	builder := strings.Builder{}

	version := vars[ServiceVersion]
	if version == "" {
		version = "unknown"
	}

	builder.WriteString(fmt.Sprintf("%s version: %s\n", appName, version))

	if branch, ok := vars[GitBranch]; ok {
		builder.WriteString(fmt.Sprintf(" git branch: %s\n", branch))
	}

	if revision, ok := vars[GitCommitSHA]; ok {
		builder.WriteString(fmt.Sprintf(" git revision: %s\n", revision))
	}

	if goVersion, ok := vars[GoVersion]; ok {
		builder.WriteString(fmt.Sprintf(" go: %s\n", goVersion))
	}

	//nolint:forbidigo
	fmt.Print(builder.String())

	ctx.Kong.Exit(0)
	return nil
}

func addVariableToKong(k *kong.Kong, name, value string) error {
	versionVars := kong.Vars{name: value}
	if err := versionVars.Apply(k); err != nil {
		return fmt.Errorf("apply %q var to kong: %w", name, err)
	}
	return nil
}

func fetchVersionFromBuildInfo() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	version := info.Main.Version

	if version == "" || version == "(devel)" {
		return ""
	}

	return version
}
