package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#2563EB")
	SecondaryColor = lipgloss.Color("#7C3AED")
	SuccessColor   = lipgloss.Color("#10B981")
	WarningColor   = lipgloss.Color("#F59E0B")
	DangerColor    = lipgloss.Color("#EF4444")
	LightColor     = lipgloss.Color("#F3F4F6")
	DarkColor      = lipgloss.Color("#1F2937")

	// App styles
	AppStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Background(DarkColor).
		Foreground(LightColor)

	// Header
	HeaderStyle = lipgloss.NewStyle().
		Background(PrimaryColor).
		Foreground(LightColor).
		Padding(0, 2).
		Height(3).
		Bold(true)

	HeaderLeftStyle = lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1)

	HeaderRightStyle = lipgloss.NewStyle().
		Faint(true)

	// Footer
	FooterStyle = lipgloss.NewStyle().
		Background(DarkColor).
		Foreground(LightColor).
		Padding(0, 2).
		Height(2)

	FooterLeftStyle = lipgloss.NewStyle().
		Faint(true)

	FooterRightStyle = lipgloss.NewStyle().
		Bold(true)

	// Navigation
	NavBarStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#374151")).
		Padding(0, 2).
		Height(2)

	NavItemStyle = lipgloss.NewStyle().
		Padding(0, 2).
		Foreground(LightColor)

	NavItemSelectedStyle = NavItemStyle.Copy().
		Background(SecondaryColor).
		Bold(true)

	// Sidebar
	SidebarStyle = lipgloss.NewStyle().
		Width(25).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor)

	// Content
	ContentStyle = lipgloss.NewStyle().
		Padding(1, 2)

	// Cards
	CardStyle = lipgloss.NewStyle().
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#4B5563")).
		Background(lipgloss.Color("#374151"))

	// Buttons
	ButtonStyle = lipgloss.NewStyle().
		Background(PrimaryColor).
		Foreground(LightColor).
		Padding(0, 3).
		Bold(true)

	// Forms
	FormStyle = lipgloss.NewStyle().
		Width(40).
		Padding(2, 3).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(PrimaryColor)

	InputStyle = lipgloss.NewStyle().
		Padding(0, 1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#6B7280")).
		Width(30)

	InputLabelStyle = lipgloss.NewStyle().
		Bold(true).
		MarginBottom(1)

	// Text
	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(PrimaryColor).
		MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
		Faint(true).
		MarginBottom(2)

	// Status
	LoadingStyle = lipgloss.NewStyle().
		Foreground(WarningColor).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(DangerColor).
		Bold(true)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(SuccessColor).
		Bold(true)

	// Shelf styles
	ShelfPlankStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("#5C4033")).
		Foreground(lipgloss.Color("#F5DEB3")).
		Bold(true).
		Padding(0, 1)

	BookStyle = lipgloss.NewStyle().
		Width(6).
		Height(8).
		Align(lipgloss.Center).
		Border(lipgloss.RoundedBorder()).
		Padding(0, 1)

	BookSelectedStyle = BookStyle.Copy().
		BorderForeground(PrimaryColor).
		Background(lipgloss.Color("#374151"))
)