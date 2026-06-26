package media

import (
	"github.com/echotik/cli/internal/media"
	"github.com/echotik/cli/internal/normalize"
	"github.com/echotik/cli/internal/output"
	"github.com/echotik/cli/internal/runner"
	"github.com/echotik/cli/internal/schema"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{Use: "media", Short: "Resolve and download EchoTik media assets"}
	cmd.AddCommand(newResolve())
	cmd.AddCommand(newDownload())
	return cmd
}

func newResolve() *cobra.Command {
	var url string
	cmd := &cobra.Command{
		Use:   "resolve",
		Short: "Resolve a media URL into an artifact",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("url", url, "example: echotik media resolve --url https://...")); err != nil {
				return err
			}
			data := normalize.Data{
				Records:   []any{map[string]string{"url": url}},
				Entities:  []normalize.Entity{},
				Relations: []normalize.Relation{},
				Artifacts: []normalize.Artifact{media.ArtifactFromURL(url)},
				Raw:       map[string]string{"url": url},
			}
			return output.Success(cmd.OutOrStdout(), data, map[string]interface{}{
				"entity":     "media",
				"capability": "basic",
				"operation":  "resolve",
				"freshness":  "local",
				"path":       "",
				"params":     map[string]string{"url": url},
			})
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "image/video/TOS URL")
	return cmd
}

func newDownload() *cobra.Command {
	var url, outputDir string
	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download a media asset and write a manifest",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := runner.Validate(schema.Require("url", url, "example: echotik media download --url https://... --output ./assets")); err != nil {
				return err
			}
			result, err := media.Download(url, outputDir)
			if err != nil {
				return output.NewError("download_error", err.Error(), "check whether the media URL expired; resolve it again if needed")
			}
			data := normalize.Data{
				Records:   []any{result},
				Entities:  []normalize.Entity{},
				Relations: []normalize.Relation{},
				Artifacts: result.Artifacts,
				Raw:       result,
			}
			return output.Success(cmd.OutOrStdout(), data, map[string]interface{}{
				"entity":     "media",
				"capability": "download",
				"operation":  "download",
				"freshness":  "local",
				"path":       result.Manifest,
				"params": map[string]string{
					"url":    url,
					"output": outputDir,
				},
			})
		},
	}
	cmd.Flags().StringVar(&url, "url", "", "image/video/TOS URL")
	cmd.Flags().StringVar(&outputDir, "output", "echotik-assets", "output directory")
	return cmd
}
