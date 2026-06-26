package creator

import (
	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{Use: "creator", Short: "TikTok Shop creator intelligence"}
	cmd.AddCommand(newBasic())
	cmd.AddCommand(newAnalytics())
	cmd.AddCommand(newLeaderboard())
	return cmd
}

func newBasic() *cobra.Command {
	cmd := &cobra.Command{Use: "basic", Short: "Creator discovery and basic information"}
	cmd.AddCommand(newSearch())
	cmd.AddCommand(newDetail())
	return cmd
}

func newSearch() *cobra.Command {
	var keyword, region, size, searchType, sortType string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search creators through EchoTik general search",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("keyword", keyword, "example: echotik creator basic search --keyword skincare --region US")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "creator",
				Capability: "basic",
				Operation:  "search",
				Freshness:  "search",
				Path:       "/api/v3/echotik/search/items",
				RecordsAs:  "creator",
				Params: schema.Clean(map[string]string{
					"sk":         keyword,
					"region":     region,
					"type":       "1",
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

func newDetail() *cobra.Command {
	var creatorID, uniqueID string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Fetch creator details by user ID or unique ID",
		RunE: func(cmd *cobra.Command, args []string) error {
			if creatorID == "" && uniqueID == "" {
				return runner.Validate(schema.Require("creator-id", creatorID, "example: echotik creator basic detail --creator-id 123 or --unique-id handle"))
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "creator",
				Capability: "basic",
				Operation:  "detail",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/influencer/detail",
				RecordsAs:  "creator",
				Params: schema.Clean(map[string]string{
					"user_ids":   creatorID,
					"unique_ids": uniqueID,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&creatorID, "creator-id", "", "creator user_id; comma-separated max 10")
	cmd.Flags().StringVar(&uniqueID, "unique-id", "", "creator unique_id; comma-separated max 10")
	return cmd
}

func newAnalytics() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Creator relation analysis"}
	cmd.AddCommand(newCreatorAnalyticsList("products", "product", "creator_promotes_product", "/api/v3/echotik/influencer/product/list", "influencer_product_sort_field"))
	cmd.AddCommand(newCreatorAnalyticsList("videos", "video", "creator_published_video", "/api/v3/echotik/influencer/video/list", "influencer_video_sort_field"))
	cmd.AddCommand(newCreatorLives())
	cmd.AddCommand(newTrends())
	return cmd
}

func newCreatorAnalyticsList(operation, recordsAs, relationType, path, sortFieldName string) *cobra.Command {
	var creatorID, uniqueID, productID, categoryID, page, pageSize, sortField, sortType string
	cmd := &cobra.Command{
		Use:   operation,
		Short: "Fetch creator related " + operation,
		RunE: func(cmd *cobra.Command, args []string) error {
			if creatorID == "" && uniqueID == "" {
				return runner.Validate(schema.Require("creator-id", creatorID, "example: echotik creator analytics "+operation+" --creator-id 123"))
			}
			fromID := creatorID
			if fromID == "" {
				fromID = uniqueID
			}
			return runner.Run(cmd, runner.Request{
				Entity:       "creator",
				Capability:   "analytics",
				Operation:    operation,
				Freshness:    "offline_t_plus_1",
				Path:         path,
				RecordsAs:    recordsAs,
				RelationFrom: &normalize.RelationNode{Type: "creator", ID: fromID},
				RelationType: relationType,
				Params: schema.Clean(map[string]string{
					"user_id":     creatorID,
					"unique_id":   uniqueID,
					"product_id":  productID,
					"category_id": categoryID,
					sortFieldName: sortField,
					"sort_type":   sortType,
					"page_num":    page,
					"page_size":   pageSize,
				}),
			})
		},
	}
	bindCreatorRelationFlags(cmd, &creatorID, &uniqueID, &page, &pageSize)
	cmd.Flags().StringVar(&productID, "product-id", "", "filter by product ID")
	cmd.Flags().StringVar(&categoryID, "category-id", "", "filter by product category ID")
	cmd.Flags().StringVar(&sortField, "sort-field", "", "EchoTik sort field for this relation")
	cmd.Flags().StringVar(&sortType, "sort-type", "1", "sort order: 0=asc, 1=desc")
	return cmd
}

func newCreatorLives() *cobra.Command {
	var creatorID, page, pageSize string
	cmd := &cobra.Command{
		Use:   "lives",
		Short: "Fetch creator live streams",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("creator-id", creatorID, "example: echotik creator analytics lives --creator-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:       "creator",
				Capability:   "analytics",
				Operation:    "lives",
				Freshness:    "offline_t_plus_1",
				Path:         "/api/v3/echotik/influencer/live/list",
				RecordsAs:    "live",
				RelationFrom: &normalize.RelationNode{Type: "creator", ID: creatorID},
				RelationType: "creator_hosted_live",
				Params: map[string]string{
					"user_id":   creatorID,
					"page_num":  page,
					"page_size": pageSize,
				},
			})
		},
	}
	cmd.Flags().StringVar(&creatorID, "creator-id", "", "creator user_id")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newTrends() *cobra.Command {
	var creatorID, startDate, endDate, page, pageSize string
	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Fetch creator trend snapshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, check := range []struct{ name, value, hint string }{
				{"creator-id", creatorID, "example: echotik creator analytics trends --creator-id 123 --start-date 2026-01-01 --end-date 2026-01-31"},
				{"start-date", startDate, "use yyyy-MM-dd"},
				{"end-date", endDate, "use yyyy-MM-dd"},
			} {
				if err := runner.Validate(schema.Require(check.name, check.value, check.hint)); err != nil {
					return err
				}
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "creator",
				Capability: "analytics",
				Operation:  "trends",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/influencer/trend",
				Params: map[string]string{
					"user_id":    creatorID,
					"start_date": startDate,
					"end_date":   endDate,
					"page_num":   page,
					"page_size":  pageSize,
				},
			})
		},
	}
	cmd.Flags().StringVar(&creatorID, "creator-id", "", "creator user_id")
	cmd.Flags().StringVar(&startDate, "start-date", "", "start date, yyyy-MM-dd")
	cmd.Flags().StringVar(&endDate, "end-date", "", "end date, yyyy-MM-dd")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newLeaderboard() *cobra.Command {
	cmd := &cobra.Command{Use: "leaderboard", Short: "Creator ranking and trend discovery"}
	cmd.AddCommand(newLeaderboardTop())
	return cmd
}

func newLeaderboardTop() *cobra.Command {
	var region, date, categoryName, productCategoryID, rankType, rankField, page, pageSize string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Fetch creator ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik creator leaderboard top --region US --date 2026-01-01")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("date", date, "rank date in yyyy-MM-dd format")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "creator",
				Capability: "leaderboard",
				Operation:  "top",
				Freshness:  "ranking",
				Path:       "/api/v3/echotik/influencer/ranklist",
				RecordsAs:  "creator",
				Params: schema.Clean(map[string]string{
					"date":                     date,
					"region":                   region,
					"influencer_category_name": categoryName,
					"product_category_id":      productCategoryID,
					"rank_type":                schema.RankType(rankType),
					"influencer_rank_field":    rankField,
					"page_num":                 page,
					"page_size":                pageSize,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&date, "date", "", "rank date in yyyy-MM-dd format")
	cmd.Flags().StringVar(&categoryName, "category-name", "", "creator category name")
	cmd.Flags().StringVar(&productCategoryID, "category-id", "", "product category ID")
	cmd.Flags().StringVar(&rankType, "type", "daily", "ranking type: daily|weekly|monthly or 1|2|3")
	cmd.Flags().StringVar(&rankField, "rank-field", "1", "ranking field: 1=total_followers_cnt, 2=total_sale_cnt")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "20", "page size")
	return cmd
}

func bindCreatorRelationFlags(cmd *cobra.Command, creatorID, uniqueID, page, pageSize *string) {
	cmd.Flags().StringVar(creatorID, "creator-id", "", "creator user_id")
	cmd.Flags().StringVar(uniqueID, "unique-id", "", "creator unique_id")
	cmd.Flags().StringVar(page, "page", "1", "page number")
	cmd.Flags().StringVar(pageSize, "page-size", "10", "page size")
}
