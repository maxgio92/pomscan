package project_test

import (
	"os"
	"testing"

	log "github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"

	"github.com/maxgio92/gopom"
	"github.com/maxgio92/pomscan/pkg/project"
)

func TestPlugin(t *testing.T) {
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
						PomPath: "./testdata/plugins/pom.xml",
						Profile: "",
					},
				},
				{
					Plugin: &gopom.Plugin{
						ArtifactID: "maven-antrun-plugin",
						GroupID:    "org.apache.maven.plugins",
					},
					Metadata: &project.PluginMetadata{
						PomPath: "./testdata/plugins/pom.xml",
						Profile: "dist",
					},
				},
			},
		},
	}
	logger := log.New(os.Stderr)
	projectList := project.NewProjectList(
		project.ListWithLogger(&logger),
		project.ListWithPomPaths("./testdata/plugins/pom.xml"),
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
