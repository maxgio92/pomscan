package project

import (
	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
	"os"
)

var (
	ErrNoProject   = errors.New("no project found")
	ErrDepNotFound = errors.New("dependency not found")
)

type ProjectList struct {
	projects []*Project
	pomPaths []string
	logger   *log.Logger
}

func NewProjectList(opts ...ProjListOption) *ProjectList {
	list := new(ProjectList)

	list.projects = make([]*Project, 0)
	list.pomPaths = make([]string, 0)
	logger := log.New(os.Stderr).Level(log.InfoLevel)
	list.logger = &logger

	for _, f := range opts {
		f(list)
	}

	for _, path := range list.pomPaths {
		list.projects = append(list.projects, NewProject(
			WithPomPath(path),
			WithLogger(list.logger),
		))
	}

	return list
}

type ProjListOption func(*ProjectList)

func ListWithPomPaths(paths ...string) ProjListOption {
	return func(plist *ProjectList) {
		for _, path := range paths {
			plist.pomPaths = append(plist.pomPaths, path)
		}
	}
}

func ListWithLogger(logger *log.Logger) ProjListOption {
	return func(p *ProjectList) {
		p.logger = logger
	}
}

func (list *ProjectList) LoadAll() error {
	if list.projects == nil || len(list.projects) == 0 {
		return ErrNoProject
	}

	var err error
	for i, _ := range list.projects {
		err = list.projects[i].Load()
		if err != nil {
			return errors.Wrapf(err, "loading project from %s", list.projects[i].pomPath)
		}
	}

	return nil
}

type Dependency struct {
	*gopom.Dependency
	Metadata *DependencyMetadata
}

type DependencyMetadata struct {
	PomPath         string
	VersionProperty *Property
}

func newDependency(dep *gopom.Dependency) *Dependency {
	return &Dependency{
		Dependency: dep,
		Metadata:   &DependencyMetadata{},
	}
}

func (list *ProjectList) SearchDirectDependency(artifactID, groupID string) ([]*Dependency, error) {
	results := make([]*Dependency, 0)
	for _, project := range list.projects {
		project := project
		var dep *gopom.Dependency
		var err error

		// Search in the current project's direct dependencies.
		if groupID != "" {
			dep, err = project.Search(groupID, artifactID)
		} else {
			dep, err = project.SearchByArtifactID(artifactID)
		}
		if err != nil {
			list.logger.Debug().Err(err).Str("pom", project.pomPath).Msg("search dependency")

			// Search in the current project's inherited dependencies.
			list.logger.Debug().Str("project", project.Name).Msg("searching between inherited dependencies")
			dep, err = project.SearchDepInDepMgmtSec(groupID, artifactID)
			if err != nil {
				list.logger.Debug().Err(err).Str("pom", project.pomPath).Msg("search dependency")
				continue
			}
		}

		prop, err := list.ResolveDepVersionProp(dep)
		if err != nil {
			list.logger.Debug().Err(err).Str("pom", project.pomPath).Msg("resolve version")
		}

		// TODO: return matched pom file path.
		res := newDependency(dep)
		res.Metadata.PomPath = project.pomPath
		res.Metadata.VersionProperty = prop
		results = append(results, res)
	}
	if len(results) == 0 {
		return nil, ErrDepNotFound
	}

	return results, nil
}
