package common

// this context provides access to provided list of flags for each command and to the global flags
// it is also possible to use it as non-persistent key-value store between actions
type Runtime interface {
	// run commands with this context
	Run() error
	// returns list of arguments provided to the command
	GetArgs() []string
	// returns names of founded flags for command
	FoundFlags() []string
	// returns names of founded global flags
	FoundGlobalFlags() []string
	// checks if flag was found
	HasFlag(name string) bool
	// checks if global flag was found
	HasGlobalFlag(name string) bool
	// returns `string` value provided for flag with `name` name for command
	StringFlag(name string) string
	// returns `string` value provided for global flag with `name` name
	StringGlobalFlag(name string) string
	// returns `int` value provided for flag with `name` name for command
	IntFlag(name string) int64
	// returns `int` value provided for global flag with `name` name
	IntGlobalFlag(name string) int64
	// returns `bool` value provided for flag with `name` name for command
	BoolFlag(name string) bool
	// returns `bool` value provided for global flag with `name` name
	BoolGlobalFlag(name string) bool
	// returns `float64` value provided for flag with `name` name for command
	FloatFlag(name string) float64
	// returns `float64` value provided for global flag with `name` name
	FloatGlobalFlag(name string) float64
	// returns `interface{}` value provided for flag with `name` name for command
	CustomFlag(name string) interface{}
	// returns `interface{}` value provided for global flag with `name` name
	CustomGlobalFlag(name string) interface{}
	// stores provided key/value pair for future use
	// returns old value stored by this key and boolean flag if the key was overridden
	Set(key interface{}, value interface{}) (oldValue interface{}, overridden bool)
	// returns value stored previously by `Set` method
	// `found` flag shows if value was found
	Get(key interface{}) (value interface{}, found bool)
}
