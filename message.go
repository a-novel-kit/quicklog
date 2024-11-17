package quicklog

// Message is a generic representation of a data that supports rendering under different formats.
type Message interface {
	// RenderTerminal renders a message in a format that is suitable for terminal output.
	RenderTerminal() string
	// RenderJSON renders a message in a format that is suitable for JSON output.
	RenderJSON() map[string]interface{}
}

// AnimatedMessage allow a Message to produce an output that is dynamically updated.
type AnimatedMessage interface {
	// RunTerminal starts rendering the message in dynamic mode. The returned channel is subscribed by the logger to
	// render the message in real-time.
	//
	// The CI flag is passed for environments where real-time outputs might result in loads of repetitive logs
	// (building docker image for example). In this case, only relevant updates should be sent to the chan.
	RunTerminal(ci bool) <-chan string

	// RunJSON is similar to RunTerminal but tailored for JSON output. The CI flag is assumed to always be true.
	RunJSON() <-chan map[string]interface{}

	// Close terminates the logger, and releases all its resources.
	Close()
}

// RenderWithChildTerminal automatically renders a parent with its child in terminal format.
func RenderWithChildTerminal(parent string, child Message) string {
	// If there is no parent message, then act as if nothing is logged. Child is an addon, not a replacement.
	if parent == "" {
		return ""
	}

	// No child = no change.
	if child == nil {
		return parent
	}

	return parent + child.RenderTerminal()
}

// RenderWithChildJSON automatically renders a parent with its child in JSON format.
func RenderWithChildJSON(parent map[string]interface{}, child Message) map[string]interface{} {
	// If there is no parent message, then act as if nothing is logged. Child is an addon, not a replacement.
	if parent == nil {
		return nil
	}

	// No child = no change.
	if child != nil {
		parent["data"] = child.RenderJSON()
	}

	return parent
}
