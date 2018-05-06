package common

// Runtime provides access to provided list of flags for each command and to the global flags
// It is also possible to use it as non-persistent key-value store between actions
type Runtime interface {
	// Run execute tasks with this context
	Run() error
	// GetArgs returns list of arguments provided to the command
	GetArgs() []string
	// CurrentCommand returns command under execution
	CurrentCommand() ParsedCommand
	// Flags returns names of founded flags for command
	Flags() []string
	// GlobalFlags returns names of founded global flags
	GlobalFlags() []string
	// HasFlag returns `true` if flag with specified name was found in list of provided arguments
	HasFlag(name string) bool
	// HasGlobalFlag returns `true` if flag with specified name was found in list of provided arguments
	HasGlobalFlag(name string) bool
	// StringFlag returns `string` value provided for flag with `name` name for command
	StringFlag(name string) string
	// StringGlobalFlag returns `string` value provided for global flag with `name` name
	StringGlobalFlag(name string) string
	// IntFlag returns `int` value provided for flag with `name` name for command
	IntFlag(name string) int64
	// IntGlobalFlag returns `int` value provided for global flag with `name` name
	IntGlobalFlag(name string) int64
	// BoolFlag returns `bool` value provided for flag with `name` name for command
	BoolFlag(name string) bool
	// BoolGlobalFlag returns `bool` value provided for global flag with `name` name
	BoolGlobalFlag(name string) bool
	// FloatFlag returns `float64` value provided for flag with `name` name for command
	FloatFlag(name string) float64
	// FloatGlobalFlag returns `float64` value provided for global flag with `name` name
	FloatGlobalFlag(name string) float64
	// CustomFlag returns `interface{}` value provided for flag with `name` name for command
	CustomFlag(name string) interface{}
	// CustomGlobalFlag returns `interface{}` value provided for global flag with `name` name
	CustomGlobalFlag(name string) interface{}
	// Set stores provided key/value pair for future use
	// Returns old value stored under this key and boolean value `true` if value was overridden
	Set(key interface{}, value interface{}) (oldValue interface{}, overridden bool)
	// Get returns value stored previously by `Set` method with specified `key`
	// Second return value `found` shows if value was found
	Get(key interface{}) (value interface{}, found bool)
}
