package pom

import (
	"encoding/xml"
	"fmt"
	"os"
	"sync"

	"github.com/pkg/errors"
)

const (
	PomFile = "pom.xml"
)

var (
	ErrDepNotFound = errors.New("dependency not found")
)

type Artifact struct {
	GroupId    string `xml:"groupId"`
	ArtifactId string `xml:"artifactId"`
	Version    string `xml:"version"`
	Scope      string `xml:"scope"`
}

type Project struct {
	XMLName     xml.Name `xml:"project"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Artifact
	Dependencies []Artifact `xml:"dependencies>dependency"`
	depCache map[string]*Artifact
	lock     sync.RWMutex
}

func NewPom(pathname string) (*Project, error) {
	file, err := os.Open(pathname)
	if err != nil {
		fmt.Println("error opening file:", err)
	}
	defer file.Close()

	var pom Project
	decoder := xml.NewDecoder(file)
	err = decoder.Decode(&pom)
	if err != nil {
		return nil, errors.Wrap(err, "error decoding Pom")
	}

	// Create a dependency cache for O(1) search.
	pom.depCache = make(map[string]*Artifact, len(pom.Dependencies))
	for _, dep := range pom.Dependencies {
		pom.depCache[fmt.Sprintf("%s.%s", dep.GroupId, dep.ArtifactId)] = &dep
	}

	return &pom, nil
}

func (p *Project) Dep(groupId, artifactId string) (*Artifact, error) {
	p.lock.RLock()
	dep, ok := p.depCache[fmt.Sprintf("%s.%s", groupId, artifactId)]
	p.lock.RUnlock()
	if ok {
		return dep, nil
	}
	return nil, ErrDepNotFound
}
