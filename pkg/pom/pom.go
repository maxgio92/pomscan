package pom

import (
	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
)

// Project is a Maven Project, following the POM model.
type Project struct {
	*gopom.Project
	pomPath string
	logger  *log.Logger
}

// TODO: provide a ProjectList.
// This is needed to traverse a project tree to resolve properties and other.

func NewProject(opts ...Option) *Project {
	project := new(Project)
	for _, f := range opts {
		f(project)
	}

	return project
}

type Option func(*Project)

func WithPomPath(path string) Option {
	return func(p *Project) {
		p.pomPath = path
	}
}

func WithLogger(logger *log.Logger) Option {
	return func(p *Project) {
		p.logger = logger
	}
}

func (p *Project) Load() error {
	pom, err := gopom.Parse(p.pomPath)
	if err != nil {
		return errors.Wrap(err, "parsing pom")
	}

	p.Project = pom

	return nil
}

func (p *Project) SearchDepInDepMgmtSec(groupID, artifactID string) (*gopom.Dependency, error) {
	if p.DependencyManagement == nil {
		return nil, errors.New("pom dependency management is empty")
	}
	if p.DependencyManagement.Dependencies == nil {
		return nil, errors.New("pom dependency management' dependency list is empty")
	}
	for _, dep := range *p.DependencyManagement.Dependencies {
		if (groupID == "" && dep.ArtifactID == artifactID) ||
			(dep.ArtifactID == artifactID && dep.GroupID == groupID) {
			{
				return &dep, nil
			}
		}
	}

	return nil, errors.New("dependency not found")
}
