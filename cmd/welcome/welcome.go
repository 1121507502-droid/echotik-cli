package welcome

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

const (
	reset  = "\x1b[0m"
	italic = "\x1b[3m"
	dim    = "\x1b[2m"
	purple = "\x1b[38;2;101;90;236m"
	muted  = "\x1b[38;2;128;123;164m"
	blue   = "\x1b[38;2;88;112;255m"
)

var italicLogo = []string{
	`   ______     __        ______     __  __     ______   __     __  __`,
	`  /\  ___\   /\ \      /\  ___\   /\ \_\ \   /\  __ \ /\ \   /\ \/ /`,
	`  \ \  __\   \ \ \____ \ \ \____  \ \  __ \  \ \ \/\ \\ \ \  \ \  _"-.`,
	`   \ \_____\  \ \_____\ \ \_____\  \ \_\ \_\  \ \_____\\ \_\  \ \_\ \_\`,
	`    \/_____/   \/_____/  \/_____/   \/_/\/_/   \/_____/ \/_/   \/_/\/_/`,
}

func New(cliVersion string) *cobra.Command {
	var noAnimation bool

	cmd := &cobra.Command{
		Use:   "welcome",
		Short: "Show the EchoTik pixel logo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if noAnimation {
				printLogo(cmd.OutOrStdout(), cliVersion, "EchoTik CLI ready")
				return nil
			}
			revealLogo(cmd.OutOrStdout(), cliVersion, "EchoTik CLI ready")
			return nil
		},
	}

	cmd.Flags().BoolVar(&noAnimation, "no-animation", false, "print a static logo")
	return cmd
}

func printLogo(w io.Writer, cliVersion, subtitle string) {
	for _, line := range logoLines(cliVersion, subtitle) {
		fmt.Fprintln(w, line)
	}
}

func revealLogo(w io.Writer, cliVersion, subtitle string) {
	lines := logoLines(cliVersion, subtitle)
	for _, line := range lines {
		fmt.Fprintln(w, line)
		if line != "" {
			time.Sleep(35 * time.Millisecond)
		}
	}
}

func logoLines(cliVersion, subtitle string) []string {
	lines := []string{""}
	for i, line := range italicLogo {
		suffix := ""
		switch i {
		case 1:
			suffix = fmt.Sprintf("  %sechotik cli %s%s", muted, displayVersion(cliVersion), reset)
		case 2:
			suffix = fmt.Sprintf("  %sready%s", blue, reset)
		}
		lines = append(lines, fmt.Sprintf("%s%s%s%s%s", italic, purple, line, reset, suffix))
	}
	lines = append(lines, "", fmt.Sprintf("%s%s%s", dim, subtitle, reset), "")
	return lines
}

func displayVersion(cliVersion string) string {
	if cliVersion == "" {
		return "vdev"
	}
	if cliVersion[0] == 'v' {
		return cliVersion
	}
	return "v" + cliVersion
}
