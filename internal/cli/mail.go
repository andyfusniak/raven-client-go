package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/andyfusniak/raven-client-go/http"
	"github.com/spf13/cobra"
)

// NewCmdListMail list mail sub command.
func NewCmdListMail() *cobra.Command {
	return &cobra.Command{
		Use:     "mail",
		Short:   "List mail",
		Aliases: []string{"mails"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			results, err := app.HTTPClient.ListMail(ctx, app.projectID)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%s\t%s\t%v\t%v\n"
			headers := []interface{}{"MAIL ID", "TEMPLATE ID", "STATUS", "TO", "CREATED AT", "SENT AT"}
			if err := renderTable(os.Stdout, results, format, headers, time.Time{}); err != nil {
				return fmt.Errorf("list mail failed to render table: %+v", err)
			}
			return nil
		},
	}
}

// NewCmdListMailLogs list mlogs sub command.
func NewCmdListMailLogs() *cobra.Command {
	return &cobra.Command{
		Use:     "mail-logs MAIL_ID",
		Short:   "List mail logs",
		Aliases: []string{"maillogs", "mlogs"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing MAIL_ID argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			mailID := args[0]

			mail, err := app.HTTPClient.GetMail(ctx, app.projectID, mailID)
			if err != nil {
				return err
			}

			mails := []http.Mail{
				*mail,
			}
			format := "%s\t%s\t%s\t%s\t%v\t%v\n"
			headers := []interface{}{"MAIL ID", "TEMPLATE ID", "STATUS", "TO", "CREATED AT", "SENT AT"}
			if err := renderTable(os.Stdout, mails, format, headers, time.Time{}); err != nil {
				return fmt.Errorf("list mail failed to render table: %+v", err)
			}

			results, err := app.HTTPClient.ListMailLogs(ctx, app.projectID, mailID)
			if err != nil {
				return err
			}

			format = "%s\t%s\t\t%v\t%s\t%v\n"
			headers = []interface{}{
				"├── MAIL LOG ID", "STATUS",
				"SMTP CODE", "MSG", "DURATION",
			}
			if err := renderTable(os.Stdout, results, format, headers, mail.CreatedAt); err != nil {
				return fmt.Errorf("list mail logs failed to render table: %+v", err)
			}
			return nil
		},
	}
}

// NewCmdGetMail get mail sub command.
func NewCmdGetMail() *cobra.Command {
	return &cobra.Command{
		Use:     "mail MAIL_ID",
		Short:   "Get a mail entry",
		Aliases: []string{"mails"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing MAIL_ID argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			mailID := args[0]
			result, err := app.HTTPClient.GetMail(ctx, app.projectID, mailID)
			if err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeMailNotFound {
						fmt.Fprintf(os.Stderr,
							"Mail %q not found - use raven list mail for a full list.\n",
							mailID)
						os.Exit(1)
					}
				}
				return err
			}

			if err := renderMail(os.Stdout, result); err != nil {
				return err
			}

			return nil
		},
	}
}

func renderMail(w io.Writer, m *http.Mail) error {
	fmt.Fprintf(w, "MAIL ID:\t%s\n", m.ID)
	fmt.Fprintf(w, "TEMPLATE ID:\t%s\n", m.TemplateID)
	fmt.Fprintf(w, "PROJECT ID:\t%s\n", m.ProjectID)
	fmt.Fprintf(w, "STATUS:\t\t%s\n", m.Status)
	fmt.Fprintf(w, "TO:\t\t%s\n", m.EmailTo)
	fmt.Fprintf(w, "FROM:\t\t%s\n", m.EmailFrom)
	fmt.Fprintf(w, "REPLY TO:\t%s\n", m.EmailReplyTo)
	fmt.Fprintf(w, "SUBJECT:\t%s\n", m.Subject)
	// fmt.Fprintf(w, "TEXT: %s\n", m.TxtDigest[:7])
	// fmt.Fprintf(w, "%s\n", t.Txt)

	fmt.Fprintf(w, "CREATED AT:\t%s\n", m.CreatedAt)
	fmt.Fprintf(w, "SENT AT:\t%s\n", renderMailSentAt(m.SentAt))
	fmt.Fprintf(w, "LAST MODIFIED:\t%v\n", m.ModifiedAt)
	fmt.Fprintln(w)

	return nil
}
