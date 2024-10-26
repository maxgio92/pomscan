package project

import "github.com/maxgio92/gopom"

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
