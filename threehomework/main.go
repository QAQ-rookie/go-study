package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	group, ctx := errgroup.WithContext(context.Background())
	group.Go(func() error {
		s := http.Server{
			Addr: ":9090",
		}

		go func() {
			<-ctx.Done()
			fmt.Println("stop http server")
			s.Shutdown(context.Background())
		}()
		return s.ListenAndServe()
	})
	group.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
		s := make(chan os.Signal, len(exitSignals))
		signal.Notify(s, exitSignals...)
		for {
			fmt.Println("signal")
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case <-s:
				fmt.Println("signal ctx exit", s)
				return errors.New("signal ctx exit")
			}
		}
	})

	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}
}
