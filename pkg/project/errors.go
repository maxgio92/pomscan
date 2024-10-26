package project

import (
	"github.com/pkg/errors"
)

var (
	ErrDepVersionNotFound             = errors.New("dependency version not found")
	ErrDepMgmtEmpty                   = errors.New("dependency management is empty")
	ErrBuildEmpty                     = errors.New("build is empty")
	ErrBuildPluginsEmpty              = errors.New("build plugins is empty")
	ErrProfilesEmpty                  = errors.New("profiles is empty")
	ErrBuildPluginMgmtEmpty        = errors.New("plugin management is empty")
	ErrBuildPluginMgmtPluginsEmpty = errors.New("plugin management plugin list is empty")
	ErrPluginNotFound              = errors.New("plugin not found")
	ErrPropNotFoundInProfiles         = errors.New("property not found in profiles")
	ErrPropNotFound                   = errors.New("property not found")
	ErrPropValueEmpty                 = errors.New("property value is empty")
	ErrPropEntriesEmpty               = errors.New("properties entry is empty")
	ErrProfilesNotFound               = errors.New("profiles are not found")
)
