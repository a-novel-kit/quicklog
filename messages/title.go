package messages

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/a-novel-kit/quicklog"
)

type titleMessage struct {
	title       string
	description string

	child quicklog.Message

	quicklog.Message
}

func (title *titleMessage) RenderTerminal() string {
	if title.title == "" {
		return ""
	}

	blockStyle := lipgloss.NewStyle().
		Width(quicklog.TermWidth).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("33")).
		Padding(0, 1)

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("33")).
		Bold(true).
		Render(title.title)

	if title.description != "" {
		content += "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("33")).
			Faint(true).
			Render(title.description)
	}

	return quicklog.RenderWithChildTerminal(blockStyle.Render(content)+"\n", title.child)
}

func (title *titleMessage) RenderJSON() map[string]interface{} {
	if title.title == "" {
		return nil
	}

	content := map[string]interface{}{
		"message": title.title,
	}

	if title.description != "" {
		content["content"] = title.description
	}

	return quicklog.RenderWithChildJSON(content, title.child)
}

// NewTitle groups together important logs under a section. Description and child are optional.
func NewTitle(title string, description string, child quicklog.Message) quicklog.Message {
	return &titleMessage{
		title:       title,
		description: description,
		child:       child,
	}
}
