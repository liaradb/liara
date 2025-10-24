package application

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func WithSignal(ctx context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	go runOnClose(cancel)
	return ctx, cancel
}

func runOnClose(cancel func()) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	cancel()
}
