package loggers_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/loggers"
	"github.com/a-novel-kit/quicklog/messages"
)

func TestTerminalLog(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewTerminal()

			logger.Log(quicklog.LevelInfo, messages.NewBase("This is an info message.", nil))

			// Ignore empty renders
			logger.Log(quicklog.LevelInfo, messages.NewBase("", nil))
			logger.Log(quicklog.LevelFatal, messages.NewBase("", nil))

			logger.Log(quicklog.LevelWarning, messages.NewBase("This is a warning message.", nil))
			logger.Log(quicklog.LevelError, messages.NewBase("This is an error message.", nil))
			logger.Log(quicklog.LevelFatal, messages.NewBase("This is a fatal message.", nil))
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
			require.Equal(
				t,
				"This is an info message.                                                        \n"+
					"This is a warning message.                                                      \n",
				res.STDOut,
			)
			require.Equal(
				t,
				"This is an error message.                                                       \n"+
					"This is a fatal message.                                                        \n",
				res.STDErr,
			)
		},
		Env: []string{"CI=true"},
	})
}

func TestTerminalLogAnimated(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewTerminal()

			logChan := make(chan string)
			animated := &fakeAnimated{outTerm: logChan}

			cleaner := logger.LogAnimated(animated)
			defer cleaner()

			logChan <- "This is an animated message."
			// Ignore empty renders.
			logChan <- ""
			logChan <- "This is another animated message."
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.Truef(t, res.Success, "stdout: %s\nstderr: %s", res.STDOut, res.STDErr)
			// Newlines are automatically appended if missing by the log library.
			// https://pkg.go.dev/log#Logger.Output
			require.Equal(t, "This is an animated message.\nThis is another animated message.\n", res.STDOut)
		},
		Env: []string{"CI=true"},
	})
}

func TestTerminalLogAnimatedNoConcurrentLog(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewTerminal()

			logChan := make(chan string)
			animated := &fakeAnimated{outTerm: logChan}

			cleaner := logger.LogAnimated(animated)
			defer cleaner()

			logChan <- "This is an animated message."

			// Attempting a concurrent log should crash the program.
			logger.Log(quicklog.LevelInfo, messages.NewBase("This is an info message.", nil))

			// Unreachable code.
			logChan <- "This is another animated message."
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success, "stdout: %s\nstderr: %s", res.STDOut, res.STDErr)

			require.Equal(t, "This is an animated message.\n", res.STDOut)
			require.Equal(t, "cannot log while an animated message is running\n", res.STDErr)
		},
		Env: []string{"CI=true"},
	})
}
