package doctor

import (
	"os"

	"github.com/echotik/cli/internal/core"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check local EchoTik CLI setup",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, user, pass, err := core.ResolveCredential()
			return output.Success(cmd.OutOrStdout(), map[string]any{
				"configPath":       core.ConfigPath(),
				"configExists":     fileExists(core.ConfigPath()),
				"hasEnvCredential": os.Getenv("ECHOTIK_USERNAME") != "" && os.Getenv("ECHOTIK_PASSWORD") != "",
				"hasCredential":    user != "" && pass != "",
				"ready":            err == nil,
				"hint":             hint(err),
			}, nil)
		},
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func hint(err error) string {
	if err == nil {
		return ""
	}
	return "run: echotik config set-credential"
}
