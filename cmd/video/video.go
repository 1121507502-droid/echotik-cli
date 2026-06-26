package video

import (
	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{Use: "video", Short: "TikTok video intelligence"}
	cmd.AddCommand(newBasic())
	cmd.AddCommand(newAnalytics())
	cmd.AddCommand(newLeaderboard())
	return cmd
}

func newBasic() *cobra.Command {
	cmd := &cobra.Command{Use: "basic", Short: "Video discovery and basic information"}
	cmd.AddCommand(newSearch())
	cmd.AddCommand(newDetail())
	return cmd
}

func newSearch() *cobra.Command {
	var keyword, region, publishTime, sortType, offset string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search videos in realtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("keyword", keyword, "example: echotik video basic search --keyword sunscreen --region US")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("region", region, "example: echotik video basic search --keyword sunscreen --region US")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "basic",
				Operation:  "search",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/video/search",
				RecordsAs:  "video",
				Params: schema.Clean(map[string]string{
					"keyword":      keyword,
					"region":       region,
					"publish_time": publishTime,
					"sort_type":    sortType,
					"offset":       offset,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword")
	cmd.Flags().StringVar(&region, "region", "", "TikTok region, e.g. US")
	cmd.Flags().StringVar(&publishTime, "publish-time", "", "EchoTik publish_time filter")
	cmd.Flags().StringVar(&sortType, "sort-type", "", "0=relevance, 1=most likes")
	cmd.Flags().StringVar(&offset, "offset", "0", "pagination offset")
	return cmd
}

func newDetail() *cobra.Command {
	var videoID string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Fetch video details",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("video-id", videoID, "example: echotik video basic detail --video-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "basic",
				Operation:  "detail",
				Freshness:  "offline_t_plus_1",
				Path:       "/api/v3/echotik/video/detail",
				RecordsAs:  "video",
				Params: map[string]string{
					"video_ids": videoID,
				},
			})
		},
	}
	cmd.Flags().StringVar(&videoID, "video-id", "", "video ID; comma-separated max 10")
	return cmd
}

func newAnalytics() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Video relation analysis"}
	cmd.AddCommand(newProducts())
	cmd.AddCommand(newComments())
	cmd.AddCommand(newTrends())
	cmd.AddCommand(newMedia())
	return cmd
}

func newProducts() *cobra.Command {
	var videoID, page, pageSize string
	cmd := &cobra.Command{
		Use:   "products",
		Short: "Fetch products related to a video",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("video-id", videoID, "example: echotik video analytics products --video-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:       "video",
				Capability:   "analytics",
				Operation:    "products",
				Freshness:    "offline_t_plus_1",
				Path:         "/api/v3/echotik/video/product/list",
				RecordsAs:    "product",
				RelationFrom: &normalize.RelationNode{Type: "video", ID: videoID},
				RelationType: "video_promotes_product",
				Params: map[string]string{
					"video_ids": videoID,
					"page_num":  page,
					"page_size": pageSize,
				},
			})
		},
	}
	cmd.Flags().StringVar(&videoID, "video-id", "", "video ID; comma-separated values supported")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "10", "page size")
	return cmd
}

func newComments() *cobra.Command {
	var videoID, offset, count string
	cmd := &cobra.Command{
		Use:   "comments",
		Short: "Fetch realtime video comments",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("video-id", videoID, "example: echotik video analytics comments --video-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "analytics",
				Operation:  "comments",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/video/comments",
				Params: schema.Clean(map[string]string{
					"video_id": videoID,
					"offset":   offset,
					"count":    count,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&videoID, "video-id", "", "video ID")
	cmd.Flags().StringVar(&offset, "offset", "0", "pagination offset/cursor")
	cmd.Flags().StringVar(&count, "count", "20", "number of comments")
	return cmd
}

func newTrends() *cobra.Command {
	var videoID string
	cmd := &cobra.Command{
		Use:   "trends",
		Short: "Fetch realtime video interaction trend",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("video-id", videoID, "example: echotik video analytics trends --video-id 123")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "analytics",
				Operation:  "trends",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/video/trend_insight",
				Params: map[string]string{
					"video_id": videoID,
				},
			})
		},
	}
	cmd.Flags().StringVar(&videoID, "video-id", "", "video ID")
	return cmd
}

func newMedia() *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:   "media",
		Short: "Resolve video media URLs from a TikTok URL",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("url", url, "example: echotik video analytics media --url https://www.tiktok.com/...")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "analytics",
				Operation:  "media",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/video/download-url",
				RecordsAs:  "video",
				Params: map[string]string{
					"url": url,
				},
			})
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "TikTok web/app video URL")
	return cmd
}

func newLeaderboard() *cobra.Command {
	cmd := &cobra.Command{Use: "leaderboard", Short: "Video ranking and trend discovery"}
	cmd.AddCommand(newLeaderboardTop())
	return cmd
}

func newLeaderboardTop() *cobra.Command {
	var region, date, productCategoryID, createdByAI, rankType, rankField, page, pageSize string
	cmd := &cobra.Command{
		Use:   "top",
		Short: "Fetch video ranking list",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("region", region, "example: echotik video leaderboard top --region US --date 2026-01-01")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("date", date, "rank date in yyyy-MM-dd format")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "video",
				Capability: "leaderboard",
				Operation:  "top",
				Freshness:  "ranking",
				Path:       "/api/v3/echotik/video/ranklist",
				RecordsAs:  "video",
				Params: schema.Clean(map[string]string{
					"date":                date,
					"region":              region,
					"product_category_id": productCategoryID,
					"created_by_ai":       createdByAI,
					"video_rank_field":    rankField,
					"rank_type":           schema.RankType(rankType),
					"page_num":            page,
					"page_size":           pageSize,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&region, "region", "", "TikTok Shop region, e.g. US")
	cmd.Flags().StringVar(&date, "date", "", "rank date in yyyy-MM-dd format")
	cmd.Flags().StringVar(&productCategoryID, "category-id", "", "product category ID")
	cmd.Flags().StringVar(&createdByAI, "created-by-ai", "", "true/false")
	cmd.Flags().StringVar(&rankType, "type", "daily", "ranking type: daily|weekly|monthly or 1|2|3")
	cmd.Flags().StringVar(&rankField, "rank-field", "1", "ranking field: 1=total_views_cnt, 2=total_video_sale_cnt")
	cmd.Flags().StringVar(&page, "page", "1", "page number")
	cmd.Flags().StringVar(&pageSize, "page-size", "20", "page size")
	return cmd
}
