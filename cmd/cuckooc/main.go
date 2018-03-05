package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/vedhavyas/cuckooc"
)

var cfile = flag.String("c", "", "path to configuration file")

func main() {
	flag.Parse()
	if *cfile == "" {
		fmt.Println("Path to config file is required")
		flag.Usage()
		os.Exit(1)
	}

	c, err := cuckooc.LoadConfig(*cfile)
	if err != nil {
		log.Fatalf("unable to load config: %v\n", err)
	}

	wg := new(sync.WaitGroup)
	defer wg.Wait()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cmdCh := make(chan cuckooc.Executor)
	gk := cuckooc.NewGatekeeper(cmdCh)
	wg.Add(1)
	go gk.Start(ctx, c, wg)

	if c.TCP != "" {
		wg.Add(1)
		go cuckooc.StartTCPServer(ctx, c, wg, cmdCh)
	}

	sigInt := make(chan os.Signal, 1)
	signal.Notify(sigInt, os.Interrupt)
	for range sigInt {
		log.Println("stopping all services...")
		break
	}
}
