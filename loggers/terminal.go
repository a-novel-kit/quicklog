package loggers

import (
	"io"
	"log"
	"os"
	"sync"

	"github.com/a-novel-kit/quicklog"
)

const CIEnv = "CI"

type terminalLogger struct {
	ci bool

	// True while an animated log is running. Prevents concurrent logs that cause management issues.
	animated bool

	quicklog.Logger
}

func (logger *terminalLogger) getDestination(level quicklog.Level) io.Writer {
	if level == quicklog.LevelError {
		return os.Stderr
	}

	return os.Stdout
}

func (logger *terminalLogger) checkAnimationLock() {
	if logger.animated {
		log.New(os.Stderr, "", 0).Fatal("cannot log while an animated message is running")
	}
}

func (logger *terminalLogger) Log(level quicklog.Level, message quicklog.Message) {
	logger.checkAnimationLock()

	rendered := message.RenderTerminal()
	if rendered == "" {
		return
	}

	if level == quicklog.LevelFatal {
		log.New(os.Stderr, "", 0).Fatal(rendered)
		return
	}

	log.New(logger.getDestination(level), "", 0).Print(rendered)
}

func (logger *terminalLogger) LogAnimated(message quicklog.AnimatedMessage) func() {
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

		stdLogger := log.New(os.Stdout, "", 0)

		for logMessage := range message.RunTerminal(logger.ci) {
			if logMessage == "" {
				continue
			}

			stdLogger.Print(logMessage)
		}
	}()

	return cleaner
}

// NewTerminal creates a new Logger that logs to the terminal.
func NewTerminal() quicklog.Logger {
	return &terminalLogger{
		ci: os.Getenv(CIEnv) == "true",
	}
}
