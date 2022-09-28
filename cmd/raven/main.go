package main

import (
	"context"
	"fmt"
	"os"

	"github.com/andyfusniak/raven-client-go/http"
	"github.com/andyfusniak/raven-client-go/internal/cli"
	"github.com/spf13/cobra"
)

var version string
var endpoint string
var gitCommit string

func main() {

	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "[run] error: %+v", err)
		os.Exit(1)
	}
}

func run() error {
	// allow RAVEN_ENDPOINT override
	if envEndpoint, ok := os.LookupEnv("RAVEN_ENDPOINT"); ok {
		endpoint = envEndpoint
	}

	ravenHTTPClient, err := http.NewClient(http.Config{
		Endpoint: endpoint,
		Timeout:  http.DefaultTimout,
	})
	if err != nil {
		return err
	}

	appv := cli.NewApp(cli.Config{
		Version:    version,
		Endpoint:   endpoint,
		GitCommit:  gitCommit,
		HTTPClient: ravenHTTPClient,
	})

	root := cobra.Command{
		Use:     "raven",
		Short:   "raven is command line tool for managing Raven Mailer projects",
		Version: version,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			ctx := cmd.Context()
			v := ctx.Value(cli.AppKey("app"))
			_ = v.(*cli.App)
		},
	}
	root.AddCommand(cli.NewCmdCreate())
	root.AddCommand(cli.NewCmdDelete())
	root.AddCommand(cli.NewCmdGet())
	root.AddCommand(cli.NewCmdList())
	root.AddCommand(cli.NewCmdUpdate())
	root.AddCommand(cli.NewCmdVersion(version, gitCommit, endpoint))

	ctx := context.WithValue(context.Background(), cli.AppKey("app"), appv)
	if err := root.ExecuteContext(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "%+v\n", err)
		os.Exit(1)
	}
	return nil
}
