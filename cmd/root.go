package cmd

import (
	"errors"
	"fmt"
	"os"

	cmdapi "github.com/echotik/cli/cmd/api"
	cmdauth "github.com/echotik/cli/cmd/auth"
	cmdconfig "github.com/echotik/cli/cmd/config"
	cmddoctor "github.com/echotik/cli/cmd/doctor"
	cmdproduct "github.com/echotik/cli/cmd/product"
	cmdshop "github.com/echotik/cli/cmd/shop"
	cmdwelcome "github.com/echotik/cli/cmd/welcome"
	"github.com/echotik/cli/internal/output"
	"github.com/spf13/cobra"
)

var version = "0.1.0-dev"

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
		Short: "EchoTik TikTok Shop intelligence CLI",
		Long: `echotik is an agent-friendly CLI for EchoTik TikTok Shop APIs.

It provides high-level shortcuts for product and shop intelligence, plus a
generic API command for direct EchoTik endpoint calls.`,
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
	root.AddCommand(cmdwelcome.New())
	root.AddCommand(&cobra.Command{
		Use:   "skills",
		Short: "Print agent skill installation guidance",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(cmd.OutOrStdout(), "Install the bundled skills from this repository's skills/ directory into your agent.")
			return nil
		},
	})
	return root
}
