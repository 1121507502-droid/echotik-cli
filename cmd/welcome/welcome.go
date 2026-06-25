package welcome

import (
	"fmt"
	"io"
	"time"

	"github.com/spf13/cobra"
)

const (
	reset   = "\x1b[0m"
	dim     = "\x1b[2m"
	purple  = "\x1b[38;5;135m"
	violet  = "\x1b[38;5;99m"
	magenta = "\x1b[38;5;201m"
	black   = "\x1b[38;5;234m"
)

var pixelLogo = []string{
	"██████╗ ██████╗██╗  ██╗ ██████╗ ████████╗██╗██╗  ██╗",
	"██╔════╝██╔════╝██║  ██║██╔═══██╗╚══██╔══╝██║██║ ██╔╝",
	"█████╗  ██║     ███████║██║   ██║   ██║   ██║█████╔╝ ",
	"██╔══╝  ██║     ██╔══██║██║   ██║   ██║   ██║██╔═██╗ ",
	"███████╗╚██████╗██║  ██║╚██████╔╝   ██║   ██║██║  ██╗",
	"╚══════╝ ╚═════╝╚═╝  ╚═╝ ╚═════╝    ╚═╝   ╚═╝╚═╝  ╚═╝",
}

func New() *cobra.Command {
	var noAnimation bool

	cmd := &cobra.Command{
		Use:   "welcome",
		Short: "Show the EchoTik pixel logo",
		RunE: func(cmd *cobra.Command, args []string) error {
			if noAnimation {
				printLogo(cmd.OutOrStdout(), 0, "EchoTik CLI ready")
				return nil
			}
			revealLogo(cmd.OutOrStdout(), "EchoTik CLI ready")
			return nil
		},
	}

	cmd.Flags().BoolVar(&noAnimation, "no-animation", false, "print a static logo")
	return cmd
}

func printLogo(w io.Writer, offset int, subtitle string) {
	fmt.Fprintln(w)
	fmt.Fprintf(w, "%s▓▓%s%s EchoTik %s%s▓▓%s\n\n", black, reset, purple, reset, black, reset)
	for _, line := range pixelLogo {
		fmt.Fprintln(w, colorizeLine(line, offset))
	}
	fmt.Fprintf(w, "\n%s%s%s\n\n", dim, subtitle, reset)
}

func revealLogo(w io.Writer, subtitle string) {
	lines := logoLines(0, subtitle)
	for _, line := range lines {
		fmt.Fprintln(w, line)
		if line != "" {
			time.Sleep(35 * time.Millisecond)
		}
	}
}

func logoLines(offset int, subtitle string) []string {
	lines := []string{
		"",
		fmt.Sprintf("%s▓▓%s%s EchoTik %s%s▓▓%s", black, reset, purple, reset, black, reset),
		"",
	}
	for _, line := range pixelLogo {
		lines = append(lines, colorizeLine(line, offset))
	}
	lines = append(lines, "", fmt.Sprintf("%s%s%s", dim, subtitle, reset), "")
	return lines
}

func colorizeLine(line string, offset int) string {
	out := ""
	for i, ch := range line {
		if ch == ' ' {
			out += " "
			continue
		}

		switch phase := (i + offset) % 6; {
		case phase < 2:
			out += purple + string(ch)
		case phase < 4:
			out += violet + string(ch)
		default:
			out += magenta + string(ch)
		}
	}
	return out + reset
}
