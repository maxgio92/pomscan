package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"k8s.io/kubectl/pkg/cmd/util"

	"github.com/maxgio92/pomscan/cmd/dependency"
	"github.com/maxgio92/pomscan/cmd/plugin"
	"github.com/maxgio92/pomscan/internal/options"
)

var (
	out = bufio.NewWriter(os.Stdout)
)

type Options struct {
	options.CommonOptions
}

func NewRootCommand(opts *options.CommonOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:               "pomscan",
		Short:             "Scan POM files",
		DisableAutoGenTag: true,
	}

	cmd.AddCommand(dependency.NewDepCmd(opts))
	cmd.AddCommand(plugin.NewPluginCmd(opts))
	opts.AddFlags(cmd.PersistentFlags())
	cmd.PersistentFlags().BoolVar(&opts.Debug, "debug", false, "Sets log level to debug")

	return cmd
}

func Execute() {
	// Ensure the buffer is flushed to output when returning.
	defer out.Flush()

	// Mark context done when signals arrive.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)

	// Graceful shutdown.
	go func() {
		<-ctx.Done()
		fmt.Println("signal received, shutting down")
		stop()
	}()

	logger := log.New(os.Stderr).Level(log.InfoLevel)

	opts := options.NewCommonOptions(
		options.WithProjectPath("."),
		options.WithLogger(&logger),
	)

	util.CheckErr(NewRootCommand(opts).ExecuteContext(ctx))
}
