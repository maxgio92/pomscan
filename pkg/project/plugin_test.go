package project_test

import (
	"os"
	"testing"

	log "github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/maxgio92/gopom"
	"github.com/maxgio92/pomscan/pkg/project"
)

func TestSearchPlugin(t *testing.T) {
	pomPath := "./testdata/plugin/pom.xml"
	testCases := []struct {
		artifactID string
		groupID    string
		found      []project.Plugin
	}{
		{
			artifactID: "maven-antrun-plugin",
			groupID:    "org.apache.maven.plugins",
			found: []project.Plugin{
				{
					Plugin: &gopom.Plugin{
						ArtifactID: "maven-antrun-plugin",
						GroupID:    "org.apache.maven.plugins",
					},
					Metadata: &project.PluginMetadata{
						PomPath: pomPath,
						Profile: "",
					},
				},
				{
					Plugin: &gopom.Plugin{
						ArtifactID: "maven-antrun-plugin",
						GroupID:    "org.apache.maven.plugins",
					},
					Metadata: &project.PluginMetadata{
						PomPath: pomPath,
						Profile: "dist",
					},
				},
			},
		},
	}
	logger := log.New(os.Stderr)
	projectList := project.NewProjectList(
		project.ListWithLogger(&logger),
		project.ListWithPomPaths(pomPath),
	)
	err := projectList.LoadAll()
	assert.Nil(t, err)

	for _, tt := range testCases {
		found, err := projectList.SearchPlugins(tt.artifactID, tt.groupID)
		assert.Nil(t, err)
		assert.Equal(t, len(tt.found), len(found))
		for k, plugin := range found {
			assert.Equal(t, tt.found[k].ArtifactID, plugin.ArtifactID)
			assert.Equal(t, tt.found[k].GroupID, plugin.GroupID)
			assert.Equal(t, tt.found[k].Metadata.PomPath, plugin.Metadata.PomPath)
			assert.Equal(t, tt.found[k].Metadata.Profile, plugin.Metadata.Profile)
		}
	}
}

func TestSearchPluginDependency(t *testing.T) {
	pomPath := "./testdata/plugin-dependency/pom.xml"
	testCases := []struct {
		artifactID string
		groupID    string
		found      []project.Dependency
	}{
		{
			artifactID: "ant-contrib",
			groupID:    "ant-contrib",
			found: []project.Dependency{
				{
					Dependency: &gopom.Dependency{
						ArtifactID: "ant-contrib",
						GroupID:    "ant-contrib",
					},
					Metadata: &project.DependencyMetadata{
						PomPath: pomPath,
						VersionProperty: &project.Property{
							Name:  "ant.contrib.version",
							Value: "1.0b3",
							Metadata: &project.PropertyMetadata{
								DeclarePath: pomPath,
							},
						},
					},
				},
			},
		},
	}
	logger := log.New(os.Stderr)
	projectList := project.NewProjectList(
		project.ListWithLogger(&logger),
		project.ListWithPomPaths(pomPath),
	)
	err := projectList.LoadAll()
	assert.Nil(t, err)

	for _, tt := range testCases {
		found, err := projectList.SearchPluginDependency(tt.artifactID, tt.groupID)
		assert.Nil(t, err)
		assert.Equal(t, len(tt.found), len(found))
		for k, plugin := range found {
			assert.Equal(t, tt.found[k].ArtifactID, plugin.ArtifactID)
			assert.Equal(t, tt.found[k].GroupID, plugin.GroupID)
			assert.NotNil(t, plugin.Metadata)
			if plugin.Metadata != nil {
				assert.Equal(t, tt.found[k].Metadata.PomPath, plugin.Metadata.PomPath)
				assert.NotNil(t, plugin.Metadata.VersionProperty)
				if plugin.Metadata.VersionProperty != nil {
					assert.Equal(t, tt.found[k].Metadata.VersionProperty.Name, plugin.Metadata.VersionProperty.Name)
					assert.Equal(t, tt.found[k].Metadata.VersionProperty.Value, plugin.Metadata.VersionProperty.Value)
					assert.NotNil(t, tt.found[k].Metadata.VersionProperty.Metadata, plugin.Metadata.VersionProperty.Metadata)
					if plugin.Metadata.VersionProperty.Metadata != nil {
						assert.Equal(t, tt.found[k].Metadata.VersionProperty.Metadata.DeclarePath,
							plugin.Metadata.VersionProperty.Metadata.DeclarePath)
					}
				}
			}
		}
	}
}
