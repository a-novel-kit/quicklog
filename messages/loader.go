package messages

import (
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/google/uuid"
	"github.com/samber/lo"

	"github.com/a-novel-kit/quicklog"
)

const EraseLineSequence = ansi.EraseEntireLine + "\r" + ansi.CursorUp1

type loaderStatus string

const (
	loaderStatusDefault loaderStatus = "running"
	loaderStatusSuccess loaderStatus = "success"
	loaderStatusError   loaderStatus = "error"
)

type Loader interface {
	quicklog.AnimatedMessage

	// Update sets a new message for the loader. It has no effect if called after Close, Success or Error.
	// If step is empty, the previous step will be re-rendered.
	Update(step string)

	// Nest adds more information to the loader in the form of an additional message.
	Nest(message quicklog.Message)

	// Success generates a success message, and closes the loader.
	// If step is empty, the previous step will be re-rendered.
	Success(step string)
	// Error generates an error message, and closes the loader.
	Error(err error)
}

type loaderMessage struct {
	renderTerminal chan string
	renderJSON     chan map[string]interface{}

	closed bool

	nested quicklog.Message

	// Keep track of the last rendered step message, for auto updates.
	lastStep string
	// Keep track of the last rendered terminal output, for post-processing.
	lastRenderedTerminal string
	// A flag to determine whether the loader is running in a CI environment.
	ci bool

	// Record the start time to show a timer after the message.
	startedAt time.Time
	// Display a custom spinner.
	spinner *spinner.Model
	// Record the last time spinner was updated. This helps trigger proper updates, according to fps parameter.
	spinnerLastUpdate time.Time
	// Allow logs to be grouped under JSON environments.
	opID uuid.UUID
	// Set the updater frequency for the elapsed timer.
	elapsedUpdateFrequency  time.Duration
	elapsedUpdateTicker     *time.Ticker
	elapsedUpdateTickerStop chan bool

	wait sync.WaitGroup
	mu   sync.Mutex

	quicklog.AnimatedMessage
}

// ==============================================================================================================
// Accessors.
// ==============================================================================================================

// Return whether the terminal channel is set.
func (loader *loaderMessage) hasTerminalChan() bool {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	if loader.renderTerminal == nil {
		return false
	}

	return !loader.closed
}

// Return whether the JSON channel is set.
func (loader *loaderMessage) hasJSONChan() bool {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	if loader.renderJSON == nil {
		return false
	}

	return !loader.closed
}

// Return the terminal channel if it is set. Otherwise, set a new one and return it.
func (loader *loaderMessage) getOrSetTerminalOutput() <-chan string {
	if !loader.hasTerminalChan() {
		loader.mu.Lock()
		loader.renderTerminal = make(chan string)
		loader.mu.Unlock()
	}

	return loader.renderTerminal
}

// Return the JSON channel if it is set. Otherwise, set a new one and return it.
func (loader *loaderMessage) getOrSetJSONOutput() <-chan map[string]interface{} {
	if !loader.hasJSONChan() {
		loader.mu.Lock()
		loader.renderJSON = make(chan map[string]interface{})
		loader.mu.Unlock()
	}

	return loader.renderJSON
}

func (loader *loaderMessage) getLastStep() string {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return loader.lastStep
}

func (loader *loaderMessage) isCI() bool {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	return loader.ci
}

// ==============================================================================================================
// Rendering.
// ==============================================================================================================

// Updates and return the loader view.
func (loader *loaderMessage) renderLoader() string {
	loader.mu.Lock()
	defer loader.mu.Unlock()

	if time.Since(loader.spinnerLastUpdate) > loader.spinner.Spinner.FPS {
		// Will be true on first render, prevent unnecessary updates.
		if loader.spinnerLastUpdate != (time.Time{}) {
			// Running the update method does not actually update the spinner but a copy of it (since it is not a pointer
			// method). So we have to assign the copu back to the original spinner afterward.
			newSpinner, _ := loader.spinner.Update(loader.spinner.Tick())
			*loader.spinner = newSpinner
		}

		loader.spinnerLastUpdate = time.Now()
	}

	return loader.spinner.View()
}

