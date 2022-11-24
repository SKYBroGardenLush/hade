package survey

import (
	"testing"

	"github.com/SKYBroGardenLush/skyscraper/framework/survey/terminal"
)

func RunTest(t *testing.T, procedure func(expectConsole), test func(terminal.Stdio) error) {
	t.Skip("warning: Windows does not support psuedoterminals")
}
