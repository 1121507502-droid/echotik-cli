package live

import (
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{Use: "live", Short: "TikTok live stream intelligence"}
	cmd.AddCommand(newBasic())
	cmd.AddCommand(newAnalytics())
	cmd.AddCommand(newLeaderboard())
	return cmd
}

func newBasic() *cobra.Command {
	cmd := &cobra.Command{Use: "basic", Short: "Live discovery and basic information"}
	cmd.AddCommand(newSearch())
	cmd.AddCommand(newDetail())
	return cmd
}

func newSearch() *cobra.Command {
	var keyword, region, offset string
	cmd := &cobra.Command{
		Use:   "search",
		Short: "Search live streams in realtime",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("keyword", keyword, "example: echotik live basic search --keyword beauty --region US")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("region", region, "example: echotik live basic search --keyword beauty --region US")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "live",
				Capability: "basic",
				Operation:  "search",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/live/search",
				RecordsAs:  "live",
				Params: schema.Clean(map[string]string{
					"keyword": keyword,
					"region":  region,
					"offset":  offset,
				}),
			})
		},
	}
	cmd.Flags().StringVar(&keyword, "keyword", "", "search keyword")
	cmd.Flags().StringVar(&region, "region", "", "TikTok region, e.g. US")
	cmd.Flags().StringVar(&offset, "offset", "0", "pagination offset/cursor")
	return cmd
}

func newDetail() *cobra.Command {
	var roomID, creatorID string
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "Fetch realtime live stream details",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("room-id", roomID, "example: echotik live basic detail --room-id 123 --creator-id 456")); err != nil {
				return err
			}
			if err := runner.Validate(schema.Require("creator-id", creatorID, "example: echotik live basic detail --room-id 123 --creator-id 456")); err != nil {
				return err
			}
			return runner.Run(cmd, runner.Request{
				Entity:     "live",
				Capability: "basic",
				Operation:  "detail",
				Freshness:  "realtime",
				Path:       "/api/v3/realtime/live/detail",
				RecordsAs:  "live",
				Params: map[string]string{
					"room_id": roomID,
					"user_id": creatorID,
				},
			})
		},
	}
	cmd.Flags().StringVar(&roomID, "room-id", "", "live room ID")
	cmd.Flags().StringVar(&creatorID, "creator-id", "", "creator user_id")
	return cmd
}

func newAnalytics() *cobra.Command {
	cmd := &cobra.Command{Use: "analytics", Short: "Live relation analysis"}
	for _, op := range []string{"products", "trends", "media"} {
		operation := op
		cmd.AddCommand(&cobra.Command{
			Use:   operation,
			Short: "Reserved live analytics operation",
			RunE: func(cmd *cobra.Command, args []string) error {
				return runner.Unsupported(cmd, "live", "analytics", operation, "use live basic detail, or product/shop/creator analytics where live relations are available")
			},
		})
	}
	return cmd
}

func newLeaderboard() *cobra.Command {
	cmd := &cobra.Command{Use: "leaderboard", Short: "Live ranking and trend discovery"}
	cmd.AddCommand(&cobra.Command{
		Use:   "top",
		Short: "Reserved live leaderboard operation",
		RunE: func(cmd *cobra.Command, args []string) error {
			return runner.Unsupported(cmd, "live", "leaderboard", "top", "EchoTik docs do not expose a live leaderboard endpoint yet")
		},
	})
	return cmd
}
