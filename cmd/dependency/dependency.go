package dependency

import (
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
	"github.com/spf13/cobra"

	"github.com/maxgio92/pomscan/internal/files"
	"github.com/maxgio92/pomscan/internal/options"
	"github.com/maxgio92/pomscan/internal/output"
	"github.com/maxgio92/pomscan/pkg/project"
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

	// Whether to search build plugin dependencies.
	SearchPlugin bool
	*options.CommonOptions
}

func NewDepCmd(opts *options.CommonOptions) *cobra.Command {
	o := &Options{"", "", false, false, opts}

	cmd := &cobra.Command{
		Use:     "dependency",
		Aliases: []string{"dep"},
		Short:   "Search an artifact through the direct runtime dependencies across the project hierarchy.",
		RunE:    o.Run,
	}

	cmd.Flags().StringVarP(&o.ArtifactID, "artifact-id", "a", "", "Filter by artifact ID.")
	cmd.MarkFlagRequired("artifact-id")
	cmd.Flags().StringVarP(&o.GroupID, "group-id", "g", "", "Filter by group ID. It must be combined with artifact ID.")
	cmd.Flags().BoolVar(&o.VersionOnly, "version-only", false, "Print only matches that have the version set. It supports properties.")
	cmd.Flags().BoolVar(&o.SearchPlugin, "plugin", false, "Whether to search an artifact through the plugin build dependencies.")

	return cmd
}

func (o *Options) Run(_ *cobra.Command, _ []string) error {
	if o.Debug {
		logger := o.Logger.Level(log.DebugLevel)
		o.Logger = &logger
	}

	pomPaths, err := files.FindFiles(o.ProjectPath, pomFile)
	if err != nil {
		return errors.Wrap(err, "find project files")
	}

	projectList := project.NewProjectList(
		project.ListWithPomPaths(pomPaths...),
		project.ListWithLogger(o.Logger),
	)
	err = projectList.LoadAll()
	if err != nil {
		return errors.Wrap(err, "loading projects")
	}

	deps := make([]*project.Dependency, 0)
	if !o.SearchPlugin {
		// Runtime dependencies.
		var err error
		deps, err = projectList.SearchDirectDependency(o.ArtifactID, o.GroupID)
		if err != nil {
			return errors.Wrap(err, "searching direct runtime dependency")
		}
	} else {
		// Build time dependencies.
		var err error
		deps, err = projectList.SearchPluginDependency(o.ArtifactID, o.GroupID)
		if err != nil {
			return errors.Wrap(err, "searching plugin build dependency")
		}
	}

	for _, dep := range deps {
		output.PrintDep(dep, o.VersionOnly)
	}

	return nil
}
