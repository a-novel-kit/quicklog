package quicklog

// Level specify the importance of the log message. Some implementations may use different channels depending
// on the log level.
type Level string

const (
	// LevelInfo is the lowest log level. It is used for general information messages.
	LevelInfo Level = "INFO"
	// LevelWarning is used for messages that are not errors but may require attention.
	LevelWarning Level = "WARNING"
	// LevelError is used for messages that indicate an error occurred.
	LevelError Level = "ERROR"
	// LevelFatal is used for messages that indicate a fatal error occurred. A logger implementation should
	// automatically exit the program, or trigger a crash, after logging a message with this level.
	LevelFatal Level = "FATAL"
)

type Logger interface {
	// Log a message with the specified level.
	Log(level Level, message Message)

	// LogAnimated logs a message that can be updated in real-time.
	//
	// Running an animated log prevents new messages from being printed until the animated log is closed.
	// Attempting to log while an animated message is running will cause a panic.
	LogAnimated(message AnimatedMessage) (cleaner func())
}
