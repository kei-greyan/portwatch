package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/user/portwatch/internal/alert"
	"github.com/user/portwatch/internal/config"
	"github.com/user/portwatch/internal/monitor"
	"github.com/user/portwatch/internal/scanner"
	"github.com/user/portwatch/internal/state"
)

func main() {
	cfgPath := flag.String("config", "", "path to config file (optional)")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: failed to load config: %v\n", err)
		os.Exit(1)
	}
	if err := config.Validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: invalid config: %v\n", err)
		os.Exit(1)
	}

	sc := scanner.New(cfg)
	st := state.New(cfg)
	al := alert.New(os.Stdout)
	mon := monitor.New(cfg, sc, st, al)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "portwatch: shutting down")
		cancel()
	}()

	if err := mon.Run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "portwatch: %v\n", err)
		os.Exit(1)
	}
}
