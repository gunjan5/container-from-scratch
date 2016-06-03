package cmd

import (
	"github.com/codegangsta/cli"
	"github.com/gunjan5/container-from-scratch/container"
	"github.com/gunjan5/container-from-scratch/server"
)

func Serve(ctx *cli.Context) error {
	server.MakeServer()
	return nil

}

func Run(ctx *cli.Context) error {

	return container.Run(ctx.Args()[:])
}

func NewRoot(ctx *cli.Context) error {

	return container.NewRoot(ctx.Args()[:])
}

func Child(ctx *cli.Context) error {

	return container.Child(ctx.Args()[:])
}
