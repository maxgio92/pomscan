package project

import (
	"github.com/maxgio92/gopom"
	"github.com/pkg/errors"
	"regexp"
)

var (
	ErrPropNotFoundInProfiles = errors.New("property not found in profiles")
	ErrPropNotFound           = errors.New("property not found")
	ErrPropValueEmpty         = errors.New("property value is empty")
	ErrPropEntriesEmpty       = errors.New("properties entry is empty")
	ErrDepVersionNotFound     = errors.New("dependency version not found")
	ErrProfilesNotFound       = errors.New("profiles are not found")
)

type Property struct {
	Name     string
	Value    string
	Metadata *PropertyMetadata
}

type PropertyMetadata struct {
	// DeclarePath is the path to the pom where the property is declared
	// and inherited by child projects.
	DeclarePath string
}

func NewProperty(name string) *Property {
	return &Property{
		Name:     name,
		Metadata: &PropertyMetadata{},
	}
}

func (list *ProjectList) ResolveDepVersionProp(dep *gopom.Dependency) (*Property, error) {
	property := new(Property)
	if dep.Version == "" {
		return property, nil
	}

	// Check if the dependency contains a property.
	pattern := `\$\{(.+)\}`
	reg := regexp.MustCompile(pattern)
	matches := reg.FindAllStringSubmatch(dep.Version, -1)

	// The version does not contain a property.
	if len(matches) == 0 {
		return property, nil
	}
	if len(matches[0]) == 0 {
		return property, nil
	}

	name := matches[0][1]
	property, err := list.resolvePropAcrossProjects(name)
	if err != nil {
		return nil, err
	}
	if property == nil {
		return nil, ErrDepVersionNotFound
	}
	//dep.Version = property.Value

	return property, nil
}

func (list *ProjectList) resolvePropAcrossProjects(name string) (*Property, error) {
	var value *string
	prop := NewProperty(name)

	// Search the property across all POMs.
	for _, project := range list.projects {
		project := project
		list.logger.Debug().Msg("resolve version")

		// Fallback to profile properties.
		if project.Properties == nil {
			var err error
			value, err = project.resolvePropertyFromProfiles(name)
			if err != nil || value == nil {
				list.logger.Debug().Err(err).Str("project", project.Name).Str("property", name).Msg("resolve version from profiles")
				continue
			}
			list.logger.Info().Str("project", project.Name).Str("property", name).Msg("resolved version from property")
			prop.Metadata.DeclarePath = project.pomPath
			break
		}
		if project.Properties.Entries == nil {
			var err error
			value, err = project.resolvePropertyFromProfiles(name)
			if err != nil || value == nil {
				list.logger.Debug().Err(err).Str("project", project.Name).Str("property", name).Msg("resolve version from profiles")
				continue
			}
			list.logger.Info().Str("project", project.Name).Str("property", name).Msg("resolved version from property")
			prop.Metadata.DeclarePath = project.pomPath
			break
		}

		var err error
		value, err = resolvePropertyFromProperties(project.Properties.Entries, name)
		if err != nil || value == nil {
			list.logger.Debug().Err(err).Str("project", project.Name).Str("property", name).Msg("resolve version")
			continue
		}
		// TODO: pick the parent POM instead of the first match.
		list.logger.Debug().Str("project", project.Name).Str("property", name).Msg("resolved version from property")
		prop.Metadata.DeclarePath = project.pomPath
		break
	}
	if value == nil {
		return nil, ErrPropNotFound
	}

	prop.Value = *value

	return prop, nil
}

// TODO: return the profile the property is declared in, besides its value.
func (project *Project) resolvePropertyFromProfiles(prop string) (*string, error) {
	var value *string
	if project.Profiles == nil {
		return nil, ErrProfilesNotFound
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
			project.logger.Debug().Err(err).Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolve version from profiles")
			continue
		}
		if value != nil {
			project.logger.Debug().Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolved version from profile property")
			// TODO: pick the default profile instead of the first match.
			break
		}
	}
	if value == nil {
		return nil, ErrPropNotFoundInProfiles
	}

	return value, nil
}

func (list *ProjectList) resolvePropertyFromProfiles(prop string) (*string, error) {
	var value *string
nextPom:
	for _, project := range list.projects {
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
				list.logger.Debug().Err(err).Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolve version from profiles")
				continue
			}
			if value != nil {
				list.logger.Info().Str("project", project.Name).Str("profile", profile.ID).Str("property", prop).Msg("resolved version from profile property")
				// TODO: pick the default profile instead of the first match.
				break nextPom
			}
		}
	}
	if value == nil {
		return nil, ErrPropNotFoundInProfiles
	}

	return value, nil
}

func resolvePropertyFromProperties(entries map[string]string, property string) (*string, error) {
	if entries == nil {
		return nil, ErrPropEntriesEmpty
	}

	value, ok := entries[property]
	if !ok {
		return nil, ErrPropNotFound
	}
	if value == "" {
		return nil, ErrPropValueEmpty
	}

	return &value, nil
}
