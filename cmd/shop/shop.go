package shop

import (
	"github.com/echotik/cli/internal/client"
	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shop",
		Short: "TikTok Shop seller intelligence",
	}
	cmd.AddCommand(newList())
	cmd.AddCommand(newRank())
	return cmd
}

func newList() *cobra.Command {
	var region, categoryID, keyword, page, pageSize string
	cmd := &cobra.Command{
		Use:   "+list",
		Short: "List shops from EchoTik offline seller library",
		RunE: func(cmd *cobra.Command, args []string) error {
			return do(cmd, "/api/v3/echotik/seller/list", clean(map[string]string{
				"region":      region,
				"category_id": categoryID,
				"keyword":     keyword,
				"page":        page,
				"page_size":   pageSize,
			}), "offline_t_plus_1")
		},
	}
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "EchoTik category ID")
	cmd.Flags().StringVar(&keyword, "keyword", "", "shop keyword")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "20", "page size")
	return cmd
}

func newRank() *cobra.Command {
	var region, categoryID, date, rankType, page, pageSize string
	cmd := &cobra.Command{
		Use:   "+rank",
		Short: "Fetch shop ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return do(cmd, "/api/v3/echotik/seller/ranklist", clean(map[string]string{
				"region":      region,
				"category_id": categoryID,
				"date":        date,
				"type":        rankType,
				"page":        page,
				"page_size":   pageSize,
			}), "ranking")
		},
	}
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "EchoTik category ID")
	cmd.Flags().StringVar(&date, "date", "", "rank date, format depends on EchoTik API")
	cmd.Flags().StringVar(&rankType, "type", "daily", "ranking type: daily|weekly|monthly")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "20", "page size")
	return cmd
}

func clean(params map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range params {
		if v != "" {
			out[k] = v
		}
	}
	return out
}

func do(cmd *cobra.Command, path string, params map[string]string, freshness string) error {
	baseURL, username, password, err := core.ResolveCredential()
	if err != nil {
		return output.NewError("authentication_error", err.Error(), "run: echotik config set-credential")
	}
	c := client.New(baseURL, username, password)
	resp, err := c.Do(cmd.Context(), client.Request{Method: "GET", Path: path, Params: params})
	if err != nil {
		return output.NewError("api_error", err.Error(), "check parameters or retry with backoff")
	}
	data := resp.JSON
	if data == nil {
		data = string(resp.Raw)
	}
	return output.Success(cmd.OutOrStdout(), data, map[string]interface{}{
		"path":       path,
		"freshness":  freshness,
		"statusCode": resp.StatusCode,
	})
}
