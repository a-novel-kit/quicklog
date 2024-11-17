package quicklog_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	testutils "github.com/a-novel-kit/test-utils"

	"github.com/a-novel-kit/quicklog"
)

func TestTermWidthDefault(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, 80, quicklog.TermWidth)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr)
		},
		Env: []string{"TERM_WIDTH="},
	})
}

func TestTermWidthOverride(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, 120, quicklog.TermWidth)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr)
		},
		Env: []string{"TERM_WIDTH=120"},
	})
}

func TestTermWidthOverrideInvalid(t *testing.T) {
	testutils.RunCMD(t, &testutils.CMDConfig{
		CmdFn: func(t *testing.T) {
			require.Equal(t, 80, quicklog.TermWidth)
		},
		MainFn: func(t *testing.T, res *testutils.CMDResult) {
			require.True(t, res.Success, res.STDErr)
		},
		Env: []string{"TERM_WIDTH=foobar"},
	})
}
