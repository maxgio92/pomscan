package dep

import (
	"github.com/maxgio92/gopom"
	"github.com/maxgio92/pomscan/internal/output"
	"github.com/maxgio92/pomscan/pkg/pom"
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

	pomPaths, err := files.FindFiles(o.ProjectPath, pomFile)
	if err != nil {
		return errors.Wrap(err, "find pom files")
	}

	projects := make([]*pom.Project, 0)
	for _, pomPath := range pomPaths {
		project := pom.NewProject(
			pom.WithPomPath(pomPath),
			pom.WithLogger(o.Logger),
		)
		err = project.Load()
		if err != nil {
			return errors.Wrap(err, "parsing pom")
		}
		projects = append(projects, project)
	}

	for i, project := range projects {
		var dep *gopom.Dependency
		var err error
		// Search in direct dependencies.
		if o.GroupID != "" {
			dep, err = project.Search(o.GroupID, o.ArtifactID)
		} else {
			dep, err = project.SearchByArtifactID(o.ArtifactID)
		}
		if err != nil {
			o.Logger.Debug().Err(err).Str("pom", pomPaths[i]).Msg("search dependency")

			// Search in inherited dependencies.
			o.Logger.Debug().Str("project", project.Name).Msg("searching between inherited dependencies")
			dep, err = project.SearchDepInDepMgmtSec(o.GroupID, o.ArtifactID)
			if err != nil {
				o.Logger.Debug().Err(err).Str("pom", pomPaths[i]).Msg("search dependency")
				continue
			}
		}

		if err := project.ResolveVersionProp(dep, projects); err != nil {
			o.Logger.Debug().Err(err).Str("pom", pomPaths[i]).Msg("resolve version")
		}

		output.PrintDep(dep, pomPaths[i], o.VersionOnly)
	}

	return nil
}
