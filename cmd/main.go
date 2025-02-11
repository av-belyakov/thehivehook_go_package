package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
)

func main() {
	var fmempprof *os.File
	var memprofile = flag.String("memprofile", "", "write memory profile to `file`")

	flag.Parse()

	if *memprofile != "" {
		log.Printf("Start profiling a Go program, file name '%s'", *memprofile)

		fmempprof, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer fmempprof.Close()

		runtime.GC()
		if err := pprof.WriteHeapProfile(fmempprof); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		sigChan := make(chan os.Signal, 1)
		osCall := <-sigChan
		log.Printf("system call:%+v", osCall)

		cancel()
		if fmempprof != nil && fmempprof.Close != nil {
			fmempprof.Close()
		}
	}()

	server(ctx)
}
