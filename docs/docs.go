package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/maxgio92/pomscan/cmd"
	"github.com/maxgio92/pomscan/internal/options"
	log "github.com/rs/zerolog"
	"github.com/spf13/cobra/doc"
)

const (
	cmdline      = "pomscan"
	docsDir      = "docs"
	fileTemplate = `---
title: %s
---	

`
)

var (
	filePrepender = func(filename string) string {
		title := strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(filename, "_", " "), ".md"), fmt.Sprintf("%s/", docsDir))
		return fmt.Sprintf(fileTemplate, title)
	}
	linkHandler = func(filename string) string {
		if filename == cmdline+".md" {
			return "_index.md"
		}
		return filename
	}
)

func main() {
	logger := log.New(os.Stderr)
	if err := doc.GenMarkdownTreeCustom(
		cmd.NewRootCommand(options.NewCommonOptions(options.WithLogger(&logger))),
		docsDir,
		filePrepender,
		linkHandler,
	); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err := os.Rename(path.Join(docsDir, cmdline+".md"), path.Join(docsDir, "_index.md"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
