// Package main is the entry point of the engine.
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
	if len(os.Args) > 1 {
		r, err := perft(os.Args[1:])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		fmt.Println(r)
		return
	}

	run()
}

func run() {
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
