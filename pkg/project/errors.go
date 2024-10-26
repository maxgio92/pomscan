package project

import (
	"github.com/pkg/errors"
)

var (
	ErrDepVersionNotFound     = errors.New("dependency version not found")
	ErrPomDepMgmtEmpty        = errors.New("pom dependency management is empty")
	ErrPropNotFoundInProfiles = errors.New("property not found in profiles")
	ErrPropNotFound           = errors.New("property not found")
	ErrPropValueEmpty         = errors.New("property value is empty")
	ErrPropEntriesEmpty       = errors.New("properties entry is empty")
	ErrProfilesNotFound       = errors.New("profiles are not found")
)
