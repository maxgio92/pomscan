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

	// Print only dependencies that have a version set in the POMs.
	// It includes cases where the version is set with properties.
	VersionOnly bool
	*options.CommonOptions
}

func NewDepCmd(opts *options.CommonOptions) *cobra.Command {
	o := &Options{"", "", false, opts}

	cmd := &cobra.Command{
		Use:   "dep",
		Short: "Get info about a dependency",
		RunE:  o.Run,
	}

	cmd.Flags().StringVarP(&o.ArtifactID, "artifact-id", "a", "", "Filter by artifact ID.")
	cmd.MarkFlagRequired("artifact-id")
	cmd.Flags().StringVarP(&o.GroupID, "group-id", "g", "", "Filter by group ID. It must be combined with artifact ID.")
	cmd.Flags().BoolVar(&o.VersionOnly, "version-only", false, "Print only matches that have the version set. It supports properties.")

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

	for i, pom := range poms {
		var dep *gopom.Dependency
		var err error
		// Search in direct dependencies.
		if o.GroupID != "" {
			dep, err = pom.Search(o.GroupID, o.ArtifactID)
		} else {
			dep, err = pom.SearchByArtifactID(o.ArtifactID)
		}
		if err != nil {
			o.Logger.Debug().Err(err).Str("pom", files[i]).Msg("search dependency")

			// Search in inherited dependencies.
			o.Logger.Debug().Str("project", pom.Name).Msg("searching between inherited dependencies")
			dep, err = searchDepInDepManagementSection(o.GroupID, o.ArtifactID, pom)
			if err != nil {
				o.Logger.Debug().Err(err).Str("pom", files[i]).Msg("search dependency")
				continue
			}
		}

		if err := o.resolveVersionProperty(dep, poms); err != nil {
			o.Logger.Debug().Err(err).Str("pom", files[i]).Msg("resolve version")
		}

		o.printDep(dep, files[i])
	}

	return nil
}

func searchDepInDepManagementSection(groupID, artifactID string, pom *gopom.Project) (*gopom.Dependency, error) {
	if pom.DependencyManagement == nil {
		return nil, errors.New("pom dependency management is empty")
	}
	if pom.DependencyManagement.Dependencies == nil {
		return nil, errors.New("pom dependency management' dependency list is empty")
	}
	for _, dep := range *pom.DependencyManagement.Dependencies {
		if (groupID == "" && dep.ArtifactID == artifactID) ||
			(dep.ArtifactID == artifactID && dep.GroupID == groupID) {
			{
				return &dep, nil
			}
		}
	}

	return nil, errors.New("dependency not found")
}

func (o *Options) resolveVersionProperty(dep *gopom.Dependency, poms []*gopom.Project) error {
	if dep.Version == "" {
		return nil
	}

	// Check if the dependency contains a property.
	pattern := `\$\{(.+)\}`

	r := regexp.MustCompile(pattern)

	matches := r.FindAllStringSubmatch(dep.Version, -1)

	// The version does not contain a property.
	if len(matches) == 0 {
		return nil
	}

	for _, pom := range poms {
		pom := pom
		o.Logger.Debug().Msg("resolve version")
		if len(matches[0]) == 0 {
			continue
		}

		prop := matches[0][1]

		value, ok := pom.Properties.Entries[prop]
		if !ok {
			continue
		}
		if value != "" {
			dep.Version = value
			o.Logger.Info().Str("project", pom.Name).Str("property", prop).Msg("resolved version")

			// TODO: pick the parent POM instead of the first match.
			break
		}
	}
	if dep.Version == "" {
		return errors.New("dependency version not found")
	}

	return nil
}

func (o *Options) printDep(dep *gopom.Dependency, pomPath string) {
	if o.VersionOnly && dep.Version == "" {
		return
	}

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
