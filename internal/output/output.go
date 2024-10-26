package output

import (
	"fmt"
	"github.com/maxgio92/gopom"
)

func PrintDep(dep *gopom.Dependency, pomPath string, versionOnly bool) {
	if versionOnly && dep.Version == "" {
		return
	}

	fmt.Printf("📦 %s.%s found\n", dep.GroupID, dep.ArtifactID)
	fmt.Println("pom:", pomPath)
	if dep.Version != "" {
		fmt.Println("version:", dep.Version)
	}
	if dep.Scope != "" {
		fmt.Println("scope:", dep.Scope)
	}
	fmt.Println()
}
