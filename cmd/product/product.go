package product

import (
	"fmt"

	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "product",
		Short: "TikTok Shop product intelligence",
		Args:  rejectArgs,
		RunE:  showHelp,
	}
	cmd.AddCommand(newBasic())
	cmd.AddCommand(newAnalytics())
	cmd.AddCommand(newLeaderboard())
	return cmd
}

func newBasic() *cobra.Command {
	cmd := &cobra.Command{Use: "basic", Short: "Product discovery and basic information", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newBasicSearch())
	cmd.AddCommand(newBasicList())
	cmd.AddCommand(newBasicDetail())
	return cmd
}

func newBasicSearch() *cobra.Command {
	var keyword, region, sortType, priceRange, live, cod, count, offset string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search realtime TikTok Shop products",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("keyword", keyword, "example: echotik product basic search --keyword sunscreen --region US")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("region", region, "example: echotik product basic search --keyword sunscreen --region US")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "product",
				Capability: "basic",
				Operation:  "search",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/product/search",
				Params: schema.Clean(map[string]string{
					"sk":          keyword,
					"region":      region,
					"sort_type":   sortType,
					"price_range": priceRange,
					"live":        live,
					"cod":         cod,
					"count":       count,
					"offset":      offset,
				}),
			})
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

func newBasicList() *cobra.Command {
	var region, categoryID, categoryL2ID, categoryL3ID, page, pageSize, sortField, sortType string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List products from EchoTik offline product library",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik product basic list --region US --page-size 20")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "product",
				Capability: "basic",
				Operation:  "list",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/product/list",
				Params: schema.Clean(map[string]string{
					"region":             region,
					"category_id":        categoryID,
					"category_l2_id":     categoryL2ID,
					"category_l3_id":     categoryL3ID,
					"product_sort_field": sortField,
					"sort_type":          sortType,
					"page_num":           page,
					"page_size":          pageSize,
				}),
			})
		},
	}
	bindRegionCategoryPage(cmd, &region, &categoryID, &categoryL2ID, &categoryL3ID, &page, &pageSize)
	cmd.Flags().StringVar(&sortField, "sort-field", "", "EchoTik product_sort_field")
	cmd.Flags().StringVar(&sortType, "sort-type", "", "sort order: 0=asc, 1=desc")
	return cmd
}

func newBasicDetail() *cobra.Command {
	var productID string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Fetch product details by product ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("product-id", productID, "example: echotik product basic detail --product-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "product",
				Capability: "basic",
				Operation:  "detail",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/product/detail",
				Params: map[string]string{
					"product_ids": productID,
				},
			})
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "product ID; use comma-separated IDs for batch detail, max 10")
	return cmd
}

func newAnalytics() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Product relation analysis", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newProductAnalyticsList("creators", "creator", "product_promoted_by_creator", "/api/v3/echotik/product/influencer/list", "product_influencer_sort_field"))
	cmd.AddCommand(newProductAnalyticsList("videos", "video", "product_promoted_by_video", "/api/v3/echotik/product/video/list", "product_video_sort_field"))
	cmd.AddCommand(newProductAnalyticsList("lives", "live", "product_promoted_in_live", "/api/v3/echotik/product/live/list", "product_live_sort_field"))
	cmd.AddCommand(newProductTrends())
	return cmd
}

func newProductAnalyticsList(operation, recordsAs, relationType, path, sortFieldName string) *cobra.Command {
	var productID, page, pageSize, sortField, sortType string
	cmd := &cobra.Command{
		Use:   operation,
		Short: "Fetch product related " + operation,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("product-id", productID, "example: echotik product analytics "+operation+" --product-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:       "product",
				Capability:   "analytics",
				Operation:    operation,
				Freshness:    "offline_t_plus_1",
				Path:         path,
				RecordsAs:    recordsAs,
				RelationFrom: &normalize.RelationNode{Type: "product", ID: productID},
				RelationType: relationType,
				Params: schema.Clean(map[string]string{
					"product_id":  productID,
					sortFieldName: sortField,
					"sort_type":   sortType,
					"page_num":    page,
					"page_size":   pageSize,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "product ID")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "EchoTik sort field for this relation")
	cmd.Flags().StringVar(&sortType, "sort-type", "1", "sort order: 0=asc, 1=desc")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newProductTrends() *cobra.Command {
	var productID, startDate, endDate, page, pageSize string
	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Fetch product trend snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("product-id", productID, "example: echotik product analytics trends --product-id 123 --start-date 2026-01-01 --end-date 2026-01-31")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("start-date", startDate, "use yyyy-MM-dd")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("end-date", endDate, "use yyyy-MM-dd")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "product",
				Capability: "analytics",
				Operation:  "trends",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/product/trend",
				Params: map[string]string{
					"product_id": productID,
					"start_date": startDate,
					"end_date":   endDate,
					"page_num":   page,
					"page_size":  pageSize,
				},
			})
		},
	}
	cmd.Flags().StringVar(&productID, "product-id", "", "product ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "start date, yyyy-MM-dd")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date, yyyy-MM-dd")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newLeaderboard() *cobra.Command {
	cmd := &cobra.Command{Use: "leaderboard", Short: "Product ranking and trend discovery", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newLeaderboardTop())
	return cmd
}

func newLeaderboardTop() *cobra.Command {
	var region, categoryID, categoryL2ID, categoryL3ID, date, rankType, rankField, page, pageSize string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Fetch product ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik product leaderboard top --region US --date 2026-01-01")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("date", date, "rank date in yyyy-MM-dd format")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "product",
				Capability: "leaderboard",
				Operation:  "top",
				Freshness:  "ranking",
				Path:       "/api/v3/echotik/product/ranklist",
				Params: schema.Clean(map[string]string{
					"region":             region,
					"category_id":        categoryID,
					"category_l2_id":     categoryL2ID,
					"category_l3_id":     categoryL3ID,
					"date":               date,
					"rank_type":          schema.RankType(rankType),
					"product_rank_field": rankField,
					"page_num":           page,
					"page_size":          pageSize,
				}),
			})
		},
	}
	bindRegionCategoryPage(cmd, &region, &categoryID, &categoryL2ID, &categoryL3ID, &page, &pageSize)
	cmd.Flags().StringVar(&date, "date", "", "rank date in yyyy-MM-dd format")
	cmd.Flags().StringVar(&rankType, "type", "daily", "ranking type: daily|weekly|monthly or 1|2|3")
	cmd.Flags().StringVar(&rankField, "rank-field", "1", "ranking field: 1=total_sale_cnt, 2=total_ifl_cnt")
	return cmd
}

func bindRegionCategoryPage(cmd *cobra.Command, region, categoryID, categoryL2ID, categoryL3ID, page, pageSize *string) {
	cmd.Flags().StringVar(region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(categoryID, "category-id", "", "EchoTik category ID")
	cmd.Flags().StringVar(categoryL2ID, "category-l2-id", "", "EchoTik level-2 category ID")
	cmd.Flags().StringVar(categoryL3ID, "category-l3-id", "", "EchoTik level-3 category ID")
	cmd.Flags().StringVar(page, "page", "1", "page number")
	cmd.Flags().StringVar(pageSize, "page-size", "20", "page size")
}

func rejectArgs(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return fmt.Errorf("unknown command %q for %q", args[0], cmd.CommandPath())
	}
	return nil
}

func showHelp(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}
