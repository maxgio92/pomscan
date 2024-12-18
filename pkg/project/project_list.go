package project

import (
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
