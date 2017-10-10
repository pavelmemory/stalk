package common

import (
	"fmt"
)

var (
	DefaultCommandStringer = func(declaration Declaration) string {
		return "name: " + declaration.GetName()
	}
)

// declaration of command to be used as part of workflow
type Declaration interface {
	// returns name given command in time of creation
	GetName() string
	// set action to be taken as main command task
	Execute(action func(ctx Runtime) error) Declaration
	// returns main action of command
	GetExecution() func(ctx Runtime) error
	// set flags applicable to this command
	Flags(flags ...Flag) Declaration
	// returns flags applicable to this command
	GetFlags() []Flag
	// set commands that can be used as child commands to current
	SubCommands(commands ...Declaration) Declaration
	// returns supported child commands of current command
	GetSubCommands() []Declaration
	// action to be executed before main command task or any child command
	Before(action func(ctx Runtime) error) Declaration
	// returns action that executes before main command task or any child command
	GetBefore() func(ctx Runtime) error
	// action to be executed after main command task or any child command
	After(action func(ctx Runtime, err error)) Declaration
	// returns action that executes after main command task or any child command
	GetAfter() func(ctx Runtime, err error)
	// action to be executed if command task ended with an error
	OnError(func(ctx Runtime, err error)) Declaration
	// returns action to be executed if command task ended with an error
	GetOnError()func(ctx Runtime, err error)
	// sets function used to convert command command to sting, `DefaultCommandStringer` used if not set
	Stringer(stringer func(command Declaration) string) Declaration
	// returns function used to convert flag to sting
	GetStringer() func(command Declaration) string
	fmt.Stringer
	// returns errors found in declaration of command
	GetDeclarationErrors() []error
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

// validates provided slice of flags and returns founded errors
func ValidateFlagDeclarations(flags []Flag) []error {
	var emptyShortcut rune
	var errs []error
	expectedFlagsByName := make(map[string]Flag)
	expectedFlagsByShortcut := make(map[rune]Flag)

	for _, flag := range flags {
		errs = append(errs, flag.GetDeclarationErrors()...)

		flagName := flag.GetName()
		if _, found := expectedFlagsByName[flagName]; found {
			errs = append(errs, FlagNameNotUniqueError(flag.String()))
		}
		expectedFlagsByName[flagName] = flag

		shortcut := flag.GetShortcut()
		if shortcut == emptyShortcut {
			continue
		}
		if _, found := expectedFlagsByShortcut[shortcut]; found {
			errs = append(errs, FlagShortcutNotUniqueError(flag.String()))
		}
		if foundFlag, found := expectedFlagsByName[string(shortcut)]; found && foundFlag != flag {
			errs = append(errs, FlagShortcutNameSameError(flag.String()+" and "+foundFlag.String()))
		}
		expectedFlagsByShortcut[shortcut] = flag
	}
	return errs
}

// validates provided slice of commands and returns founded errors
func ValidateCommandDeclarations(commands []Declaration) []error {
	var errs []error
	cmdDeclByName := make(map[string]Declaration)
	for _, cmd := range commands {
		cmdName := cmd.GetName()
		if _, found := cmdDeclByName[cmdName]; found {
			errs = append(errs, CommandNameNotUniqueError(cmdName))
		}
		errs = append(errs, cmd.GetDeclarationErrors()...)

		if cmd.GetExecution() == nil && len(cmd.GetSubCommands()) == 0 {
			errs = append(errs, ActionInvalidError("'"+cmdName + "' command has no action to execute neither sub-commands"))
		}
		cmdDeclByName[cmdName] = cmd
	}
	return errs
}