// Updates and return the time elapsed since the loader started running.
func (loader *loaderMessage) renderTimeElapsed() string {
	// Compute the time elapsed since the loader started running.
	loader.mu.Lock()
	timeElapsedRaw := time.Since(loader.startedAt)
	loader.mu.Unlock()

	// Prevent the display of values with large fractions.
	if timeElapsedRaw >= 10*time.Second {
		timeElapsedRaw = timeElapsedRaw.Round(time.Second)
	} else if timeElapsedRaw >= 10*time.Millisecond {
		timeElapsedRaw = timeElapsedRaw.Round(time.Millisecond)
	}

	return timeElapsedRaw.String()
}

// Send a new message to the terminal channel, if set.
//
// If no terminal channel is set, this method is a no-op.
func (loader *loaderMessage) updateTerminalOutput(step string, status loaderStatus) {
	if !loader.hasTerminalChan() {
		return
	}

	if step == "" {
		step = loader.getLastStep()
	}

	prefix := lo.Switch[loaderStatus, string](status).
		Case(loaderStatusSuccess, lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render("✓")).
		Case(loaderStatusError, lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render("✗")).
		DefaultF(func() string { return loader.renderLoader() })

	message := lo.Switch[loaderStatus, string](status).
		Case(loaderStatusSuccess, lipgloss.NewStyle().Foreground(lipgloss.Color("46")).Render(step)).
		Case(loaderStatusError, lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Render(step)).
		Default(lipgloss.NewStyle().Foreground(lipgloss.Color("15")).Render(step))

	mainMessage := prefix + " " + message

	timeElapsed := lipgloss.NewStyle().
		Faint(true).
		Foreground(lipgloss.Color("15")).
		Render(loader.renderTimeElapsed())

	timeElapsedMargin := lo.Max([]int{
		1,
		quicklog.TermWidth -
			((lipgloss.Width(mainMessage) + lipgloss.Width(timeElapsed)) % quicklog.TermWidth),
	})

	fullMessage := lipgloss.NewStyle().
		Width(quicklog.TermWidth).
		Render(mainMessage+lipgloss.NewStyle().MarginLeft(timeElapsedMargin).Render(timeElapsed)) + "\n"

	if loader.nested != nil {
		fullMessage += loader.nested.RenderTerminal()
	}

	loader.mu.Lock()
	lastRendered := loader.lastRenderedTerminal
	loader.mu.Unlock()

	// Only overwrite the previous message outside CI environment.
	if !loader.isCI() {
		// Get the previous content so it can be erased.
		loader.mu.Lock()
		loader.lastRenderedTerminal = fullMessage

		if lastRendered != "" {
			fullMessage = strings.Repeat(EraseLineSequence, lipgloss.Height(lastRendered)) + fullMessage
		}

		loader.mu.Unlock()
	}

	loader.renderTerminal <- fullMessage
}

// Send a new message to the JSON channel, if set.
//
// If no JSON channel is set, this method is a no-op.
func (loader *loaderMessage) updateJSONOutput(step string, status loaderStatus) {
	if !loader.hasJSONChan() {
		return
	}

	if step == "" {
		step = loader.getLastStep()
	}

	loader.mu.Lock()
	elapsedTime := time.Since(loader.startedAt)
	loader.mu.Unlock()

	output := map[string]interface{}{
		"message":       step,
		"elapsed":       elapsedTime.String(),
		"elapsed_nanos": elapsedTime.Nanoseconds(),
		"op_id":         loader.opID.String(),
		"status":        string(status),
	}

	if loader.nested != nil {
		loader.mu.Lock()
		output["data"] = loader.nested.RenderJSON()
		loader.mu.Unlock()
	}

	loader.renderJSON <- output
}

// ==============================================================================================================
// Loader state management.
// ==============================================================================================================

// Updates the value of the latest rendered step.
func (loader *loaderMessage) setLastStep(step string) {
	if step == "" {
		return
	}

	loader.mu.Lock()
	defer loader.mu.Unlock()
	loader.lastStep = step
}

