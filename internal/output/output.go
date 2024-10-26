package output

import (
	"fmt"

	"github.com/maxgio92/pomscan/pkg/project"
)

const (
	pomFile                    = "pom-file"
	artifactID                 = "artifact-id"
	groupID                    = "group-id"
	scope                      = "scope"
	version                    = "version"
	versionPropertyDeclarePath = "version-property-declare-path"
	versionPropertyName        = "version-property-name"
	versionPropertyValue       = "version-property-value"
)

func PrintDep(dep *project.Dependency, versionOnly bool) {
	if versionOnly && dep.Version == "" {
		return
	}

	fmt.Printf("📦 %s.%s\n", dep.GroupID, dep.ArtifactID)
	fmt.Println(artifactID, ":", dep.ArtifactID)
	fmt.Println(groupID, ":", dep.GroupID)
	fmt.Println(pomFile, ":", dep.Metadata.PomPath)
	if dep.Version != "" {
		fmt.Println(version, ":", dep.Version)
	}
	if dep.Scope != "" {
		fmt.Println(scope, ":", dep.Scope)
	}
	if dep.Metadata.VersionProperty != nil {
		fmt.Println(versionPropertyName, ":", dep.Metadata.VersionProperty.Name)
		fmt.Println(versionPropertyValue, ":", dep.Metadata.VersionProperty.Value)
		if dep.Metadata.VersionProperty.Metadata != nil {
			fmt.Println(versionPropertyDeclarePath, ":", dep.Metadata.VersionProperty.Metadata.DeclarePath)
		}
	}
	fmt.Println()
}

func PrintPlugin(dep *project.Plugin, versionOnly bool) {
	if versionOnly && dep.Version == "" {
		return
	}

	fmt.Printf("📦 %s.%s\n", dep.GroupID, dep.ArtifactID)
	fmt.Println(artifactID, ":", dep.ArtifactID)
	fmt.Println(groupID, ":", dep.GroupID)
	fmt.Println(pomFile, ":", dep.Metadata.PomPath)
	if dep.Version != "" {
		fmt.Println(version, ":", dep.Version)
	}
	if dep.Metadata.VersionProperty != nil {
		fmt.Println(versionPropertyName, ":", dep.Metadata.VersionProperty.Name)
		fmt.Println(versionPropertyValue, ":", dep.Metadata.VersionProperty.Value)
		if dep.Metadata.VersionProperty.Metadata != nil {
			fmt.Println(versionPropertyDeclarePath, ":", dep.Metadata.VersionProperty.Metadata.DeclarePath)
		}
	}
	fmt.Println()
}
