package flag

import (
	"github.com/pavelmemory/stalk/common"
	"testing"
)

func TestDefaultFlagStringerProvider(t *testing.T) {
	for index, scenario := range []struct {
		expected string
		flag     common.Flag
	}{
		{"[--create] [STRING]", String("create")},
		{"[--create|-c] [STRING]", String("create").Shortcut('c')},
		{"[--create|-c]+ [STRING]", String("create").Shortcut('c').Required(true)},
		{"[--create|-c] <STRING, def>", StringWithDefault("create", "def").Shortcut('c')},
		{"[--verbose]", Signal("verbose")},
		{"[--verbose|-v]", Signal("verbose").Shortcut('v')},
		{"[--verbose|-v]", Signal("verbose").Shortcut('v')},
	} {
		actual := DefaultFlagStringerProvider(scenario.flag)
		if scenario.expected != actual {
			t.Error("index:", index, "\nexpected:\n", scenario.expected, "\nactual:\n", actual)
		}
	}
}
