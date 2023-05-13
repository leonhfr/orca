package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"

	"github.com/leonhfr/orca/search"
	"github.com/leonhfr/orca/uci"
)

var (
	name    = "Orca"
	version = "0.0.0"
	author  = "Leon Hollender"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	// graceful shutdown when context canceled
	// sending EOF to the UCI scanner by closing the pipe
	r, pipe := io.Pipe()
	go func() { _, _ = io.Copy(pipe, os.Stdin) }()
	go func() { <-ctx.Done(); pipe.Close() }()

	s := uci.NewState(
		fmt.Sprintf("%s v%s", name, version),
		author,
		os.Stdout,
	)
	uci.Run(ctx, search.Engine{}, r, s)
}
