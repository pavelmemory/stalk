package flag

import (
	"github.com/pavelmemory/stalk/common"
	"testing"
)

func TestDefaultFlagStringer(t *testing.T) {
	for index, scenario := range []struct {
		expected string
		flag     common.Flag
	}{
		/*1*/ {"[--create]? [STRING]", String("create")},
		/*2*/ {"[--create|-c]? [STRING]", String("create").WithShortcut('c')},
		/*3*/ {"[--create|-c] [STRING]", String("create").WithShortcut('c').Required(true)},
		/*4*/ {"[--create|-c]? <STRING, def>", StringWithDefault("create", "def").WithShortcut('c')},
		/*5*/ {"[--verbose]?", Signal("verbose")},
		/*6*/ {"[--verbose|-v]?", Signal("verbose").WithShortcut('v')},
		/*7*/ {"[--verbose|-v]?", Signal("verbose").WithShortcut('v')},
	} {
		actual := DefaultFlagStringer(scenario.flag)
		if scenario.expected != actual {
			t.Error("index:", index + 1, "\nexpected:\n", scenario.expected, "\nactual:\n", actual)
		}
	}
}
