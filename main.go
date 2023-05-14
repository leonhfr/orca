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

	s := uci.NewState(
		fmt.Sprintf("%s v%s", name, version),
		author,
		os.Stdout,
	)
	uci.Run(ctx, search.Engine{}, os.Stdin, s)
}
