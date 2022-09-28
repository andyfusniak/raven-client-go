package cli

import (
	"errors"
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

// NewCmdGet get sub command.
func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get a single resource",
	}
	cmd.AddCommand(NewCmdGetTemplate())
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

// NewCmdListGroups list groups sub command.
func NewCmdListGroups() *cobra.Command {
	return &cobra.Command{
		Use:     "groups",
		Short:   "List groups",
		Aliases: []string{"group"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			results, err := app.HTTPClient.ListGroups(ctx)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%s\t%s\n"
			headers := []interface{}{"Group ID", "Name", "Created", "Last Modified"}
			if err := renderTable(os.Stdout, results, format, headers); err != nil {
				return fmt.Errorf("list groups failed to render table: %+v", err)
			}

			return nil
		},
	}
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

// NewCmdListTemplates list templates sub command.
func NewCmdListTemplates() *cobra.Command {
	return &cobra.Command{
		Use:     "templates",
		Short:   "List templates",
		Aliases: []string{"template"},
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

// NewCmdGetTemplate get template sub command.
func NewCmdGetTemplate() *cobra.Command {
	return &cobra.Command{
		Use:     "template TEMPLATE_ID",
		Short:   "Get a single template",
		Aliases: []string{"templates"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing TEMPLATE_ID argument")
			}
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			templateID := args[0]
			result, err := app.HTTPClient.GetTemplate(ctx, templateID)
			if err != nil {
				if t, ok := err.(*http.APIError); ok {
					if t.Code == http.ErrCodeTemplateNotFound {
						fmt.Fprintf(os.Stderr,
							"Template %q not found - use raven list templates for a full list.\n",
							templateID)
						os.Exit(1)
					}
				}
				return err
			}

			if err := renderTemplate(os.Stdout, result); err != nil {
				return err
			}

			return nil
		},
	}
}

const checkMark = "✔"
const crossMark = "✘"

const heavyGreenCheckMark = "✅"
const heavyRedCrossMark = "❌"

func renderTemplate(w io.Writer, t *http.Template) error {
	fmt.Fprintf(w, "Template ID:\t\t%s %s\n", t.ID, isTemplateReady(t.TxtTemplateCompiledOK, t.HTMLTemplateCompiledOK))
	fmt.Fprintf(w, "Project ID:\t\t%s\n", t.ProjectID)
	fmt.Fprintf(w, "Group ID:\t\t%s\n", t.GroupID)
	fmt.Fprintf(w, "Last modified:\t\t%v\n", t.ModifiedAt)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "TEXT: %s %s\n", t.TxtDigest[:7], checkOrCrossWithStatus(t.TxtTemplateCompiledOK, true))
	fmt.Fprintf(w, "%s\n", t.Txt)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "HTML: %s %s\n", t.HTMLDigest[:7], checkOrCrossWithStatus(t.HTMLTemplateCompiledOK, true))
	fmt.Fprintf(w, "%s\n", t.HTML)

	return nil
}

func checkOrCrossWithStatus(v bool, useHeavy bool) string {
	if useHeavy {
		if v {
			return heavyGreenCheckMark + " Ready"
		}
		return heavyRedCrossMark + " Failed to compile"
	}

	if v {
		return checkMark + " Ready"
	}
	return crossMark + " Failed to compile"
}

func isTemplateReady(txt, html bool) string {
	if txt && html {
		return heavyGreenCheckMark + " Ready"
	}
	return heavyRedCrossMark + " Failed to compile"
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
