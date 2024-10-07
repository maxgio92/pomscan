package pom_test

import (
	"testing"

	"github.com/maxgio92/pomscan/pkg/pom"
)

func TestNewPom(t *testing.T) {
	testCases := []struct {
		name     string
		expected *pom.Project
	}{
		{name: "testdata/pom.xml", expected: &pom.Project{
			Name:        "My Project",
			Description: "A simple Maven project",
			Artifact: pom.Artifact{
				GroupId:    "com.example.my-project",
				ArtifactId: "my-project",
				Version:    "1.0-SNAPSHOT",
			},
			Dependencies: []pom.Artifact{
				{
					GroupId:    "org.apache.foo",
					ArtifactId: "foo",
					Version:    "1.0.1",
					Scope:      "",
				},
				{
					GroupId:    "org.apache.bar",
					ArtifactId: "bar",
					Version:    "1.0.2",
					Scope:      "",
				},
			},
		}},
		{name: "testdata/subproject/pom.xml", expected: &pom.Project{
			Name:        "My sub project",
			Description: "A simple Maven project",
			Artifact: pom.Artifact{
				GroupId:    "com.example.my-sub-project",
				ArtifactId: "my-sub-project",
				Version:    "1.0-SNAPSHOT",
			},
			Dependencies: []pom.Artifact{
				{
					GroupId:    "org.apache.foo",
					ArtifactId: "foo",
					Version:    "1.0.1",
					Scope:      "",
				},
				{
					GroupId:    "org.apache.bar",
					ArtifactId: "bar",
					Version:    "1.0.2",
					Scope:      "",
				},
			},
		}},
		{name: "testdata/subproject/subsubproject/pom.xml", expected: &pom.Project{
			Name:        "My sub sub project",
			Description: "A simple Maven project",
			Artifact: pom.Artifact{
				GroupId:    "com.example.my-sub-sub-project",
				ArtifactId: "my-sub-sub-project",
				Version:    "1.0-SNAPSHOT",
			},
			Dependencies: []pom.Artifact{
				{
					GroupId:    "org.apache.foo",
					ArtifactId: "foo",
					Version:    "1.0.1",
					Scope:      "",
				},
				{
					GroupId:    "org.apache.bar",
					ArtifactId: "bar",
					Version:    "1.0.2",
					Scope:      "",
				},
			},
		}},
	}

	for _, tc := range testCases {
		pom, err := pom.NewPom(tc.name)
		if err != nil {
			t.Fatalf("should not fail creating pom from %s", tc.name)
		}
		if pom.GroupId != tc.expected.GroupId {
			t.Fatalf("expected %s GroupId, found %s", tc.expected.GroupId, pom.GroupId)
		}
		if pom.ArtifactId != tc.expected.ArtifactId {
			t.Fatalf("expected %s ArtifactId, found %s", tc.expected.ArtifactId, pom.ArtifactId)
		}
		if pom.Version != tc.expected.Version {
			t.Fatalf("expected %s Version, found %s", tc.expected.Version, pom.Version)
		}
		if len(pom.Dependencies) != len(tc.expected.Dependencies) {
			t.Fatalf("expected %d deps, found %d", len(tc.expected.Dependencies), len(pom.Dependencies))
		}
		for i, dep := range pom.Dependencies {
			if dep.GroupId != tc.expected.Dependencies[i].GroupId {
				t.Fatalf("expected dep %s GroupId, found %s", tc.expected.Dependencies[i].GroupId, dep.GroupId)
			}
			if dep.ArtifactId != tc.expected.Dependencies[i].ArtifactId {
				t.Fatalf("expected dep %s ArtifactId, found %s", tc.expected.Dependencies[i].ArtifactId, dep.ArtifactId)
			}
			if dep.Version != tc.expected.Dependencies[i].Version {
				t.Fatalf("expected dep %s Version, found %s", tc.expected.Dependencies[i].Version, dep.Version)
			}
			if dep.Scope != tc.expected.Dependencies[i].Scope {
				t.Fatalf("expected dep %s Scope, found %s", tc.expected.Dependencies[i].Scope, dep.Scope)
			}
		}
	}
}
