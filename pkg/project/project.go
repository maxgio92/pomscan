package project

import (
	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
	log "github.com/rs/zerolog"
)

// Project is a Maven Project, following the POM model.
// TODO: do not export. Leave the ProjectList as the main API.
type Project struct {
	*gopom.Project
	pomPath string
	logger  *log.Logger
}

// TODO: provide a ProjectList.
// This is needed to traverse a project tree to resolve properties and other.

func NewProject(opts ...ProjOption) *Project {
	project := new(Project)
	for _, f := range opts {
		f(project)
	}

	return project
}

type ProjOption func(*Project)

func WithPomPath(path string) ProjOption {
	return func(p *Project) {
		p.pomPath = path
	}
}

func WithLogger(logger *log.Logger) ProjOption {
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
