package profiling

import (
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"

	log "github.com/sirupsen/logrus"
)

func SupportMemoryProfiling(filename string, onSignal os.Signal) {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, onSignal)

		for {
			<-ch

			f, err := os.Create(filename)
			if err != nil {
				log.WithError(err).Errorf("Failed to create memory profile")
				return
			}
			runtime.GC() // get up-to-date statistics
			if err := pprof.WriteHeapProfile(f); err != nil {
				log.WithError(err).Errorf("Failed to write memory profile")
				return
			}
			f.Close()

			log.WithField("filename", filename).Info("Wrote memory profile")
		}
	}()
}
