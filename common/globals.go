package common

import (
	"errors"
	"fmt"
	"strconv"
)

var (
	ErrorNotAllRequiredFlags  = errors.New("not all required flags provided")
	ErrorNotAllRequiredValues = errors.New("not all required values provided")
	ErrorNotImplemented       = errors.New("it is not implemented yet")
	ErrorFlagSyntax           = errors.New("wrong flag syntax")
	ErrorFlagNotSupported     = errors.New("flag not supported")
)

// common interface that describes all common parts of different flag types
type Flag interface {
	// name of the flag
	GetName() string
	// short name of the flag
	Shortcut(value string) Flag
	// returns short name of the flag
	GetShortcut() string
	// set if this flag is required to be provided
	Required(value bool) Flag
	// returns `true` if this flag is required to be provided
	IsRequired() bool
	// returns `true` if this flag has provided default value
	HasDefault() bool
	// parse provided `value` to specific to flag type and save as it's value
	Proceed(value string) error
	// returns `true` if this flag does'n need value to be provided
	IsSignal() bool
	// sets function used to convert flag to sting, `DefaultFlagStringer` used if not set
	Stringer(stringer func(flag Flag) string) Flag
	// returns function used to convert flag to sting
	GetStringer() func(flag Flag) string
	fmt.Stringer
}

var (
	DefaultFlagStringer = func(flag Flag) string {
		return "name: " + flag.GetName() +
			", shortcut: " + flag.GetShortcut() +
			", is signal: " + strconv.FormatBool(flag.IsSignal()) +
			", is required: " + strconv.FormatBool(flag.IsRequired()) +
			", has default: " + strconv.FormatBool(flag.HasDefault())
	}

	DefaultCommandStringer = func(declaration Declaration) string {
		return "name: " + declaration.GetName()
	}
)

// anyone can implement this interface and use it in creation of `Workflow`
// if flag types provided by this tool is not enough
type Custom interface {
	Flag
	// returns value parsed by `Proceed` method and represents flag's value
	Value() interface{}
}

// helper interface
type ParsedString interface {
	// returns `string` value
	StringValue() string
}

// helper interface
type ParsedInt interface {
	// returns `int64` value
	IntValue() int64
}

// helper interface
type ParsedBool interface {
	// returns `bool` value
	BoolValue() bool
}

// helper interface
type ParsedFloat interface {
	// returns `float64` value
	FloatValue() float64
}

// declaration of declaration to be used as part of workflow
type Declaration interface {
	// returns name given declaration in time of creation
	GetName() string
	// set action to be taken as main declaration task
	Execute(action func(ctx Runtime) error) Declaration
	// returns main action of declaration
	GetExecution() func(ctx Runtime) error
	// set flags applicable to this declaration
	Flags(flags ...Flag) Declaration
	// returns flags applicable for this declaration
	GetFlags() []Flag
	// set commands that can be used as child commands to current
	SubCommands(commands ...Declaration) Declaration
	// returns supported child commands of current declaration
	GetSubCommands() []Declaration
	// action to be executed before main declaration task or any child declaration
	Before(action func(ctx Runtime) error) Declaration
	// returns action that executes before main declaration task or any child declaration
	GetBefore() func(ctx Runtime) error
	// action to be executed after main declaration task or any child declaration
	After(action func(ctx Runtime, err error)) Declaration
	// returns action that executes after main declaration task or any child declaration
	GetAfter() func(ctx Runtime, err error)
	// sets function used to convert command declaration to sting, `DefaultCommandStringer` used if not set
	Stringer(stringer func(declaration Declaration) string) Declaration
	// returns function used to convert flag to sting
	GetStringer() func(declaration Declaration) string
	fmt.Stringer
}

// represents command parsed from provided arguments list with supported flags and sub-commands
type Parsed interface {
	Declaration
	// sets commands that can be used as child commands to current
	SubCommand(command Parsed)
	// returns child command founded in provided arguments list
	GetSubCommand() Parsed
	// sets flags founded in provided arguments list
	FoundFlags(flags []Flag)
	// returns flags founded in provided arguments list
	GetFoundFlags() map[string]Flag
}

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
