package cronjob

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	s.s = cron.New(opts...)
	return &s
}

func (s *CronjobServer) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	return s.s.AddFunc(spec, cmd)

}

func (s *CronjobServer) Run() error {
	s.ctx, s.cancel = context.WithCancel(context.Background())

	eg, ctx := errgroup.WithContext(s.ctx)
	eg.Go(func() error {
		s.s.Run()
		return nil
	})

	eg.Go(func() error {
		<-ctx.Done()
		s.Stop()
		s.cancel()
		return nil
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	eg.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-c:
				s.Stop()
				s.cancel()
				return nil

			}
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (s *CronjobServer) Stop() {
	ctx := s.s.Stop()
	select {
	case <-ctx.Done():
		log.Println("cron stoped")
	case <-time.After(3 * time.Second):
		log.Println("waiting too lang,killed")
	}

}
