package pom

import (
	"regexp"

	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
)

func (p *Project) ResolveVersionProp(dep *gopom.Dependency, projects []*Project) error {
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
	if len(matches[0]) == 0 {
		return nil
	}

	prop := matches[0][1]
	value, err := p.resolvePropAcrossProjects(projects, prop)
	if err != nil {
		return err
	}
	if value == nil {
		return errors.New("dependency version not found")
	}
	dep.Version = *value

	return nil
}

func (p *Project) resolvePropAcrossProjects(projects []*Project, prop string) (*string, error) {
	var value *string

	// Search the property across all POMs.
	for _, project := range projects {
		project := project
		p.logger.Debug().Msg("resolve version")

		// Fallback to profile properties.
		if project.Properties == nil {
			var err error
			value, err = p.resolvePropertyFromProfiles(projects, prop)
			if err != nil || value == nil {
				p.logger.Debug().Err(err).Str("project", project.Name).Str("property", prop).Msg("resolve version from profiles")
				continue
			}
			break
		}
		if project.Properties.Entries == nil {
			var err error
			value, err = p.resolvePropertyFromProfiles(projects, prop)
			if err != nil || value == nil {
				p.logger.Debug().Err(err).Str("project", project.Name).Str("property", prop).Msg("resolve version from profiles")
				continue
			}
			break
		}

		var err error
		value, err = resolvePropertyFromProperties(project.Properties.Entries, prop)
		if err != nil || value == nil {
			p.logger.Debug().Err(err).Str("project", project.Name).Str("property", prop).Msg("resolve version")
			continue
		}
		p.logger.Info().Str("project", project.Name).Str("property", prop).Msg("resolved version from property")
		// TODO: pick the parent POM instead of the first match.
		break
	}

	return value, nil
}

func (p *Project) resolvePropertyFromProfiles(projects []*Project, prop string) (*string, error) {
	var value *string
nextPom:
	for _, project := range projects {
		project := project
		if project.Profiles == nil {
			continue
		}
		for _, profile := range *project.Profiles {
			if profile.Properties == nil {
				continue
			}
			if profile.Properties.Entries == nil {
				continue
			}

			var err error
			value, err = resolvePropertyFromProperties(profile.Properties.Entries, prop)
			if err != nil {
				p.logger.Debug().Err(err).Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolve version from profiles")
				continue
			}
			if value != nil {
				p.logger.Info().Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolved version from profile property")
				// TODO: pick the default profile instead of the first match.
				break nextPom
			}
		}
	}
	if value == nil {
		return nil, errors.New("property not found in profiles")
	}

	return value, nil
}

func resolvePropertyFromProperties(entries map[string]string, property string) (*string, error) {
	if entries == nil {
		return nil, errors.New("entries is empty")
	}

	value, ok := entries[property]
	if !ok {
		return nil, errors.New("property not found")
	}
	if value == "" {
		return nil, errors.New("value is empty")
	}

	return &value, nil
}
