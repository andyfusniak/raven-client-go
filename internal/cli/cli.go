package cli

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/andyfusniak/raven-client-go/http"
	"github.com/spf13/cobra"
)

// AppKey context key
type AppKey string

// App cli tool.
type App struct {
	version    string
	endpoint   string
	gitCommit  string
	HTTPClient *http.Client
}

// Config parameters to create a new CLI app.
type Config struct {
	Version    string
	Endpoint   string
	GitCommit  string
	HTTPClient *http.Client
}

// NewApp creates a new CLI application.
func NewApp(c Config) *App {
	return &App{
		version:    c.Version,
		endpoint:   c.Endpoint,
		gitCommit:  c.GitCommit,
		HTTPClient: c.HTTPClient,
	}
}

// Version returns the cli application version.
func (a *App) Version() string {
	return a.version
}

// NewCmdList list sub command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "List resources",
		Aliases: []string{"lists"},
	}
	cmd.AddCommand(NewCmdListTemplates())
	return cmd
}

// NewCmdListTemplates list templates sub command.
func NewCmdListTemplates() *cobra.Command {
	return &cobra.Command{
		Use:   "templates",
		Short: "List templates",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			results, err := app.HTTPClient.ListTemplates(ctx)
			if err != nil {
				return err
			}

			tw := new(tabwriter.Writer).Init(os.Stdout, 0, 8, 2, ' ', 0)
			format := "%s\t%s\t%s\n"
			headers := []interface{}{
				"Template ID",
				"Group ID",
				"Created",
			}

			fmt.Fprintf(tw, format, headers...)
			for _, v := range results {
				var params []interface{}
				params = []interface{}{
					v.ID,
					v.GroupID,
					v.CreatedAt,
				}
				fmt.Fprintf(tw, format, params...)
			}

			if err := tw.Flush(); err != nil {
				return err
			}
			var plural string
			if len(results) > 1 {
				plural = "s"
			}
			fmt.Printf("%d templates%s in your project\n", len(results), plural)
			return nil
		},
	}
}

// NewCmdVersion returns an instance of the version sub command.
func NewCmdVersion(version, gitCommit, endpoint string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Raven CLI Tool version",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// prevent root level PersistentPreRun
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Raven CLI tool %s (build %s) built with endpoint %s\n",
				version, gitCommit, endpoint)
		},
	}
}
