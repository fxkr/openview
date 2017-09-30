package profiling

import (
	"os"
	"os/signal"
	"runtime/pprof"

	log "github.com/sirupsen/logrus"
)

func SupportCPUProfiling(filename string, onSignal os.Signal) {
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, onSignal)

		var active bool

		for {
			<-ch

			if !active {
				f, err := os.Create(filename)
				if err != nil {
					log.WithError(err).Errorf("Failed to create CPU profile")
					return
				}
				if err := pprof.StartCPUProfile(f); err != nil {
					log.WithError(err).Errorf("Failed to start CPU profile")
					return
				}
				log.Info("Started CPU profile")
			} else {
				pprof.StopCPUProfile()
				log.Info("Stopped CPU profile")
			}

			active = !active
		}
	}()
}
