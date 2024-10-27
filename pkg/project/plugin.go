package project

import (
	"fmt"
	"github.com/maxgio92/gopom"
)

var ()

// Plugin represents a build-time plugin.
type Plugin struct {
	*gopom.Plugin
	Metadata *PluginMetadata
}

type PluginMetadata struct {
	PomPath         string
	Profile         string
	VersionProperty *Property
}

func (list *ProjectList) SearchPlugins(artifactID, groupID string) ([]*Plugin, error) {
	result := make([]*Plugin, 0)
	for _, project := range list.projects {
		project := project

		// Not inherited plugins.
		// TODO: search in build.plugins
		plugin, err := project.searchPluginInBuild(artifactID, groupID)
		if err != nil {
			list.logger.Debug().Err(err).Str("project", project.Name).Msg("search plugin in build")
		}
		if plugin != nil {
			result = append(result, plugin)
		}

		// TODO: search in profiles.build.plugins
		plugin, err = project.searchPluginInProfileBuild(artifactID, groupID)
		if err != nil {
			list.logger.Debug().Err(err).Str("project", project.Name).Msg("search plugin in profiles build")
		}
		if plugin != nil {
			result = append(result, plugin)
		}

		// Inheritable and overridable plugins.
		// TODO: search in build.pluginManagement.plugins
		plugin, err = project.searchPluginInBuildPluginMgmt(artifactID, groupID)
		if err != nil {
			list.logger.Debug().Err(err).Str("project", project.Name).Msg("search plugin in build plugin management")
		}
		if plugin != nil {
			result = append(result, plugin)
		}

		// TODO: search in profiles.build.pluginManagement.plugins
		plugin, err = project.searchPluginInProfileBuildPluginMgmt(artifactID, groupID)
		if err != nil {
			list.logger.Debug().Err(err).Str("project", project.Name).Msg("search plugin in profile build plugin management")
		}
		if plugin != nil {
			result = append(result, plugin)
		}
	}
	if len(result) == 0 {
		return nil, ErrPluginNotFound
	}

	return result, nil
}

func (p *Project) searchPluginInBuild(artifactID, groupID string) (*Plugin, error) {
	if p.Build == nil {
		return nil, ErrBuildEmpty
	}
	if p.Build.Plugins == nil {
		return nil, ErrBuildPluginsEmpty
	}
	for _, plugin := range *p.Build.Plugins {
		if (groupID == "" && plugin.ArtifactID == artifactID) ||
			(groupID != "" && plugin.GroupID == groupID && plugin.ArtifactID == artifactID) {
			metadata := &PluginMetadata{
				PomPath: p.pomPath,
			}

			return &Plugin{
				Plugin:   &plugin,
				Metadata: metadata,
			}, nil
		}
	}

	return nil, ErrPluginNotFound
}

func (p *Project) searchPluginInBuildPluginMgmt(artifactID, groupID string) (*Plugin, error) {
	if p.Build == nil {
		return nil, ErrBuildEmpty
	}
	if p.Build.PluginManagement == nil {
		return nil, ErrBuildPluginMgmtEmpty
	}
	if p.Build.PluginManagement.Plugins == nil {
		return nil, ErrBuildPluginMgmtPluginsEmpty
	}
	for _, plugin := range *p.Build.PluginManagement.Plugins {
		if (groupID == "" && plugin.ArtifactID == artifactID) ||
			(groupID != "" && plugin.GroupID == groupID && plugin.ArtifactID == artifactID) {
			metadata := &PluginMetadata{
				PomPath: p.pomPath,
			}

			return &Plugin{
				Plugin:   &plugin,
				Metadata: metadata,
			}, nil
		}
	}

	return nil, ErrPluginNotFound
}

func (p *Project) searchPluginInProfileBuild(artifactID, groupID string) (*Plugin, error) {
	if p.Profiles == nil {
		return nil, ErrProfilesEmpty
	}
	for _, profile := range *p.Profiles {
		profile := profile
		if profile.Build == nil {
			continue
		}
		if profile.Build.Plugins == nil {
			continue
		}
		for _, plugin := range *profile.Build.Plugins {
			if (groupID == "" && plugin.ArtifactID == artifactID) ||
				(groupID != "" && plugin.GroupID == groupID && plugin.ArtifactID == artifactID) {
				metadata := &PluginMetadata{
					PomPath: p.pomPath,
					Profile: profile.ID,
				}

				return &Plugin{
					Plugin:   &plugin,
					Metadata: metadata,
				}, nil
			}
		}
	}

	return nil, ErrPluginNotFound
}