// Close the previous ticker if any.
func (loader *loaderMessage) closeTicker() {
	if loader.elapsedUpdateTickerStop != nil {
		loader.elapsedUpdateTickerStop <- true
	}

	if loader.elapsedUpdateTicker != nil {
		loader.elapsedUpdateTicker.Stop()
	}

	loader.wait.Wait()
}

// Return the elapsed update ticker if it is set. Otherwise, set a new one and return it.
func (loader *loaderMessage) getOrSetTicker() *time.Ticker {
	if loader.elapsedUpdateTicker == nil {
		loader.mu.Lock()
		loader.elapsedUpdateTicker = time.NewTicker(loader.elapsedUpdateFrequency)
		loader.elapsedUpdateTickerStop = make(chan bool)
		loader.mu.Unlock()
	}

	return loader.elapsedUpdateTicker
}

// Periodically send new messages to the terminal channel, independently of user updates.
func (loader *loaderMessage) runAutoTerminalUpdates() {
	ticker := loader.getOrSetTicker()
	loader.wait.Add(1)

	go func() {
		for {
			select {
			case <-ticker.C:
				loader.updateTerminalOutput("", loaderStatusDefault)
			case <-loader.elapsedUpdateTickerStop:
				loader.wait.Done()
				return
			}
		}
	}()
}

// ==============================================================================================================
// Public methods.
// ==============================================================================================================

func (loader *loaderMessage) Nest(message quicklog.Message) {
	loader.mu.Lock()
	loader.nested = message
	loader.mu.Unlock()
}

func (loader *loaderMessage) Update(step string) {
	loader.updateTerminalOutput(step, loaderStatusDefault)
	loader.updateJSONOutput(step, loaderStatusDefault)
	loader.setLastStep(step)
}

func (loader *loaderMessage) Success(step string) {
	loader.closeTicker()

	loader.updateTerminalOutput(step, loaderStatusSuccess)
	loader.updateJSONOutput(step, loaderStatusSuccess)
	loader.setLastStep(step)
}

func (loader *loaderMessage) Error(err error) {
	loader.closeTicker()

	message := err.Error()

	loader.updateTerminalOutput(message, loaderStatusError)
	loader.updateJSONOutput(message, loaderStatusError)
	loader.setLastStep(message)
}

func (loader *loaderMessage) Close() {
	loader.closeTicker()

	if loader.hasTerminalChan() {
		close(loader.renderTerminal)
	}
	if loader.hasJSONChan() {
		close(loader.renderJSON)
	}

	loader.mu.Lock()
	loader.closed = true
	loader.mu.Unlock()
}

func (loader *loaderMessage) RunTerminal(isCI bool) <-chan string {
	loader.mu.Lock()
	loader.ci = isCI
	loader.mu.Unlock()

	channel := loader.getOrSetTerminalOutput()
	// Trigger initial rendering.
	go loader.Update("")

	// If outside CI environment, run periodic updates on our own. Otherwise, let the Update method provide relevant
	// updates.
	if !isCI {
		loader.runAutoTerminalUpdates()
	}

	return channel
}

func (loader *loaderMessage) RunJSON() <-chan map[string]interface{} {
	channel := loader.getOrSetJSONOutput()
	// Trigger initial rendering.
	go loader.Update("")

	return channel
}

// ==============================================================================================================
// Greeter.
// ==============================================================================================================

type LoaderConfig struct {
	// Optional.

	OpID            *uuid.UUID
	UpdateFrequency *time.Duration

	// Required.

	Spinner spinner.Model
}

var LoaderConfigDefault = LoaderConfig{
	Spinner: func() spinner.Model {
		loaderSpinner := spinner.New()
		loaderSpinner.Spinner = spinner.Meter
		loaderSpinner.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))

		return loaderSpinner
	}(),
}

func NewLoader(step string, config *LoaderConfig) Loader {
	loader := &loaderMessage{
		spinner:                &config.Spinner,
		lastStep:               step,
		opID:                   lo.Ternary(config.OpID != nil, lo.FromPtr(config.OpID), uuid.New()),
		startedAt:              time.Now(),
		elapsedUpdateFrequency: lo.CoalesceOrEmpty(lo.FromPtr(config.UpdateFrequency), 50*time.Millisecond),
	}

	return loader
}
