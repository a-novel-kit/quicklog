package messages_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/quicklog/messages"
)

func TestErrorTerminal(t *testing.T) {
	testData := []struct {
		name string

		err     error
		message string

		expect string
	}{
		{
			name: "SimpleErrorMessage",

			err:     nil,
			message: "Hello, world!",

			expect: "Hello, world!                                                                   \n",
		},
		{
			name: "SimpleError",

			err: errors.New("this is an error"),

			expect: "this is an error                                                                \n",
		},
		{
			name: "ErrorAndMessage",

			err:     errors.New("this is an error"),
			message: "Hello, world!",

			expect: "Hello, world!                                                                   \n" +
				"this is an error                                                                \n",
		},
		{
			name: "NoMessage",

			expect: "",
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewError(testCase.err, testCase.message)
			require.Equal(t, testCase.expect, message.RenderTerminal())
		})
	}
}

func TestErrorJSON(t *testing.T) {
	testData := []struct {
		name string

		err     error
		message string

		expect map[string]interface{}
	}{
		{
			name: "SimpleErrorMessage",

			err:     nil,
			message: "Hello, world!",

			expect: map[string]interface{}{
				"message": "Hello, world!",
			},
		},
		{
			name: "SimpleError",

			err: errors.New("this is an error"),

			expect: map[string]interface{}{
				"message": "this is an error",
			},
		},
		{
			name: "ErrorAndMessage",

			err:     errors.New("this is an error"),
			message: "Hello, world!",

			expect: map[string]interface{}{
				"message": "Hello, world!",
				"error":   "this is an error",
			},
		},
		{
			name: "NoMessage",

			expect: nil,
		},
	}

	for _, testCase := range testData {
		t.Run(testCase.name, func(t *testing.T) {
			message := messages.NewError(testCase.err, testCase.message)
			require.Equal(t, testCase.expect, message.RenderJSON())
		})
	}
}
