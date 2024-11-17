package messages

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/a-novel-kit/quicklog"
)

type errorMessage struct {
	err     error
	message string

	quicklog.Message
}

func (err *errorMessage) RenderTerminal() string {
	mainStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Width(quicklog.TermWidth)
	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).Background(lipgloss.Color("52")).
		Width(quicklog.TermWidth)

	if err.err == nil && err.message == "" {
		return ""
	}

	if err.message == "" {
		return mainStyle.Render(err.err.Error()) + "\n"
	}

	if err.err == nil {
		return messageStyle.Render(err.message) + "\n"
	}

	return messageStyle.Render(err.message) + "\n" + mainStyle.Render(err.err.Error()) + "\n"
}

func (err *errorMessage) RenderJSON() map[string]interface{} {
	if err.err == nil && err.message == "" {
		return nil
	}

	if err.message == "" {
		return map[string]interface{}{
			"message": err.err.Error(),
		}
	}

	if err.err == nil {
		return map[string]interface{}{
			"message": err.message,
		}
	}

	return map[string]interface{}{
		"message": err.message,
		"error":   err.err.Error(),
	}
}

// NewError creates a new error message.
func NewError(err error, message string) quicklog.Message {
	return &errorMessage{
		err:     err,
		message: message,
	}
}
