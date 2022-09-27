package cli

import (
	"fmt"
	"io"
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
	cmd.AddCommand(NewCmdListProjects())
	cmd.AddCommand(NewCmdListTemplates())
	return cmd
}

// NewCmdListProjects list projects sub command.
func NewCmdListProjects() *cobra.Command {
	return &cobra.Command{
		Use:   "projects",
		Short: "List projects",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			results, err := app.HTTPClient.ListProjects(ctx, "dXjtqjXayte0fvfAXfTv3uyIjSj2")
			if err != nil {
				return err
			}

			format := "%s\t%s\t%s\n"
			headers := []interface{}{"Project ID", "Name", "Created"}
			if err := renderTable(os.Stdout, results, format, headers); err != nil {
				return fmt.Errorf("list projects failed to render table: %+v", err)
			}

			return nil
		},
	}
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

			format := "%s\t%s\t%s\n"
			headers := []interface{}{"Template ID", "Group ID", "Created"}
			if err := renderTable(os.Stdout, results, format, headers); err != nil {
				return fmt.Errorf("list templates failed to render table: %+v", err)
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

func renderTable(w io.Writer, results interface{}, format string, headers []interface{}) error {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)

	fmt.Fprintf(tw, format, headers...)

	switch list := results.(type) {
	case []http.Project:
		for _, v := range list {
			row := []interface{}{
				v.ID,
				v.Name,
				v.CreatedAt,
			}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Template:
		for _, v := range list {
			row := []interface{}{
				v.ID,
				v.GroupID,
				v.CreatedAt,
			}
			fmt.Fprintf(tw, format, row...)
		}

	default:
		return fmt.Errorf("unknown results type")
	}

	if err := tw.Flush(); err != nil {
		return err
	}
	return nil
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
