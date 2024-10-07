package options

import (
	log "github.com/rs/zerolog"
	"github.com/spf13/pflag"
)

type CommonOptions struct {
	ProjectPath string
	Logger      *log.Logger
	Debug       bool
}

type CommonOption func(opts *CommonOptions)

func NewCommonOptions(opts ...CommonOption) *CommonOptions {
	o := new(CommonOptions)

	for _, f := range opts {
		f(o)
	}

	return o
}

func WithProjectPath(path string) CommonOption {
	return func(opts *CommonOptions) {
		opts.ProjectPath = path
	}
}

func WithLogger(logger *log.Logger) CommonOption {
	return func(opts *CommonOptions) {
		opts.Logger = logger
	}
}

func (o *CommonOptions) AddFlags(flags *pflag.FlagSet) {
	flags.StringVarP(&o.ProjectPath, "project-path", "p", ".", "Project path")
}
