package flag

import (
	"github.com/pavelmemory/stalk/common"
	"testing"
)

func TestDefaultFlagUsageProvider(t *testing.T) {
	for index, scenario := range []struct {
		expected string
		flag     common.Flag
	}{
		{"flag 'create' needs to be used as:\n\t--create <VALUE>", String("create")},
		{"flag 'create' needs to be used as:\n\t--create <VALUE>\nor as a shortcut:\n\t-c <VALUE>", String("create").Shortcut('c')},
		{"flag 'verbose' needs to be used as:\n\t--verbose", Signal("verbose")},
		{"flag 'verbose' needs to be used as:\n\t--verbose\nor as a shortcut:\n\t-v", Signal("verbose").Shortcut('v')},
	} {
		actual := DefaultFlagUsageProvider(scenario.flag)
		if scenario.expected != actual {
			t.Error("index:", index, "\nexpected:\n", scenario.expected, "\nactual:\n", actual)
		}
	}
}

func TestDefaultFlagStringerProvider(t *testing.T) {
	for index, scenario := range []struct {
		expected string
		flag     common.Flag
	}{
		{"name: 'create'", String("create")},
		{"name: 'create', shortcut: 'c'", String("create").Shortcut('c')},
		{"name: 'create', shortcut: 'c', required", String("create").Shortcut('c').Required(true)},
		{"name: 'create', shortcut: 'c', required, default: 'def'", StringWithDefault("create", "def").Shortcut('c').Required(true)},
		{"name: 'verbose', signal", Signal("verbose")},
		{"name: 'verbose', shortcut: 'v', signal", Signal("verbose").Shortcut('v')},
		{"name: 'verbose', shortcut: 'v', signal, required", Signal("verbose").Shortcut('v').Required(true)},
		{"name: 'verbose', shortcut: 'v', signal, required, default: 'true'", SignalSetByDefault("verbose").Shortcut('v').Required(true)},
	} {
		actual := DefaultFlagStringerProvider(scenario.flag)
		if scenario.expected != actual {
			t.Error("index:", index, "\nexpected:\n", scenario.expected, "\nactual:\n", actual)
		}
	}
}
