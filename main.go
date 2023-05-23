package main

import (
	"context"
	"fmt"
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

	e := search.NewEngine()
	c := uci.NewController(
		fmt.Sprintf("%s v%s", name, version),
		author,
		os.Stdout,
	)
	c.Run(ctx, e, os.Stdin)
}
