package messages_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/messages"
)

func TestTitleTerminal(t *testing.T) {
	testCases := []struct {
		name string

		title       string
		description string
		child       quicklog.Message

		expect string
	}{
		{
			name: "SimpleTitle",

			title: "Hello, world!",

			expect: "╭────────────────────────────────────────────────────────────────────────────────╮\n" +
				"│ Hello, world!                                                                  │\n" +
				"╰────────────────────────────────────────────────────────────────────────────────╯\n",
		},
		{
			name: "TitleAndDescription",

			title:       "Hello, world!",
			description: "This is a description.",

			expect: "╭────────────────────────────────────────────────────────────────────────────────╮\n" +
				"│ Hello, world!                                                                  │\n" +
				"│ This is a description.                                                         │\n" +
				"╰────────────────────────────────────────────────────────────────────────────────╯\n",
		},
		{
			name: "TitleAndChild",

			title: "Hello, world!",
			child: messages.NewBase("Child message", nil),

			expect: "╭────────────────────────────────────────────────────────────────────────────────╮\n" +
				"│ Hello, world!                                                                  │\n" +
				"╰────────────────────────────────────────────────────────────────────────────────╯\n" +
				"Child message                                                                   \n",
		},
		{
			name: "NoMessage",

			expect: "",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewTitle(testCase.title, testCase.description, testCase.child)
			require.Equal(t, testCase.expect, message.RenderTerminal())
		})
	}
}

func TestTitleJSON(t *testing.T) {
	testCases := []struct {
		name string

		title       string
		description string
		child       quicklog.Message

		expect map[string]interface{}
	}{
		{
			name: "SimpleTitle",

			title: "Hello, world!",

			expect: map[string]interface{}{
				"message": "Hello, world!",
			},
		},
		{
			name: "TitleAndDescription",

			title:       "Hello, world!",
			description: "This is a description.",

			expect: map[string]interface{}{
				"message": "Hello, world!",
				"content": "This is a description.",
			},
		},
		{
			name: "TitleAndChild",

			title: "Hello, world!",
			child: messages.NewBase("Child message", nil),

			expect: map[string]interface{}{
				"message": "Hello, world!",
				"data": map[string]interface{}{
					"message": "Child message",
				},
			},
		},
		{
			name: "NoMessage",

			expect: nil,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewTitle(testCase.title, testCase.description, testCase.child)
			require.Equal(t, testCase.expect, message.RenderJSON())
		})
	}
}
