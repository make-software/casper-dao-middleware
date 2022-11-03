package exec

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func RunGracefulShutDownListener(ctx context.Context, cancel context.CancelFunc) context.Context {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		osCall := <-c
		log.Printf("Stop api system call:%+v", osCall)
		cancel()
	}()

	return ctx
}
