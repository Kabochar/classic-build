package main

import (
	"context"
	"flag"

	"v1/client"
	"v1/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const version = "v1"

func main() {
	flag.Parse()

	root := &cobra.Command{
		Use:     "chat",
		Version: version,
		Short:   "chat demo",
	}
	ctx := context.Background()

	root.AddCommand(server.NewServerStartCmd(ctx, version))
	root.AddCommand(client.NewCmd(ctx))

	if err := root.Execute(); err != nil {
		logrus.WithError(err).Fatal("Could not run command")
	}
}
