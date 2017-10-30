package main

import (
	"flag"
	"fmt"
	"os"
	"syscall"

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
	var resourcedir = flag.String("resourcedir", "", "path to resource files (read-only)")
	var cachedir = flag.String("cachedir", "", "path to cache directory (read-write)")
	var imagedir = flag.String("imagedir", "", "path to image files (read-only)")

	var cpuprofile = flag.String("cpuprofile", "", "write cpu profile `file` (write-only)")
	var memprofile = flag.String("memprofile", "", "write memory profile to `file` (write-only)")

	var listen = flag.String("listen", ":3000", "`address:port` to listen on")

	flag.Parse()

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
