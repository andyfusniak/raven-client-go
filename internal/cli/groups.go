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

// NewCmdCreateGroup create group sub command.
func NewCmdCreateGroup() *cobra.Command {
	return &cobra.Command{
		Use:     "group NAME",
		Short:   "Create group",
		Aliases: []string{"groups"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			name := args[0]
			result, err := app.HTTPClient.CreateGroup(ctx, app.projectID, name)
			if err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeProjectNotFound {
						os.Exit(1)
					}
					fmt.Fprintf(os.Stderr, "untrapped api error: %#v\n", terr)
					os.Exit(1)
				}

				fmt.Fprintf(os.Stderr, "unknown error: %+v\n", err)
				os.Exit(1)
			}

			renderGroup(os.Stderr, result)

			return nil
		},
	}
}

// NewCmdGetGroup get a group sub command.
func NewCmdGetGroup() *cobra.Command {
	return &cobra.Command{
		Use:     "group GROUP_ID",
		Short:   "Get a group",
		Aliases: []string{"groups"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing GROUP_ID argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			groupID := args[0]
			result, err := app.HTTPClient.GetGroup(ctx, app.projectID, groupID)
			if err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeGroupNotFound {
						fmt.Fprintf(os.Stderr,
							"group %q not found - use raven list groups for a full list.\n",
							groupID)
						os.Exit(1)
					}
				}
				return err
			}

			if err := renderGroup(os.Stdout, result); err != nil {
				return err
			}

			return nil
		},
	}
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

			results, err := app.HTTPClient.ListGroups(ctx, app.projectID)
			if err != nil {
				return err
			}

			format := "%s\t%s\t%s\t%s\n"
			headers := []interface{}{"GROUP ID", "NAME", "CREATED", "LAST MODIFIED"}
			if err := renderTable(os.Stdout, results, format, headers, time.Time{}); err != nil {
				return fmt.Errorf("list groups failed to render table: %+v", err)
			}

			return nil
		},
	}
}

// NewCmdDeleteGroup delete group sub command.
func NewCmdDeleteGroup() *cobra.Command {
	return &cobra.Command{
		Use:     "group GROUP_ID",
		Short:   "Delete a group",
		Aliases: []string{"groups"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing GROUP_ID argument")
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			app := ctx.Value(AppKey("app")).(*App)

			groupID := args[0]
			if err := app.HTTPClient.DeleteGroup(ctx, app.projectID, groupID); err != nil {
				if terr, ok := err.(*http.APIError); ok {
					if terr.Code == http.ErrCodeGroupNotFound {
						fmt.Fprintf(os.Stderr, "group %s not found\n", groupID)
						os.Exit(1)
					}
					if terr.Code == http.ErrCodeGroupIDInvalid {
						fmt.Fprintf(os.Stderr, "group id is an invalid format\n")
						os.Exit(1)
					}
					fmt.Printf("%#v\n", terr)
					os.Exit(1)
				}

				fmt.Printf("unknown error: %+v\n", err)
				os.Exit(1)
			}

			return nil
		},
	}
}

func renderGroup(w io.Writer, g *http.Group) error {
	fmt.Fprintf(w, "GROUP ID:\t%s\n", g.ID)
	fmt.Fprintf(w, "PROJECT ID:\t%s\n", g.ProjectID)
	fmt.Fprintf(w, "NAME:\t\t%s\n", g.Name)
	fmt.Fprintf(w, "CREATED:\t%s\n", g.CreatedAt)
	fmt.Fprintf(w, "LAST MODIFIED:\t%v\n", g.ModifiedAt)
	return nil
}
