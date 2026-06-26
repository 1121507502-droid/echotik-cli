package shop

import (
	"fmt"

	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "shop",
		Short: "TikTok Shop seller intelligence",
		Args:  rejectArgs,
		RunE:  showHelp,
	}
	cmd.AddCommand(newBasic())
	cmd.AddCommand(newAnalytics())
	cmd.AddCommand(newLeaderboard())
	return cmd
}

func newBasic() *cobra.Command {
	cmd := &cobra.Command{Use: "basic", Short: "Shop discovery and basic information", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newSearch())
	cmd.AddCommand(newList())
	cmd.AddCommand(newDetail())
	return cmd
}

func newSearch() *cobra.Command {
	var keyword, region, size, searchType, sortType string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search shops through EchoTik general search",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("keyword", keyword, "example: echotik shop basic search --keyword skincare --region US")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "shop",
				Capability: "basic",
				Operation:  "search",
				Freshness:  "search",
				Path:       "/api/v3/echotik/search/items",
				Params: schema.Clean(map[string]string{
					"sk":         keyword,
					"region":     region,
					"type":       "3",
					"size":       size,
					"searchType": searchType,
					"sortType":   sortType,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword")
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&size, "size", "10", "result size, max 30")
	cmd.Flags().StringVar(&searchType, "search-type", "", "0=fuzzy, 1=precise")
	cmd.Flags().StringVar(&sortType, "sort-type", "", "EchoTik general search sortType")
	return cmd
}

func newList() *cobra.Command {
	var region, categoryID, categoryL2ID, categoryL3ID, page, pageSize, sortField, sortType string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List shops from EchoTik offline seller library",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik shop basic list --region US --page-size 20")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "shop",
				Capability: "basic",
				Operation:  "list",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/seller/list",
				Params: schema.Clean(map[string]string{
					"region":            region,
					"category_id":       categoryID,
					"category_l2_id":    categoryL2ID,
					"category_l3_id":    categoryL3ID,
					"seller_sort_field": sortField,
					"sort_type":         sortType,
					"page_num":          page,
					"page_size":         pageSize,
				}),
			})
		},
	}
	bindRegionCategoryPage(cmd, &region, &categoryID, &categoryL2ID, &categoryL3ID, &page, &pageSize)
	cmd.Flags().StringVar(&sortField, "sort-field", "", "EchoTik seller_sort_field")
	cmd.Flags().StringVar(&sortType, "sort-type", "", "sort order: 0=asc, 1=desc")
	return cmd
}

func newDetail() *cobra.Command {
	var sellerID string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Fetch shop details by seller ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("seller-id", sellerID, "example: echotik shop basic detail --seller-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "shop",
				Capability: "basic",
				Operation:  "detail",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/seller/detail",
				Params: map[string]string{
					"seller_id": sellerID,
				},
			})
		},
	}
	cmd.Flags().StringVar(&sellerID, "seller-id", "", "seller/shop ID")
	return cmd
}

func newAnalytics() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Shop relation analysis", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newShopAnalyticsList("products", "product", "shop_sells_product", "/api/v3/echotik/seller/product/list", "seller_product_sort_field"))
	cmd.AddCommand(newShopAnalyticsList("creators", "creator", "shop_promoted_by_creator", "/api/v3/echotik/seller/influencer/list", "seller_influencer_sort_field"))
	cmd.AddCommand(newShopAnalyticsList("videos", "video", "shop_promoted_by_video", "/api/v3/echotik/seller/video/list", "seller_video_sort_field"))
	cmd.AddCommand(newShopAnalyticsList("lives", "live", "shop_promoted_in_live", "/api/v3/echotik/seller/live/list", "seller_live_sort_field"))
	cmd.AddCommand(newTrends())
	return cmd
}

func newShopAnalyticsList(operation, recordsAs, relationType, path, sortFieldName string) *cobra.Command {
	var sellerID, page, pageSize, sortField, sortType string
	cmd := &cobra.Command{
		Use:   operation,
		Short: "Fetch shop related " + operation,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("seller-id", sellerID, "example: echotik shop analytics "+operation+" --seller-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:       "shop",
				Capability:   "analytics",
				Operation:    operation,
				Freshness:    "offline_t_plus_1",
				Path:         path,
				RecordsAs:    recordsAs,
				RelationFrom: &normalize.RelationNode{Type: "shop", ID: sellerID},
				RelationType: relationType,
				Params: schema.Clean(map[string]string{
					"seller_id":   sellerID,
					sortFieldName: sortField,
					"sort_type":   sortType,
					"page_num":    page,
					"page_size":   pageSize,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&sellerID, "seller-id", "", "seller/shop ID")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "EchoTik sort field for this relation")
	cmd.Flags().StringVar(&sortType, "sort-type", "1", "sort order: 0=asc, 1=desc")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newTrends() *cobra.Command {
	var sellerID, startDate, endDate, page, pageSize string
	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Fetch shop trend snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("seller-id", sellerID, "example: echotik shop analytics trends --seller-id 123 --start-date 2026-01-01 --end-date 2026-01-31")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("start-date", startDate, "use yyyy-MM-dd")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("end-date", endDate, "use yyyy-MM-dd")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "shop",
				Capability: "analytics",
				Operation:  "trends",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/seller/trend",
				Params: map[string]string{
					"seller_id":  sellerID,
					"start_date": startDate,
					"end_date":   endDate,
					"page_num":   page,
					"page_size":  pageSize,
				},
			})
		},
	}
	cmd.Flags().StringVar(&sellerID, "seller-id", "", "seller/shop ID")
	cmd.Flags().StringVar(&startDate, "start-date", "", "start date, yyyy-MM-dd")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date, yyyy-MM-dd")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newLeaderboard() *cobra.Command {
	cmd := &cobra.Command{Use: "leaderboard", Short: "Shop ranking and trend discovery", Args: rejectArgs, RunE: showHelp}
	cmd.AddCommand(newLeaderboardTop())
	return cmd
}

func newLeaderboardTop() *cobra.Command {
	var region, categoryID, categoryL2ID, categoryL3ID, date, rankType, rankField, page, pageSize string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Fetch shop ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik shop leaderboard top --region US --date 2026-01-01")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("date", date, "rank date in yyyy-MM-dd format")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "shop",
				Capability: "leaderboard",
				Operation:  "top",
				Freshness:  "ranking",
				Path:       "/api/v3/echotik/seller/ranklist",
				Params: schema.Clean(map[string]string{
					"region":            region,
					"category_id":       categoryID,
					"category_l2_id":    categoryL2ID,
					"category_l3_id":    categoryL3ID,
					"date":              date,
					"rank_type":         schema.RankType(rankType),
					"seller_rank_field": rankField,
					"page_num":          page,
					"page_size":         pageSize,
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
