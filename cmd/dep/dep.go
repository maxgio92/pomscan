package dep

import (
	"fmt"
	"regexp"

	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/maxgio92/pomscan/internal/files"
	"github.com/maxgio92/pomscan/internal/options"
)

const (
	pomFile = "pom.xml"
)

type Options struct {
	ArtifactID string
	GroupID    string
	*options.CommonOptions
}

func NewDepCmd(opts *options.CommonOptions) *cobra.Command {
	o := &Options{"", "", opts}

	cmd := &cobra.Command{
		Use:   "dep",
		Short: "Get info about a dependency",
		RunE:  o.Run,
	}

	cmd.Flags().StringVarP(&o.ArtifactID, "artifact-id", "a", "", "Artifact ID")
	cmd.Flags().StringVarP(&o.GroupID, "group-id", "g", "", "Group ID")

	return cmd
}

func (o *Options) Run(_ *cobra.Command, _ []string) error {
	if o.Debug {
		logger := o.Logger.Level(log.DebugLevel)
		o.Logger = &logger
	}

	files, err := files.FindFiles(o.ProjectPath, pomFile)
	if err != nil {
		return errors.Wrap(err, "find pom files")
	}

	poms := make([]*gopom.Project, 0)
	for _, pomPath := range files {
		pom, err := gopom.Parse(pomPath)
		if err != nil {
			return errors.Wrap(err, "parsing pom")
		}
		poms = append(poms, pom)
	}

	for _, pom := range poms {
		dep, err := pom.Search(o.GroupID, o.ArtifactID)
		if err != nil {
			fmt.Println(errors.Wrap(err, "search dependency"))
		}

		if err := o.resolveVersionVariable(dep, poms); err != nil {
			o.Logger.Warn().Err(err).Msg("resolve version")
		}

		printDep(dep, "")
	}

	return nil
}

func (o *Options) resolveVersionVariable(dep *gopom.Dependency, poms []*gopom.Project) error {
	// Check if the dependency contains a variable.
	pattern := `\$\{(.+)\}`

	r := regexp.MustCompile(pattern)

	matches := r.FindAllStringSubmatch(dep.Version, -1)
	if len(matches) > 0 && len(matches[0]) > 1 {
		o.Logger.Debug().Msg("resolve version")

		varname := matches[0][1]

		value, ok := poms[0].Properties.Entries[varname]
		if !ok {
			return errors.New("dependency version not found")
		}
		dep.Version = value
	}

	// TODO: search the variable declaration recursively.

	return nil
}

func printDep(dep *gopom.Dependency, pomPath string) {
	fmt.Printf("📦 %s.%s found\n", dep.GroupID, dep.ArtifactID)
	fmt.Println("pom:", pomPath)
	if dep.Version != "" {
		fmt.Println("version:", dep.Version)
	}
	if dep.Scope != "" {
		fmt.Println("scope:", dep.Scope)
	}
	fmt.Println()
}
