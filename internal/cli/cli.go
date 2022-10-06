package cli

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/andyfusniak/raven-client-go/http"
	"github.com/spf13/cobra"
)

const checkMark = "✔"
const crossMark = "✘"

const heavyGreenCheckMark = "✅"
const heavyRedCrossMark = "❌"

const redCircle = "🔴"
const redSquare = "🟥"

const greenCircle = "🟢"
const greenSquare = "🟩"

const blueCircle = "🔵"
const blueSquare = "🟦"

const yellowCircle = "🟡"
const yellowSquare = "🟨"

const rewind = "↻"
const publishedEmoji = "📝"
const envelopeEmoji = "🖃"

// AppKey context key
type AppKey string

// App cli tool.
type App struct {
	version    string
	endpoint   string
	gitCommit  string
	projectID  string
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
		projectID:  "the-cloud-company",
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
	cmd.AddCommand(NewCmdListProjects())
	cmd.AddCommand(NewCmdListTransports())
	cmd.AddCommand(NewCmdListGroups())
	cmd.AddCommand(NewCmdListTemplates())
	cmd.AddCommand(NewCmdListMail())
	cmd.AddCommand(NewCmdListMailLogs())
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
	cmd.AddCommand(NewCmdGetMail())
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
			headers := []interface{}{"PROJECT ID", "NAME", "CREATED"}
			if err := renderTable(os.Stdout, results, format, headers, time.Time{}); err != nil {
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

			results, err := app.HTTPClient.ListTransports(ctx, app.projectID)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%v\t%s\t%v\n"
			headers := []interface{}{
				"TRANSPORT ID",
				"HOST:PORT",
				"ACTIVE",
				"USERNAME",
				"CREATED",
			}
			if err := renderTable(os.Stdout, results, format, headers, time.Time{}); err != nil {
				return fmt.Errorf("list transports failed to render table: %+v", err)
			}
			return nil
		},
	}
}

func renderTable(w io.Writer, results interface{}, format string, headers []interface{}, rel time.Time) error {
	tw := new(tabwriter.Writer).Init(w, 0, 8, 2, ' ', 0)

	fmt.Fprintf(tw, format, headers...)

	switch list := results.(type) {
	case []http.Project:
		for _, v := range list {
			row := []interface{}{v.ID, v.Name, v.CreatedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Transport:
		for _, v := range list {
			row := []interface{}{
				v.ID,
				fmt.Sprintf("%s:%d", v.Host, v.Port),
				renderTransportActive(v.Active),
				v.Username,
				v.CreatedAt,
			}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Group:
		for _, v := range list {
			row := []interface{}{v.ID, v.Name, v.CreatedAt, v.ModifiedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Template:
		for _, v := range list {
			row := []interface{}{v.ID, v.GroupID, v.CreatedAt}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.Mail:
		for _, v := range list {
			row := []interface{}{
				v.ID,
				v.TemplateID,
				renderMailStatus(v.Status),
				v.EmailTo,
				v.CreatedAt,
				renderMailSentAt(v.SentAt),
			}
			fmt.Fprintf(tw, format, row...)
		}
	case []http.MailLog:
		size := len(list)
		for i, v := range list {
			row := []interface{}{
				renderIndent(i, size) + " " + v.ID,
				renderMailStatus(v.Status),
				v.SMTPCode,
				v.Msg,
				renderRelativeTime(rel, v.CreatedAt),
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

func renderTransportActive(a bool) string {
	if a {
		return "✔"
	}
	return " "
}

func renderIndent(i, size int) string {
	if i >= size-1 {
		return "└──"
	}
	return "├──"
}

func renderRelativeTime(u, t time.Time) string {
	return fmt.Sprintf("%s", t.Sub(u))
}

func renderMailSentAt(t *time.Time) string {
	if t == nil {
		return "Not sent"
	}
	return t.Format(time.RFC1123)
}

func renderMailStatus(s string) string {
	switch s {
	case "pending":
		return "↓ " + strings.Title(s) // 💾, ⛁, ↓, ▼
	case "published":
		return "→ " + strings.Title(s) // 🖧, →
	case "received":
		return "← " + strings.Title(s) // or 🖃
	case "delivered":
		return "✔ " + strings.Title(s)
	}
	return s
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
