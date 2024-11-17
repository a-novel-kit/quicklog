package messages_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/quicklog"
	"github.com/a-novel-kit/quicklog/messages"
)

func TestBaseMessageTerminal(t *testing.T) {
	testCases := []struct {
		name string

		message string
		child   quicklog.Message

		expect string
	}{
		{
			name: "SimpleMessage",

			message: "Hello, world!",

			expect: "Hello, world!                                                                   \n",
		},
		{
			name: "NoMessage",

			expect: "",
		},
		{
			name: "WithChild",

			message: "Hello, world!",
			child:   messages.NewBase("Child message", nil),

			expect: "Hello, world!                                                                   \n" +
				"Child message                                                                   \n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewBase(testCase.message, testCase.child)
			require.Equal(t, testCase.expect, message.RenderTerminal())
		})
	}
}

func TestBaseMessageJSON(t *testing.T) {
	testCases := []struct {
		name string

		message string
		child   quicklog.Message

		expect map[string]interface{}
	}{
		{
			name: "SimpleMessage",

			message: "Hello, world!",

			expect: map[string]interface{}{
				"message": "Hello, world!",
			},
		},
		{
			name: "NoMessage",

			expect: nil,
		},
		{
			name: "WithChild",

			message: "Hello, world!",
			child:   messages.NewBase("Child message", nil),

			expect: map[string]interface{}{
				"message": "Hello, world!",
				"data": map[string]interface{}{
					"message": "Child message",
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewBase(testCase.message, testCase.child)
			require.Equal(t, testCase.expect, message.RenderJSON())
		})
	}
}
