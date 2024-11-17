package quicklog_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/a-novel-kit/quicklog"
)

type dummyMessage struct{}

func (d *dummyMessage) RenderTerminal() string {
	return "dummy"
}

func (d *dummyMessage) RenderJSON() map[string]interface{} {
	return map[string]interface{}{"dummy": true}
}

func TestRenderWithChildTerminal(t *testing.T) {
	testCases := []struct {
		name string

		parent string
		child  quicklog.Message

		expect string
	}{
		{
			name: "ParentEmpty",

			parent: "",
			child:  &dummyMessage{},

			expect: "",
		},
		{
			name: "ChildNil",

			parent: "parent",
			child:  nil,

			expect: "parent",
		},
		{
			name: "ParentAndChild",

			parent: "parent",
			child:  &dummyMessage{},

			expect: "parentdummy",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := quicklog.RenderWithChildTerminal(testCase.parent, testCase.child)
			require.Equal(t, testCase.expect, result)
		})
	}
}

func TestRenderWithChildJSON(t *testing.T) {
	testCases := []struct {
		name string

		parent map[string]interface{}
		child  quicklog.Message

		expect map[string]interface{}
	}{
		{
			name: "ParentEmpty",

			parent: nil,
			child:  &dummyMessage{},

			expect: nil,
		},
		{
			name: "ChildNil",

			parent: map[string]interface{}{"parent": true},
			child:  nil,

			expect: map[string]interface{}{"parent": true},
		},
		{
			name: "ParentAndChild",

			parent: map[string]interface{}{"parent": true},
			child:  &dummyMessage{},

			expect: map[string]interface{}{
				"parent": true,
				"data":   map[string]interface{}{"dummy": true},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := quicklog.RenderWithChildJSON(testCase.parent, testCase.child)
			require.Equal(t, testCase.expect, result)
		})
	}
}
