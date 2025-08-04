package ctrfreak

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"golang.org/x/sys/unix"


	containerd "github.com/containerd/containerd/v2/client"
	"github.com/containerd/errdefs"
	"github.com/containerd/log"
)

func canIgnoreSignal(s os.Signal) bool {
	return s == unix.SIGURG
}

type killer interface {
	Kill(context.Context, syscall.Signal, ...containerd.KillOpts) error
}

// ForwardAllSignals forwards signals
func ForwardAllSignals(ctx context.Context, task killer) chan os.Signal {
	sigc := make(chan os.Signal, 128)
	signal.Notify(sigc)
	go func() {
		for s := range sigc {
			if canIgnoreSignal(s) {
				log.L.Debugf("Ignoring signal %s", s)
				continue
			}
			log.L.Debug("forwarding signal ", s)
			if err := task.Kill(ctx, s.(syscall.Signal)); err != nil {
				if errdefs.IsNotFound(err) {
					log.L.WithError(err).Debugf("Not forwarding signal %s", s)
					return
				}
				log.L.WithError(err).Errorf("forward signal %s", s)
			}
		}
	}()
	return sigc
}

// StopCatch stops and closes a channel
func StopCatch(sigc chan os.Signal) {
	signal.Stop(sigc)
	close(sigc)
}