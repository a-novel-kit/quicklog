package quicklog

import (
	"log"
	"os"
	"strconv"

	"github.com/samber/lo"
)

// TermWidthEnv is the name of the environment variable that can be used to override the TermWidth value.
const TermWidthEnv = "TERM_WIDTH"

func getTermWidth() int {
	envWidth := os.Getenv(TermWidthEnv)
	if envWidth == "" {
		return 0
	}

	parsedWidth, err := strconv.Atoi(envWidth)
	if err != nil {
		log.Printf("Failed to parse TERM_WIDTH environment variable. Using default width: %s\n", err)
		return 0
	}

	return parsedWidth
}

// TermWidth is the default rendering width used for terminal output.
//
// It may be overridden by setting the TermWidthEnv environment variable.
var TermWidth = lo.CoalesceOrEmpty(getTermWidth(), 80)
