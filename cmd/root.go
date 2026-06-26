package cmd

import (
	"errors"
	"os"

	cmdapi "github.com/echotik/cli/cmd/api"
	cmdauth "github.com/echotik/cli/cmd/auth"
	cmdconfig "github.com/echotik/cli/cmd/config"
	cmdcreator "github.com/echotik/cli/cmd/creator"
	cmddoctor "github.com/echotik/cli/cmd/doctor"
	cmdlive "github.com/echotik/cli/cmd/live"
	cmdmedia "github.com/echotik/cli/cmd/media"
	cmdproduct "github.com/echotik/cli/cmd/product"
	cmdshop "github.com/echotik/cli/cmd/shop"
	cmdskills "github.com/echotik/cli/cmd/skills"
	cmdvideo "github.com/echotik/cli/cmd/video"
	cmdwelcome "github.com/echotik/cli/cmd/welcome"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

var version = "0.2.0-dev"

func Execute() int {
	root := NewRootCommand()
	if err := root.Execute(); err != nil {
		var cliErr *output.CLIError
		if errors.As(err, &cliErr) {
			_ = output.WriteCLIError(os.Stderr, cliErr)
			return output.ExitCode(cliErr.Type)
		}
		_ = output.Failure(os.Stderr, "command_error", err.Error(), "")
		return 1
	}
	return 0
}

func NewRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "echotik",
		Short: "EchoTik TikTok data intelligence CLI",
		Long: `echotik is an agent-friendly CLI for EchoTik TikTok data APIs.

It organizes access as entity/capability/operation commands for agent workflows,
plus a generic API command for direct EchoTik endpoint calls.`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Version:       version,
	}

	root.AddCommand(cmdconfig.New())
	root.AddCommand(cmdauth.New())
	root.AddCommand(cmddoctor.New())
	root.AddCommand(cmdapi.New())
	root.AddCommand(cmdproduct.New())
	root.AddCommand(cmdshop.New())
	root.AddCommand(cmdcreator.New())
	root.AddCommand(cmdvideo.New())
	root.AddCommand(cmdlive.New())
	root.AddCommand(cmdmedia.New())
	root.AddCommand(cmdwelcome.New())
	root.AddCommand(cmdskills.New())
	return root
}
