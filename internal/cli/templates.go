package cli

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/andyfusniak/raven-client-go/http"
	"github.com/spf13/cobra"
)

// NewCmdCreateTemplate update a template
func NewCmdCreateTemplate() *cobra.Command {
	return &cobra.Command{
		Use:     "templates TEMPLATE_ID...",
		Short:   "Create a new template",
		Aliases: []string{"template"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("must contain at least one TEMPLATE_ID")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			for _, filename := range args {
				// fmt.Printf("Arg[%d] = %s\n", i, filename)

				if !fileExists(filename) {
					fmt.Fprintf(os.Stderr, "file %s does not exist\n", filename)
					os.Exit(1)
				}

				if !acceptedFileExtension(filename) {
					fmt.Fprint(os.Stderr, "only files with .txt or .html extensions are supported\n")
					os.Exit(1)
				}

				txt, err := os.ReadFile(filename)
				if err != nil {
					fmt.Fprintf(os.Stderr, "failed to read file %s\n", filename)
				}

				templateID := baseFilenameWithoutExt(filename)

				fmt.Printf("%s\n", templateID)
				result, err := app.HTTPClient.CreateTemplate(ctx, &http.CreateTemplateParams{
					ID:        templateID,
					ProjectID: app.projectID,
					GroupID:   "wC6yNEg79ZQVFQ62y3PD",
					Txt:       string(txt),
				})
				if err != nil {
					if terr, ok := err.(*http.APIError); ok {
						if terr.Code == http.ErrCodeTemplateExists {
							fmt.Fprintf(os.Stderr, "template %s already exists\n", templateID)
							os.Exit(1)
						}
						if terr.Code == http.ErrCodeGroupNotFound {
							fmt.Fprintf(os.Stderr,
								"target group could not be found. Use raven list groups.")
							os.Exit(1)
						}
						fmt.Fprintf(os.Stderr, "%#v\n", terr)
						os.Exit(1)
					}

					fmt.Fprintf(os.Stderr, "unknown error: %+v", err)
				}

				fmt.Printf("%#v\n", result)
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

			results, err := app.HTTPClient.ListTemplates(ctx, app.projectID)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%s\n"
			headers := []interface{}{"TEMPLATE ID", "GROUP ID", "CREATED"}
			if err := renderTable(os.Stdout, results, format, headers, time.Time{}); err != nil {
				return fmt.Errorf("list templates failed to render table: %+v", err)
			}
			return nil
		},
	}
}

// NewCmdGetTemplate get template sub command.
func NewCmdGetTemplate() *cobra.Command {
	return &cobra.Command{
		Use:     "template TEMPLATE_ID",
		Short:   "Get a template",
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
			result, err := app.HTTPClient.GetTemplate(ctx, app.projectID, templateID)
			if err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeTemplateNotFound {
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

// NewCmdDeleteTemplate get template sub command.
func NewCmdDeleteTemplate() *cobra.Command {
	return &cobra.Command{
		Use:     "template TEMPLATE_ID",
		Short:   "Delete a template",
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
			if err := app.HTTPClient.DeleteTemplate(ctx, app.projectID, templateID); err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeTemplateNotFound {
						fmt.Fprintf(os.Stderr, "template %s not found\n", templateID)
						os.Exit(1)
					}

					fmt.Printf("%#v\n", terr)
				}

				fmt.Printf("unknown error: %+v\n", err)
			}

			return nil
		},
	}
}

func fileExists(f string) bool {
	if _, err := os.Stat(f); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func acceptedFileExtension(f string) bool {
	ext := filepath.Ext(f)
	return ext == ".txt" || ext == ".html"
}

func baseFilenameWithoutExt(filename string) string {
	return strings.TrimSuffix(filepath.Base(filename), filepath.Ext(filename))
}

func renderTemplate(w io.Writer, t *http.Template) error {
	fmt.Fprintf(w, "Template ID:\t\t%s\n", t.ID)
	fmt.Fprintf(w, "Project ID:\t\t%s\n", t.ProjectID)
	fmt.Fprintf(w, "Group ID:\t\t%s\n", t.GroupID)
	fmt.Fprintf(w, "Last modified:\t\t%v\n", t.ModifiedAt)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "TEXT: %s\n", t.TxtDigest[:7])
	fmt.Fprintf(w, "%s\n", t.Txt)
	fmt.Fprintln(w)

	fmt.Fprintf(w, "HTML: %s\n", t.HTMLDigest[:7])
	fmt.Fprintf(w, "%s\n", t.HTML)

	return nil
}
