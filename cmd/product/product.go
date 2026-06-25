package product

import (
	"github.com/echotik/cli/internal/client"
	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "TikTok Shop product intelligence",
	}
	cmd.AddCommand(newList())
	cmd.AddCommand(newSearch())
	cmd.AddCommand(newRank())
	return cmd
}

type commonFlags struct {
	Region     string
	Keyword    string
	CategoryID string
	Page       string
	PageSize   string
}

func bindCommon(cmd *cobra.Command, f *commonFlags) {
	cmd.Flags().StringVar(&f.Region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&f.Keyword, "keyword", "", "search keyword")
	cmd.Flags().StringVar(&f.CategoryID, "category-id", "", "EchoTik category ID")
	cmd.Flags().StringVar(&f.Page, "page", "1", "page number")
	cmd.Flags().StringVar(&f.PageSize, "page-size", "20", "page size")
}

func newList() *cobra.Command {
	var f commonFlags
	cmd := &cobra.Command{
		Use:   "+list",
		Short: "List products from EchoTik offline product library",
		RunE: func(cmd *cobra.Command, args []string) error {
			params := baseParams(f)
			return do(cmd, "/api/v3/echotik/product/list", params, "offline_t_plus_1")
		},
	}
	bindCommon(cmd, &f)
	return cmd
}

func newSearch() *cobra.Command {
	var region, keyword, sortType, priceRange, live, cod, count, offset string
	cmd := &cobra.Command{
		Use:   "+search",
		Short: "Search TikTok Shop products in realtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			if keyword == "" {
				return output.NewError("validation_error", "--keyword is required", "example: echotik product +search --keyword sunscreen --region US")
			}
			params := clean(map[string]string{
				"sk":          keyword,
				"region":      region,
				"sort_type":   sortType,
				"price_range": priceRange,
				"live":        live,
				"cod":         cod,
				"count":       count,
				"offset":      offset,
			})
			return do(cmd, "/api/v3/realtime/product/search", params, "realtime")
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword; sent to EchoTik as sk")
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&sortType, "sort-type", "", "sort type: 1=PRICE_ASC, 2=PRICE_DESC, 3=BEST_SELLERS, 4=RELEVANCE")
	cmd.Flags().StringVar(&priceRange, "price-range", "", "price range filter, e.g. 0,100")
	cmd.Flags().StringVar(&live, "live", "", "display live broadcast products: 0=no, 1=yes")
	cmd.Flags().StringVar(&cod, "cod", "", "display cash-on-delivery products: 0=no, 1=yes")
	cmd.Flags().StringVar(&count, "count", "10", "number of results")
	cmd.Flags().StringVar(&offset, "offset", "0", "result offset")
	return cmd
}

func newRank() *cobra.Command {
	var region, categoryID, categoryL2ID, categoryL3ID, date, rankType, rankField, page, pageSize string
	cmd := &cobra.Command{
		Use:   "+rank",
		Short: "Fetch product ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			apiRankType := mapRankType(rankType)
			params := map[string]string{
				"region":             region,
				"category_id":        categoryID,
				"category_l2_id":     categoryL2ID,
				"category_l3_id":     categoryL3ID,
				"date":               date,
				"rank_type":          apiRankType,
				"product_rank_field": rankField,
				"page_num":           page,
				"page_size":          pageSize,
			}
			return do(cmd, "/api/v3/echotik/product/ranklist", clean(params), "ranking")
		},
	}
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "EchoTik category ID")
	cmd.Flags().StringVar(&categoryL2ID, "category-l2-id", "", "EchoTik level-2 category ID")
	cmd.Flags().StringVar(&categoryL3ID, "category-l3-id", "", "EchoTik level-3 category ID")
	cmd.Flags().StringVar(&date, "date", "", "rank date in yyyy-MM-dd format")
	cmd.Flags().StringVar(&rankType, "type", "daily", "ranking type: daily|weekly|monthly or 1|2|3")
	cmd.Flags().StringVar(&rankField, "rank-field", "1", "ranking field: 1=total_sale_cnt, 2=total_ifl_cnt")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "20", "page size")
	return cmd
}

func baseParams(f commonFlags) map[string]string {
	return clean(map[string]string{
		"region":      f.Region,
		"keyword":     f.Keyword,
		"category_id": f.CategoryID,
		"page_num":    f.Page,
		"page_size":   f.PageSize,
	})
}

func mapRankType(raw string) string {
	switch raw {
	case "daily":
		return "1"
	case "weekly":
		return "2"
	case "monthly":
		return "3"
	default:
		return raw
	}
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
