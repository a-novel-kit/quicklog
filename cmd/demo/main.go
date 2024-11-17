package main

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/quicklog/messages"
)

var examples = map[string][]string{
	"Title": {
		messages.NewTitle("Amazing quicklog demo!", "Looks cozy in there...", nil).RenderTerminal(),
		messages.NewTitle("I am a title, just a bit lonely sabishii...", "", nil).RenderTerminal(),
	},
	"Error": {
		messages.NewError(fmt.Errorf("We made a fucky wucky!"), "Oopsie Woopsie!!").RenderTerminal(),
		messages.NewError(nil, "Oopsie Woopsie!!").RenderTerminal(),
		messages.NewError(fmt.Errorf("We made a fucky wucky!"), "").RenderTerminal(),
	},
	"Base": {
		messages.NewBase(
			"I am an example of base message! See how I am beautifully cropped for better readability. "+
				"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Phasellus rutrum eu arcu id tincidunt. "+
				"Maecenas rutrum porta felis, id scelerisque lectus dapibus et. Vivamus blandit tristique mauris, ac "+
				"finibus dui sollicitudin nec. Sed malesuada augue sit amet libero lacinia venenatis. Nam tincidunt "+
				"lectus lorem, eget tempus tellus euismod dapibus. Ut orci sapien, sagittis non accumsan at, "+
				"bibendum at elit. Nunc facilisis orci tortor, ac egestas leo interdum in. Maecenas pharetra enim "+
				"ac ante interdum, eu euismod ligula sodales.",
			nil,
		).RenderTerminal(),
	},
	"Loader": {
		func() string {
			loader := messages.NewLoader("On my way to do your mom!", &messages.LoaderConfigDefault)
			defer loader.Close()

			var message string
			defer testutils.CaptureChan(loader.RunTerminal(true), &message)()

			return message
		}(),
		func() string {
			loader := messages.NewLoader("On my way to do your mom!", &messages.LoaderConfigDefault)
			defer loader.Close()

			var message string
			defer testutils.CaptureChan(loader.RunTerminal(true), &message)()

			loader.Error(errors.New("Oops met your dad instead!"))

			return message
		}(),
		func() string {
			loader := messages.NewLoader("On my way to do your mom!", &messages.LoaderConfigDefault)
			defer loader.Close()

			var message string
			defer testutils.CaptureChan(loader.RunTerminal(true), &message)()

			loader.Success("Looks like you have a new mommy!")

			return message
		}(),
	},
}

func main() {
	fmt.Println()

	var orderedKeys []string
	for key := range examples {
		orderedKeys = append(orderedKeys, key)
	}
	sort.Strings(orderedKeys)

	for _, name := range orderedKeys {
		fmt.Println(lipgloss.NewStyle().Faint(true).Render(name + ":"))

		for _, message := range examples[name] {
			fmt.Println("\t" + strings.ReplaceAll(message, "\n", "\n\t"))
		}
	}
}
