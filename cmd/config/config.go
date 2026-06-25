package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Configure EchoTik credentials",
	}
	cmd.AddCommand(newSetCredential())
	cmd.AddCommand(newShow())
	return cmd
}

func newSetCredential() *cobra.Command {
	var username, password, baseURL string
	cmd := &cobra.Command{
		Use:   "set-credential",
		Short: "Store EchoTik Basic Auth credentials",
		RunE: func(cmd *cobra.Command, args []string) error {
			reader := bufio.NewReader(os.Stdin)
			if username == "" {
				fmt.Fprint(cmd.ErrOrStderr(), "EchoTik username: ")
				text, _ := reader.ReadString('\n')
				username = strings.TrimSpace(text)
			}
			if password == "" {
				fmt.Fprint(cmd.ErrOrStderr(), "EchoTik password: ")
				text, _ := reader.ReadString('\n')
				password = strings.TrimSpace(text)
			}
			if baseURL == "" {
				baseURL = core.DefaultBaseURL
			}
			if username == "" || password == "" {
				return output.NewError("validation_error", "missing username or password", "provide both --username and --password")
			}
			if err := core.Save(&core.Config{BaseURL: baseURL, Username: username, Password: password}); err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), map[string]any{
				"configured": true,
				"baseUrl":    baseURL,
				"configPath": core.ConfigPath(),
			}, nil)
		},
	}
	cmd.Flags().StringVar(&username, "username", "", "EchoTik username")
	cmd.Flags().StringVar(&password, "password", "", "EchoTik password")
	cmd.Flags().StringVar(&baseURL, "base-url", core.DefaultBaseURL, "EchoTik API base URL")
	return cmd
}

func newShow() *cobra.Command {
	return &cobra.Command{
		Use:   "show",
		Short: "Show current CLI configuration without secrets",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := core.Load()
			if err != nil {
				return err
			}
			return output.Success(cmd.OutOrStdout(), map[string]any{
				"baseUrl":       cfg.BaseURL,
				"hasCredential": cfg.Username != "" && cfg.Password != "",
				"configPath":     core.ConfigPath(),
			}, nil)
		},
	}
}
