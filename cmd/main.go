package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill, syscall.SIGINT)

	go func() {
		log.Printf("system call:%+v", <-ctx.Done())

		stop()
	}()

	server(ctx)
}
