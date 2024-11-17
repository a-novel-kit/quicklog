package loggers

import (
	"log"
	"os"
	"sync"

	"github.com/rs/zerolog"

	"github.com/a-novel-kit/quicklog"
)

type zerologLogger struct {
	// True while an animated log is running. Prevents concurrent logs that cause management issues.
	animated bool

	logger zerolog.Logger

	quicklog.Logger
}

func (logger *zerologLogger) getEvent(level quicklog.Level) *zerolog.Event {
	switch level {
	case quicklog.LevelError:
		return logger.logger.Error()
	case quicklog.LevelWarning:
		return logger.logger.Warn()
	case quicklog.LevelFatal:
		return logger.logger.Fatal()
	default:
		return logger.logger.Info()
	}
}

func (logger *zerologLogger) checkAnimationLock() {
	if logger.animated {
		log.New(os.Stderr, "", 0).Fatal("cannot log while an animated message is running")
	}
}

func (logger *zerologLogger) Log(level quicklog.Level, message quicklog.Message) {
	logger.checkAnimationLock()

	rendered := message.RenderJSON()
	if rendered == nil {
		return
	}

	event := logger.getEvent(level)
	event.Fields(rendered).Msg("")
}

func (logger *zerologLogger) LogAnimated(message quicklog.AnimatedMessage) func() {
	logger.checkAnimationLock()

	logger.animated = true

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(1)

	cleaner := func() {
		logger.animated = false
		message.Close()
		waitGroup.Wait()
	}

	go func() {
		defer waitGroup.Done()

		for logMessage := range message.RunJSON() {
			if logMessage == nil {
				continue
			}

			logger.logger.Info().Fields(logMessage).Msg("")
		}
	}()

	return cleaner
}

// NewZerolog creates a new logger using the zerolog library.
func NewZerolog(logger zerolog.Logger) quicklog.Logger {
	return &zerologLogger{
		logger: logger,
	}
}