func (p *Project) searchPluginInProfileBuildPluginMgmt(artifactID, groupID string) (*Plugin, error) {
	if p.Profiles == nil {
		return nil, ErrProfilesEmpty
	}
	for _, profile := range *p.Profiles {
		if profile.Build == nil {
			continue
		}
		if profile.Build.PluginManagement == nil {
			continue
		}
		if profile.Build.PluginManagement.Plugins == nil {
			continue
		}
		for _, plugin := range *profile.Build.PluginManagement.Plugins {
			if (groupID == "" && plugin.ArtifactID == artifactID) ||
				(groupID != "" && plugin.GroupID == groupID && plugin.ArtifactID == artifactID) {
				metadata := &PluginMetadata{
					PomPath: p.pomPath,
					Profile: profile.ID,
				}

				return &Plugin{
					Plugin:   &plugin,
					Metadata: metadata,
				}, nil
			}
		}
	}

	return nil, ErrPluginNotFound
}

func (list *ProjectList) SearchPluginDependency(artifactID, groupID string) ([]*Dependency, error) {
	result := make([]*Dependency, 0)

	deps, err := list.searchPluginDepInBuild(artifactID, groupID)
	if err != nil {
		list.logger.Debug().Err(err).Msg("search plugin dependency in build")
	}
	if deps != nil {
		result = append(result, deps...)
	}

	deps, err = list.searchPluginDepInProfileBuild(artifactID, groupID)
	if err != nil {
		list.logger.Debug().Err(err).Msg("search plugin dependency in profile build")
	}
	if deps != nil {
		result = append(result, deps...)
	}

	deps, err = list.searchPluginDepInBuildPluginMgmt(artifactID, groupID)
	if err != nil {
		list.logger.Debug().Err(err).Msg("search plugin dependency in build plugin management")
	}
	if deps != nil {
		result = append(result, deps...)
	}

	deps, err = list.searchPluginDepInProfileBuildPluginMgmt(artifactID, groupID)
	if err != nil {
		list.logger.Debug().Err(err).Msg("search plugin dependency in profile build plugin management")
	}
	if deps != nil {
		result = append(result, deps...)
	}

	if len(result) == 0 {
		return nil, ErrDepNotFound
	}

	return result, nil
}

func (list *ProjectList) searchPluginDepInBuild(artifactID, groupID string) ([]*Dependency, error) {
	dependencies := make([]*Dependency, 0)
	for _, project := range list.projects {
		if project.Build == nil {
			return nil, ErrBuildEmpty
		}
		if project.Build.Plugins == nil {
			return nil, ErrBuildPluginsEmpty
		}
		for _, plugin := range *project.Build.Plugins {
			if plugin.Dependencies == nil {
				continue
			}
			for _, dependency := range *plugin.Dependencies {
				if (groupID == "" && dependency.ArtifactID == artifactID) ||
					(groupID != "" && dependency.GroupID == groupID && artifactID == dependency.ArtifactID) {

					// Resolve the version of the dependency.
					property, err := list.ResolveDepVersionProp(&dependency)
					if err != nil {
						list.logger.Debug().Err(err).
							Str("project", project.Name).
							Str("plugin", fmt.Sprintf("%s:%s", plugin.GroupID, plugin.ArtifactID)).
							Str("property", dependency.Version).
							Msg("cannot resolve version property in plugin dependency in build")
					}

					dependencies = append(dependencies, &Dependency{
						Dependency: &dependency,
						Metadata: &DependencyMetadata{
							PomPath:         project.pomPath,
							VersionProperty: property,
						},
					})
				}
			}
		}
	}
	if len(dependencies) == 0 {
		return nil, ErrDepNotFound
	}

	return dependencies, nil
}

