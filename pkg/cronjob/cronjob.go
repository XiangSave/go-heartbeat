package cronjob

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/robfig/cron/v3"
	"golang.org/x/sync/errgroup"
)

type CronjobServer struct {
	ctx    context.Context
	cancel context.CancelFunc
	s      *cron.Cron
}

func New(opts ...cron.Option) *CronjobServer {
	var s CronjobServer
	s.s = cron.New()
	return &s
}

func (s *CronjobServer) Run() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	eg, ctx := errgroup.WithContext(s.ctx)
	eg.Go(func() error {
		defer fmt.Println("Listen defer")
		return s.s.Run()
	})

	eg.Go(func() error {
		<-ctx.Done()
		return s.s.Stop()
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				return s.s.Stop()
				// s.cancel()

			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nkl
}
