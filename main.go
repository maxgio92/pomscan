package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/maxgio92/gopom"
	"github.com/maxgio92/pomscan/internal/files"
)

const (
	pomFile = "pom.xml"
)

func main() {
	var path, artifactId, groupId string
	flag.StringVar(&path, "project-path", ".", "The path to the project")
	flag.StringVar(&artifactId, "artifact-id", "", "The artifactId of the dependency")
	flag.StringVar(&groupId, "group-id", "", "The groupId of the dependency")
	flag.Parse()
	if groupId == "" || artifactId == "" {
		fmt.Println("Dependency id missing")
		flag.Usage()
		os.Exit(1)
	}

	files, err := files.FindFiles(".", pomFile)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		pom, err := gopom.Parse(file)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		dep, err := pom.Search(groupId, artifactId)
		if err == nil {
			fmt.Printf("%s.%s found\n", groupId, artifactId)
			fmt.Println("version:", dep.Version)
			fmt.Println("scope:", dep.Scope)
			fmt.Println("pom:", file)
		}
	}
}
