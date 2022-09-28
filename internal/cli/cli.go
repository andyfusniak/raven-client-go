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

// NewCmdCreate create sub command.
func NewCmdCreate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create resource(s)",
	}
	cmd.AddCommand(NewCmdCreateGroup())
	cmd.AddCommand(NewCmdCreateTemplate())
	return cmd
}

// NewCmdList list sub command.
func NewCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List resources",
	}
	cmd.AddCommand(NewCmdListGroups())
	cmd.AddCommand(NewCmdListProjects())
	cmd.AddCommand(NewCmdListTemplates())
	cmd.AddCommand(NewCmdListTransports())
	return cmd
}

// NewCmdGet get sub command.
func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single resource",
	}
	cmd.AddCommand(NewCmdGetGroup())
	cmd.AddCommand(NewCmdGetTemplate())
	return cmd
}

// NewCmdUpdate for updating resources.
func NewCmdUpdate() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update resources",
	}
	// cmd.AddCommand(NewCmdUpdateTemplates())
	return cmd
}

// NewCmdDelete delete sub command.
func NewCmdDelete() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
	}
	cmd.AddCommand(NewCmdDeleteGroup())
	cmd.AddCommand(NewCmdDeleteTemplate())
	return cmd
}

// NewCmdListProjects list projects sub command.
func NewCmdListProjects() *cobra.Command {
	return &cobra.Command{
		Use:     "projects",
		Short:   "List projects",
		Aliases: []string{"project"},
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

// NewCmdListTransports list transports sub command.
func NewCmdListTransports() *cobra.Command {
	return &cobra.Command{
		Use:     "transports",
		Short:   "List transports",
		Aliases: []string{"transport"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			results, err := app.HTTPClient.ListTransports(ctx)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%v\t%s\t%v\n"
			headers := []interface{}{"Transport ID", "Host:Port", "IsActive", "Username", "Created"}
			if err := renderTable(os.Stdout, results, format, headers); err != nil {
				return fmt.Errorf("list transports failed to render table: %+v", err)
			}
			return nil
		},
	}
}

func renderTable(w io.Writer, results interface{}, format string, headers []interface{}) error {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)

	fmt.Fprintf(tw, format, headers...)

	switch list := results.(type) {
	case []http.Group:
		for _, v := range list {
			row := []interface{}{v.ID, v.Name, v.CreatedAt, v.ModifiedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Project:
		for _, v := range list {
			row := []interface{}{v.ID, v.Name, v.CreatedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Template:
		for _, v := range list {
			row := []interface{}{v.ID, v.GroupID, v.CreatedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Transport:
		for _, v := range list {
			row := []interface{}{v.ID, fmt.Sprintf("%s:%d", v.Host, v.Port), v.IsActive, v.Username, v.CreatedAt}
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