func (list *ProjectList) searchPluginDepInBuildPluginMgmt(artifactID, groupID string) ([]*Dependency, error) {
	dependencies := make([]*Dependency, 0)
	for _, project := range list.projects {
		if project.Build == nil {
			return nil, ErrBuildEmpty
		}
		if project.Build.PluginManagement == nil {
			return nil, ErrBuildPluginMgmtEmpty
		}
		if project.Build.PluginManagement.Plugins == nil {
			return nil, ErrBuildPluginMgmtPluginsEmpty
		}
		for _, plugin := range *project.Build.PluginManagement.Plugins {
			if plugin.Dependencies == nil {
				continue
			}
			for _, dependency := range *plugin.Dependencies {
				if (groupID == "" && dependency.ArtifactID == artifactID) ||
					(groupID != "" && dependency.GroupID == groupID && artifactID == dependency.ArtifactID) {

					// Resolve the version of the dependency.
					property, err := list.ResolveDepVersionProp(&dependency)
					if err != nil {
						list.logger.Debug().Err(err).
							Str("project", project.Name).
							Str("plugin", fmt.Sprintf("%s:%s", plugin.GroupID, plugin.ArtifactID)).
							Str("property", dependency.Version).
							Msg("cannot resolve version property in plugin dependency in build plugin management")
					}

					dependencies = append(dependencies, &Dependency{
						Dependency: &dependency,
						Metadata: &DependencyMetadata{
							PomPath:         project.pomPath,
							VersionProperty: property,
						},
					})
				}
			}
		}
	}
	if len(dependencies) == 0 {
		return nil, ErrDepNotFound
	}

	return dependencies, nil
}

func (list *ProjectList) searchPluginDepInProfileBuild(artifactID, groupID string) ([]*Dependency, error) {
	dependencies := make([]*Dependency, 0)
	for _, project := range list.projects {
		if project.Profiles == nil {
			return nil, ErrProfilesEmpty

		}
		for _, profile := range *project.Profiles {
			profile := profile
			if profile.Build == nil {
				continue
			}
			if profile.Build.Plugins == nil {
				continue
			}
			for _, plugin := range *profile.Build.Plugins {
				if plugin.Dependencies == nil {
					continue
				}
				for _, dependency := range *plugin.Dependencies {
					if (groupID == "" && dependency.ArtifactID == artifactID) ||
						(groupID != "" && dependency.GroupID == groupID && artifactID == dependency.ArtifactID) {

						// Resolve the version of the dependency.
						property, err := list.ResolveDepVersionProp(&dependency)
						if err != nil {
							list.logger.Debug().Err(err).
								Str("project", project.Name).
								Str("plugin", fmt.Sprintf("%s:%s", plugin.GroupID, plugin.ArtifactID)).
								Str("property", dependency.Version).
								Msg("cannot resolve version property in plugin dependency in build plugin management")
						}

						dependencies = append(dependencies, &Dependency{
							Dependency: &dependency,
							Metadata: &DependencyMetadata{
								PomPath:         project.pomPath,
								VersionProperty: property,
							},
						})
					}
				}
			}
		}
	}
	if len(dependencies) == 0 {
		return nil, ErrDepNotFound
	}

	return dependencies, nil
}

func (list *ProjectList) searchPluginDepInProfileBuildPluginMgmt(artifactID, groupID string) ([]*Dependency, error) {
	dependencies := make([]*Dependency, 0)
	for _, project := range list.projects {
		if project.Profiles == nil {
			return nil, ErrProfilesEmpty

		}
		for _, profile := range *project.Profiles {
			profile := profile
			if profile.Build == nil {
				continue
			}
			if profile.Build.PluginManagement == nil {
				continue
			}
			if profile.Build.PluginManagement.Plugins == nil {
				continue
			}
			for _, plugin := range *profile.Build.PluginManagement.Plugins {
				if plugin.Dependencies == nil {
					continue
				}
				for _, dependency := range *plugin.Dependencies {
					if (groupID == "" && dependency.ArtifactID == artifactID) ||
						(groupID != "" && dependency.GroupID == groupID && artifactID == dependency.ArtifactID) {

						// Resolve the version of the dependency.
						property, err := list.ResolveDepVersionProp(&dependency)
						if err != nil {
							list.logger.Debug().Err(err).
								Str("project", project.Name).
								Str("plugin", fmt.Sprintf("%s:%s", plugin.GroupID, plugin.ArtifactID)).
								Str("property", dependency.Version).
								Msg("cannot resolve version property in plugin dependency in build plugin management")
						}

						dependencies = append(dependencies, &Dependency{
							Dependency: &dependency,
							Metadata: &DependencyMetadata{
								PomPath:         project.pomPath,
								VersionProperty: property,
							},
						})
					}
				}
			}
		}
	}
	if len(dependencies) == 0 {
		return nil, ErrDepNotFound
	}

	return dependencies, nil
}
