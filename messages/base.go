package messages

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/a-novel-kit/quicklog"
)

type baseMessage struct {
	message string

	child quicklog.Message

	quicklog.Message
}

func (base *baseMessage) RenderTerminal() string {
	if base.message == "" {
		return ""
	}

	content := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Width(quicklog.TermWidth).
		Render(base.message)

	return quicklog.RenderWithChildTerminal(content+"\n", base.child)
}

func (base *baseMessage) RenderJSON() map[string]interface{} {
	if base.message == "" {
		return nil
	}

	content := map[string]interface{}{
		"message": base.message,
	}

	return quicklog.RenderWithChildJSON(content, base.child)
}

// NewBase groups together important logs under a section. Description and child are optional.
func NewBase(message string, child quicklog.Message) quicklog.Message {
	return &baseMessage{
		message: message,
		child:   child,
	}
}
