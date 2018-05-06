package common

import (
	"bytes"
	"fmt"
)

// Error type will be returned if declaration errors were found or some errors will appear on processing of args
type Error struct {
	// Cause represents general cause of error and must be used for comparing
	Cause ErrorCode
	// ContextMessage contains additional information about the error
	// This value depends on the context: flag name declaration error, invalid action, etc...
	ContextMessage string
}

// Error returns string representation of the error
func (e Error) Error() string {
	switch {
	case e.ContextMessage == "" && e.Cause == errorNotSpecified:
		return ""
	case e.Cause == errorNotSpecified:
		return e.ContextMessage
	case e.ContextMessage == "":
		return e.Cause.String()
	default:
		return e.Cause.String() + ": " + e.ContextMessage
	}
}

// NotImplementedError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotImplementedError(msg string) Error {
	return Error{Cause: ErrorNotImplemented, ContextMessage: msg}
}

// NotAllRequiredValuesError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotAllRequiredValuesError(msg string) Error {
	return Error{Cause: ErrorNotAllRequiredValues, ContextMessage: msg}
}

// NotAllRequiredFlagsError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotAllRequiredFlagsError(msg string) Error {
	return Error{Cause: ErrorNotAllRequiredFlags, ContextMessage: msg}
}

// FlagSyntaxError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagSyntaxError(msg string) Error {
	return Error{Cause: ErrorFlagSyntax, ContextMessage: msg}
}

// FlagNotSupportedError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNotSupportedError(msg string) Error {
	return Error{Cause: ErrorFlagNotSupported, ContextMessage: msg}
}

// FlagShortcutInvalidError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutInvalidError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutInvalid, ContextMessage: msg}
}

// FlagShortcutNotUniqueError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutNotUniqueError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutNotUnique, ContextMessage: msg}
}

// FlagNameNotUniqueError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNameNotUniqueError(msg string) Error {
	return Error{Cause: ErrorFlagNameNotUnique, ContextMessage: msg}
}

// FlagNameInvalidError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNameInvalidError(msg string) Error {
	return Error{Cause: ErrorFlagNameInvalid, ContextMessage: msg}
}

// FlagShortcutNameSameError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutNameSameError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutNameSame, ContextMessage: msg}
}

// FlagRequiredAndHasDefaultError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagRequiredAndHasDefaultError(msg string) Error {
	return Error{Cause: ErrorFlagRequiredAndHasDefault, ContextMessage: msg}
}

// FlagSignalAndRequiredError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagSignalAndRequiredError(msg string) Error {
	return Error{Cause: ErrorFlagSignalAndRequired, ContextMessage: msg}
}

// CommandNameInvalidError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func CommandNameInvalidError(msg string) Error {
	return Error{Cause: ErrorCommandNameInvalid, ContextMessage: msg}
}

// CommandNameNotUniqueError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func CommandNameNotUniqueError(msg string) Error {
	return Error{Cause: ErrorCommandNameNotUnique, ContextMessage: msg}
}

// ActionInvalidError returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func ActionInvalidError(msg string) Error {
	return Error{Cause: ErrorActionInvalid, ContextMessage: msg}
}

// ErrorCode represents general cases of errors
type ErrorCode byte

const (
	errorNotSpecified ErrorCode = iota

	// ErrorNotImplemented signals that this functionality not yet implemented
	ErrorNotImplemented
	// ErrorNotAllRequiredValues signals that not all values required by declarations were passed for processing
	ErrorNotAllRequiredValues
	// ErrorNotAllRequiredFlags signals that not all flags that declared with 'Required(true)' were passed for processing
	ErrorNotAllRequiredFlags
	// ErrorFlagSyntax signals that flag name passed for processing has invalid syntax
	ErrorFlagSyntax
	// ErrorFlagNotSupported signals that flag passed for processing was not declared
	ErrorFlagNotSupported
	// ErrorFlagShortcutInvalid signals that flag declaration contains invalid shortcut value
	ErrorFlagShortcutInvalid
	// ErrorFlagShortcutNotUnique signals that flag declarations have collision by shortcut value
	ErrorFlagShortcutNotUnique
	// ErrorFlagShortcutNameSame signals that declared flag name value has collision with declared flag shortcut value (with not the same flag)
	ErrorFlagShortcutNameSame
	// ErrorFlagNameInvalid signals that flag declaration contains invalid name value
	ErrorFlagNameInvalid
	// ErrorFlagNameNotUnique signals that flag declarations have collision by name value
	ErrorFlagNameNotUnique
	// ErrorFlagRequiredAndHasDefault signals that flag declaration defined as required, but has provided default value that make no sense
	ErrorFlagRequiredAndHasDefault
	// ErrorFlagSignalAndRequired signals that flag declaration defined as signal flag, so no sense to make it required, please use Bool flag to do so
	ErrorFlagSignalAndRequired
	// ErrorCommandNameInvalid signals that command declaration contains invalid name value
	ErrorCommandNameInvalid
	// ErrorCommandNameNotUnique signals that command declarations have collision by name value
	ErrorCommandNameNotUnique
	// ErrorActionInvalid signals that action is not a valid action (usually 'nil' value)
	ErrorActionInvalid
)

// String returns string representation for ErrorCode values
// or panic in case ErrorCode is not one of declared above
func (ec ErrorCode) String() string {
	if name, found := errorCodeNames[ec]; found {
		return name
	}
	panic("unknown value for ErrorCode type: " + fmt.Sprintf("%#v", ec))
}

var errorCodeNames = map[ErrorCode]string{
	ErrorNotImplemented:       "it is not implemented yet",
	ErrorNotAllRequiredValues: "not all required values provided",

	ErrorNotAllRequiredFlags:   "not all required flags provided",
	ErrorFlagSyntax:            "wrong flag syntax",
	ErrorFlagNotSupported:      "flag not supported",
	ErrorFlagShortcutInvalid:   "invalid flag shortcut",
	ErrorFlagShortcutNotUnique: "flag shortcut is not unique",

	ErrorFlagShortcutNameSame: "flag shortcut same to flag name",

	ErrorFlagNameInvalid:   "invalid flag name",
	ErrorFlagNameNotUnique: "flag name is not unique",

	ErrorCommandNameInvalid:   "invalid command name",
	ErrorCommandNameNotUnique: "command name is not unique",

	ErrorActionInvalid: "invalid action",
}

// DeclarationErrors is an abstraction under error slice that used to pass found declaration errors as a single error
type DeclarationErrors []error

// Error returns string representation of errors separated by new line
func (de DeclarationErrors) Error() string {
	buf := bytes.Buffer{}
	var separator rune
	for _, err := range de[:len(de)-1] {
		buf.WriteRune(separator)
		separator = '\n'
		buf.WriteString(err.Error())
	}
	buf.WriteString(de[len(de)-1].Error())
	return buf.String()
}
