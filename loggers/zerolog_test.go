package loggers_test

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/loggers"
	"github.com/a-novel-kit/quicklog/messages"
)

func TestZerologLog(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewZerolog(zerolog.New(os.Stdout))

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
				"{\"level\":\"info\",\"message\":\"This is an info message.\"}\n"+
					"{\"level\":\"warn\",\"message\":\"This is a warning message.\"}\n"+
					"{\"level\":\"error\",\"message\":\"This is an error message.\"}\n"+
					"{\"level\":\"fatal\",\"message\":\"This is a fatal message.\"}\n",
				res.STDOut,
			)
		},
	})
}

func TestZerologLogAnimated(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewZerolog(zerolog.New(os.Stdout))

			logChan := make(chan map[string]interface{})
			animated := &fakeAnimated{outJSON: logChan}

			cleaner := logger.LogAnimated(animated)
			defer cleaner()

			logChan <- messages.NewBase("This is an animated message.", nil).RenderJSON()
			// Ignore empty renders.
			logChan <- nil
			logChan <- messages.NewBase("This is another animated message.", nil).RenderJSON()
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success)
			require.Equal(
				t,
				"{\"level\":\"info\",\"message\":\"This is an animated message.\"}\n"+
					"{\"level\":\"info\",\"message\":\"This is another animated message.\"}\n",
				res.STDOut,
			)
		},
	})
}

func TestZerologLogAnimatedNoConcurrentLog(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			logger := loggers.NewZerolog(zerolog.New(os.Stdout))

			logChan := make(chan map[string]interface{})
			animated := &fakeAnimated{outJSON: logChan}

			cleaner := logger.LogAnimated(animated)
			defer cleaner()

			logChan <- messages.NewBase("This is an animated message.", nil).RenderJSON()

			// Attempting a concurrent log should crash the program.
			logger.Log(quicklog.LevelInfo, messages.NewBase("This is an info message.", nil))

			// Unreachable code.
			logChan <- messages.NewBase("This is another animated message.", nil).RenderJSON()
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.False(t, res.Success)
			require.Equal(
				t,
				"{\"level\":\"info\",\"message\":\"This is an animated message.\"}\n",
				res.STDOut,
			)
			require.Equal(t, "cannot log while an animated message is running\n", res.STDErr)
		},
	})
}
