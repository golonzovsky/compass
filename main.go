package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"
	"github.com/golonzovsky/comPass/cmd"
)

func main() {
	errorLogger := log.New(log.WithOutput(os.Stderr))

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() { signal.Stop(c) }()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		select {
		case <-c: // first signal, cancel context
			cancel()
		case <-ctx.Done():
		}
		<-c // second signal, hard exit
		errorLogger.Error("second interrupt, exiting")
		os.Exit(1)
	}()

	if err := cmd.NewRootCmd().ExecuteContext(ctx); err != nil {
		errorLogger.Error(err.Error())
		os.Exit(1)
	}

}
