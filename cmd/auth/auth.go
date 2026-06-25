package auth

import (
	"errors"

	"github.com/echotik/cli/internal/client"
	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "auth",
		Short: "Check EchoTik authentication",
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "status",
		Short: "Validate configured EchoTik credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			baseURL, username, password, err := core.ResolveCredential()
			if err != nil {
				return output.NewError("authentication_error", err.Error(), "run: echotik config set-credential")
			}
			c := client.New(baseURL, username, password)
			resp, err := c.Do(cmd.Context(), client.Request{
				Method: "GET",
				Path:   "/api/v3/echotik/product/list",
				Params: map[string]string{"region": "US", "page_num": "1", "page_size": "1"},
			})
			if err != nil {
				var httpErr *client.HTTPError
				if errors.As(err, &httpErr) && httpErr.StatusCode == 401 {
					return output.NewError("authentication_error", "invalid EchoTik credential", "run: echotik config set-credential")
				}
				return err
			}
			return output.Success(cmd.OutOrStdout(), map[string]any{
				"authenticated": true,
				"baseUrl":       baseURL,
				"statusCode":    resp.StatusCode,
			}, nil)
		},
	})
	return cmd
}
