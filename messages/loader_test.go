package messages_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/quicklog/messages"
)

var dummySpinner = spinner.Spinner{
	Frames: []string{"u", "w", "o"},
	FPS:    100 * time.Second,
}

var dummySpinnerModel = func() spinner.Model {
	spinnerModel := spinner.New()
	spinnerModel.Spinner = dummySpinner
	return spinnerModel
}()

var dummyOpID = uuid.MustParse("10000000-1000-1000-1000-100000000000")

var loaderTestConfig = &messages.LoaderConfig{
	Spinner:         dummySpinnerModel,
	OpID:            &dummyOpID,
	UpdateFrequency: lo.ToPtr(100 * time.Millisecond),
}

func TestLoaderTerminal(t *testing.T) {
	t.Run("RenderInitialMessage", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		testutils.RequireChan(t, loader.RunTerminal(true), func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})
	})

	t.Run("RenderUpdates", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})

		go loader.Update("updated message")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u updated message .+\n$`), value)
		})

		go loader.Update("updated message 2")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u updated message 2 .+\n$`), value)
		})
	})

	t.Run("RenderError", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})

		go loader.Error(errors.New("error message"))

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^✗ error message .+\n$`), value)
		})
	})

	t.Run("RenderSuccess", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})

		go loader.Success("success message")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^✓ success message .+\n$`), value)
		})
	})

	t.Run("RenderNested", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})

		loader.Nest(messages.NewBase("child message", nil))
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\nchild message\s+\n$`), value)
		})

		loader.Nest(nil)
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})
	})

	t.Run("RenderTimeElapsed", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		// Initial value should be rendered in nanoseconds.
		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message\s+\d{1,3}\.\d{1,3}µs\n$`), value)
		})

		// Waiting 10ms, time should be rounded to ms.
		time.Sleep(10 * time.Millisecond)
		// Update should be rendered in milliseconds.
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message\s+\d{1,3}ms\n$`), value)
		})

		// Waiting 3s, time should still be rounded to ms.
		time.Sleep(3 * time.Second)
		// Update should be rendered in milliseconds.
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message\s+\d{1,3}\.\d{0,3}s\n$`), value)
		})

		// Waiting 10s, time should be rounded to seconds.
		time.Sleep(7 * time.Second)
		// Update should be rendered in seconds.
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message\s+\d{1,2}s\n$`), value)
		})
	})

	t.Run("Loader", func(t *testing.T) {
		cfg := *loaderTestConfig
		cfg.Spinner.Spinner.FPS = 100 * time.Millisecond
		loader := messages.NewLoader("initial message", &cfg)
		defer loader.Close()

		channel := loader.RunTerminal(true)

		// The spinner view should not move for 100ms.
		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})

		time.Sleep(100 * time.Millisecond)
		go loader.Update("")

		// The spinner view should move after 100ms.
		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^w initial message .+\n$`), value)
		})

		time.Sleep(100 * time.Millisecond)
		go loader.Update("")

		// The spinner view should move after 100ms.
		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^o initial message .+\n$`), value)
		})

		// Back to square one.
		time.Sleep(100 * time.Millisecond)
		go loader.Update("")

		// The spinner view should move after 100ms.
		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message .+\n$`), value)
		})
	})

	t.Run("NonCIEnvironment", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunTerminal(false)

		testutils.RequireChanC(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^u initial message\s+\d{1,3}\.\d{1,3}µs\n$`), value)
		}, 50*time.Millisecond, 5*time.Millisecond)

		// After 100ms, the timer should be automatically updated.
		time.Sleep(100 * time.Millisecond)

		testutils.RequireChanC(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^.+u initial message\s+1\d{2}ms\n$`), value)
		}, 50*time.Millisecond, 5*time.Millisecond)

		// Update the message.
		go loader.Update("updated message")

		// The message should be updated.
		testutils.RequireChanC(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^.+u updated message\s+1\d{2}ms\n$`), value)
		}, 50*time.Millisecond, 5*time.Millisecond)

		// Wait a bit more.
		time.Sleep(100 * time.Millisecond)

		testutils.RequireChanC(t, channel, func(collect *assert.CollectT, value string) {
			assert.Regexp(collect, regexp.MustCompile(`^.+u updated message\s+2\d{2}ms\n$`), value)
		}, 50*time.Millisecond, 5*time.Millisecond)
	})
}

func TestLoaderJSON(t *testing.T) {
	t.Run("RenderInitialMessage", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		testutils.RequireChan(t, loader.RunJSON(), func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})
	})

	t.Run("RenderUpdates", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunJSON()

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})

		go loader.Update("updated message")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "updated message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})

		go loader.Update("updated message 2")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "updated message 2", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})
	})

	t.Run("RenderError", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunJSON()

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})

		go loader.Error(errors.New("error message"))

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "error message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "error", value["status"])
		})
	})

	t.Run("RenderSuccess", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunJSON()

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})

		go loader.Success("success message")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "success message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "success", value["status"])
		})
	})

	t.Run("RenderNested", func(t *testing.T) {
		loader := messages.NewLoader("initial message", loaderTestConfig)
		defer loader.Close()

		channel := loader.RunJSON()

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})

		loader.Nest(messages.NewBase("child message", nil))
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
			assert.Equal(collect, map[string]interface{}{"message": "child message"}, value["data"])
		})

		loader.Nest(nil)
		go loader.Update("")

		testutils.RequireChan(t, channel, func(collect *assert.CollectT, value map[string]interface{}) {
			assert.Equal(collect, "initial message", value["message"])
			assert.Regexp(collect, regexp.MustCompile(`^\d{1,3}(\.\d+)?(µs|ms|s)$`), value["elapsed"])
			assert.NotEmpty(collect, value["elapsed_nanos"])
			assert.Equal(collect, dummyOpID.String(), value["op_id"])
			assert.Equal(collect, "running", value["status"])
		})
	})
}
