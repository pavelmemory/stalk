package common

import "fmt"

var (
	EmptyNameMessage    = "<empty name>"
	ShortcutNotProvided rune
)

// common interface that describes all common parts of different flag types
type Flag interface {
	// GetName returns name of the flag
	GetName() string
	// Shortcut sets short name of the flag
	WithShortcut(value rune) Flag
	// GetShortcut returns short name of the flag
	GetDeclaredShortcut() rune
	// Required sets this flag as required
	Required(value bool) Flag
	// IsRequired returns `true` if this flag is required
	IsDeclaredRequired() bool
	// HasDefault returns `true` if this flag has default value
	HasDefault() bool
	// Parse parses provided `value` to the specific to flag type and saves it as value
	Parse(value string) error
	// IsSignal returns `true` if this flag does'n need value to be provided
	IsDeclaredSignal() bool
	// Stringer sets function used to convert flag to sting, `DefaultFlagStringer` used if not set
	WithStringer(stringer func(flag Flag) string) Flag
	// GetStringer returns function used to convert flag to sting
	GetDeclaredStringer() func(flag Flag) string
	fmt.Stringer

	// Description sets logical description for this flag
	WithDescription(value string) Flag
	// GetDescription returns description message for this flag
	GetDeclaredDescription() string
	// GetDeclarationErrors returns errors found in declaration of flag
	GetDeclarationErrors() []error
}

// Custom interface can be used for user-defined specific flags and used in creation of `Workflow`
type Custom interface {
	Flag
	// Value returns value parsed by `Parse` method and represents flag's value
	Value() interface{}
}

// ParsedString helper interface that supply `string` value
type ParsedString interface {
	// returns `string` value
	StringValue() string
}

// ParsedInt helper interface that supply `int64` value
type ParsedInt interface {
	// IntValue returns parsed value
	IntValue() int64
}

// ParsedBool helper interface that supply `bool` value
type ParsedBool interface {
	// BoolValue returns parsed value
	BoolValue() bool
}

// ParsedFloat helper interface that supply `float64` value
type ParsedFloat interface {
	// FloatValue returns parsed value
	FloatValue() float64
}
