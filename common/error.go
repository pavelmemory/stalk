package common

import (
	"bytes"
	"fmt"
)

// this error type will be returned if declaration errors were found or some errors will appear on processing of args
type Error struct {
	// represents general cause of error and must be used for comparing
	Cause ErrorCode
	// contains additional information about the error
	// this value depends on context: flag name declaration error, invalid action, etc.
	ContextMessage string
}

// returns string representation of the error
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

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotImplementedError(msg string) Error {
	return Error{Cause: ErrorNotImplemented, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotAllRequiredValuesError(msg string) Error {
	return Error{Cause: ErrorNotAllRequiredValues, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func NotAllRequiredFlagsError(msg string) Error {
	return Error{Cause: ErrorNotAllRequiredFlags, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagSyntaxError(msg string) Error {
	return Error{Cause: ErrorFlagSyntax, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNotSupportedError(msg string) Error {
	return Error{Cause: ErrorFlagNotSupported, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutInvalidError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutInvalid, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutNotUniqueError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutNotUnique, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNameNotUniqueError(msg string) Error {
	return Error{Cause: ErrorFlagNameNotUnique, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagNameInvalidError(msg string) Error {
	return Error{Cause: ErrorFlagNameInvalid, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func FlagShortcutNameSameError(msg string) Error {
	return Error{Cause: ErrorFlagShortcutNameSame, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func CommandNameInvalidError(msg string) Error {
	return Error{Cause: ErrorCommandNameInvalid, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func CommandNameNotUniqueError(msg string) Error {
	return Error{Cause: ErrorCommandNameNotUnique, ContextMessage: msg}
}

// returns Error with corresponding function name ErrorCode and provided msg as ContextMessage
func ActionInvalidError(msg string) Error {
	return Error{Cause: ErrorActionInvalid, ContextMessage: msg}
}

// represents general cases of errors
type ErrorCode byte

const (
	errorNotSpecified ErrorCode = iota

	// signals that this functionality not yet implemented
	ErrorNotImplemented
	// signals that not all values required by declarations were passed for processing
	ErrorNotAllRequiredValues
	// signals that not all flags that declared with 'Required(true)' were passed for processing
	ErrorNotAllRequiredFlags
	// signals that flag name passed for processing has invalid syntax
	ErrorFlagSyntax
	// signals that flag passed for processing was not declared
	ErrorFlagNotSupported
	// signals that flag declaration contains invalid shortcut value
	ErrorFlagShortcutInvalid
	// signals that flag declarations have collision by shortcut value
	ErrorFlagShortcutNotUnique
	// signals that declared flag name value has collision with declared flag shortcut value (with not the same flag)
	ErrorFlagShortcutNameSame
	// signals that flag declaration contains invalid name value
	ErrorFlagNameInvalid
	// signals that flag declarations have collision by name value
	ErrorFlagNameNotUnique
	// signals that command declaration contains invalid name value
	ErrorCommandNameInvalid
	// signals that command declarations have collision by name value
	ErrorCommandNameNotUnique
	// signals that action is not a valid action (usually 'nil' value)
	ErrorActionInvalid
)

// returns string representation for ErrorCode values
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

// abstraction under error slice that used to pass found declaration errors as a single error
type DeclarationErrors []error

// returns string representation of errors separated by new line
func (de DeclarationErrors) Error() string {
	buf := bytes.Buffer{}
	for _, err := range de[:len(de)-1] {
		buf.WriteString(err.Error())
		buf.WriteRune('\n')
	}
	buf.WriteString(de[len(de)-1].Error())
	return buf.String()
}
