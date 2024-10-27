package project

import "github.com/maxgio92/gopom"

// Dependency represents a runtime dependency.
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
		Metadata:   new(DependencyMetadata),
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

		// Resolve dependency version.
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

func (p *Project) SearchDepInDepMgmtSec(groupID, artifactID string) (*gopom.Dependency, error) {
	if p.DependencyManagement == nil {
		return nil, ErrDepMgmtEmpty
	}
	if p.DependencyManagement.Dependencies == nil {
		return nil, ErrDepMgmtEmpty
	}
	for _, dep := range *p.DependencyManagement.Dependencies {
		if (groupID == "" && dep.ArtifactID == artifactID) ||
			(dep.ArtifactID == artifactID && dep.GroupID == groupID) {
			{
				return &dep, nil
			}
		}
	}

	return nil, ErrDepNotFound
}
