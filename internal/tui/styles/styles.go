package styles

import (
	"fmt"

	catppuccin "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (

	// https://github.com/catppuccin/catppuccin/blob/main/docs/style-guide.md

	// Styles
	BasicStyle = lipgloss.NewStyle().
			Foreground(catppuccinToLipgloss(catppuccin.Macchiato.Lavender())).
			Background(catppuccinToLipgloss(catppuccin.Macchiato.Base()))

		// LogoStyle is the style for the logo of the application.
	HeaderComponent = BasicStyle.
			PaddingRight(2).
			Width(100).
			MaxWidth(100).
			Align(lipgloss.Left)

	// TableStyle is the style for the table of the application.
	TableStyle = BasicStyle.
			BorderBackground(catppuccinToLipgloss(catppuccin.Macchiato.Base())).
			Border(lipgloss.RoundedBorder(), true)

	WelcomeStyle = TableStyle.
			Bold(true).
			AlignVertical(lipgloss.Center).
			Align(lipgloss.Center)

	TableHeader = BasicStyle.
			Bold(true).
			AlignVertical(lipgloss.Center).
			Background(catppuccinToLipgloss(catppuccin.Macchiato.Overlay2()))

	TableHeaderSelected = TableHeader.
				Background(catppuccinToLipgloss(catppuccin.Macchiato.Overlay1())).
				Foreground(catppuccinToLipgloss(catppuccin.Macchiato.Lavender()))
)

func catppuccinToLipgloss(c catppuccin.Color) lipgloss.Color {
	rgb := c.RGB
	// to hex string
	hex := fmt.Sprintf("#%02x%02x%02x", rgb[0], rgb[1], rgb[2])

	return lipgloss.Color(hex)
}
