package main

import (
	"fmt"
	"os"
	"syscall"

	"github.com/namsral/flag"
	"github.com/pkg/errors"

	"github.com/fxkr/openview/backend"
	"github.com/fxkr/openview/backend/util/profiling"
	"github.com/fxkr/openview/backend/util/safe"
)

func main() {
	err := run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%+v", errors.WithStack(err))
		os.Exit(1)
	}
}

func run() error {
	fs := flag.NewFlagSetWithEnvPrefix(os.Args[0], "OPENVIEW", flag.ContinueOnError)
	fs.String(flag.DefaultConfigFlagname, "", "path to config file (read-only)")

	var resourcedir = fs.String("resourcedir", "", "path to resource files (read-only)")
	var cachedir = fs.String("cachedir", "", "path to cache directory (read-write)")
	var imagedir = fs.String("imagedir", "", "path to image files (read-only)")

	var cpuprofile = fs.String("cpuprofile", "", "write cpu profile `file` (write-only)")
	var memprofile = fs.String("memprofile", "", "write memory profile to `file` (write-only)")

	var listen = fs.String("listen", ":3000", "`address:port` to listen on")

	err := fs.Parse(os.Args[1:])
	if err != nil {
		return errors.WithStack(err)
	}

	if *resourcedir == "" {
		return errors.New("-resourcedir is mandatory")
	}
	if *cachedir == "" {
		return errors.New("-cachedir is mandatory")
	}
	if *imagedir == "" {
		return errors.New("-imagedir is mandatory")
	}
	if *memprofile != "" {
		profiling.SupportMemoryProfiling(*memprofile, syscall.SIGUSR1)
	}
	if *cpuprofile != "" {
		profiling.SupportCPUProfiling(*cpuprofile, syscall.SIGUSR2)
	}

	app, err := backend.NewApplication(&backend.Config{
		ResourceDir: safe.UnsafeNewPath(*resourcedir),
		CacheDir:    safe.UnsafeNewPath(*cachedir),
		ImageDir:    safe.UnsafeNewPath(*imagedir),

		ListenAddress: *listen,
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return errors.WithStack(app.Run())
}
