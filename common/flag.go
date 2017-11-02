package common

import "fmt"

var (
	EmptyNameMessage    = "<empty name>"
	ShortcutNotProvided rune
)

// common interface that describes all common parts of different flag types
type Flag interface {
	// name of the flag
	GetName() string
	// short name of the flag
	Shortcut(value rune) Flag
	// returns short name of the flag
	GetShortcut() rune
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

	DescriptionProvider(provider func(flag Flag) string) Flag
	GetDescriptionProvider() func(flag Flag) string
	// sets description message that will be printed in `help`
	Description(value string) Flag
	// returns description message if it was set
	GetDescription() string
	// returns errors found in declaration of flag
	GetDeclarationErrors() []error
}

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
