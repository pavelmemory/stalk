package common

import (
	"fmt"
)

// CommandDeclaration is a declaration of command to be used as part of workflow
type CommandDeclaration interface {
	// GetName returns name given to command at creation time
	GetName() string
	// WithAction sets action to be taken as a main command task
	WithAction(action func(ctx Runtime) error) CommandDeclaration
	// GetDeclaredAction returns main command task
	GetDeclaredAction() func(ctx Runtime) error
	// WithFlags sets flags supported by this command
	WithFlags(flags ...Flag) CommandDeclaration
	// GetDeclaredFlags returns flags supported by this command
	GetDeclaredFlags() []Flag
	// WithSubCommands sets commands that can be used as child commands
	WithSubCommands(commands ...CommandDeclaration) CommandDeclaration
	// GetDeclaredSubCommands returns supported child commands of current command
	GetDeclaredSubCommands() []CommandDeclaration
	// WithBefore sets an action to be executed before this command task or any child command
	WithBefore(action func(ctx Runtime) error) CommandDeclaration
	// GetDeclaredBefore returns action that executes before main command task or any child command
	GetDeclaredBefore() func(ctx Runtime) error
	// WithAfter sets action to be executed after main command task or any child command
	WithAfter(action func(ctx Runtime, err error)) CommandDeclaration
	// GetDeclaredAfter returns action that executes after main command task or any child command
	GetDeclaredAfter() func(ctx Runtime, err error)
	// WithOnError sets action to be executed if command task return an error
	WithOnError(func(ctx Runtime, err error)) CommandDeclaration
	// GetDeclaredOnError returns action to be executed if command task return an error
	GetDeclaredOnError() func(ctx Runtime, err error)
	// WithStringer sets function used to convert command to sting, `DefaultCommandStringer` used if not set
	WithStringer(stringer func(command CommandDeclaration) string) CommandDeclaration
	// GetDeclaredStringer returns function used to convert command to sting
	GetDeclaredStringer() func(command CommandDeclaration) string
	fmt.Stringer
	// WithDescription sets logical description for this command
	WithDescription(value string) CommandDeclaration
	// GetDescription returns description for this command
	GetDeclaredDescription() string
	// GetDeclarationErrors returns errors found in declaration of command
	GetDeclarationErrors() []error
}

// Parsed represents command parsed from provided arguments list with supported flags and sub-commands
type ParsedCommand interface {
	CommandDeclaration
	// SubCommand sets command that used as child command to current
	SubCommand(command ParsedCommand)
	// GetSubCommand returns child command
	GetSubCommand() ParsedCommand
	// Flags sets flags that would be used to execute this command
	Flags(flags []Flag)
	// GetFlags returns flags
	GetFlags() map[string]Flag
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

		shortcut := flag.GetDeclaredShortcut()
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

// ValidateCommandDeclarations validates provided slice of commands and returns founded errors
func ValidateCommandDeclarations(commands []CommandDeclaration) []error {
	var errs []error
	cmdByName := make(map[string]CommandDeclaration)
	for _, cmd := range commands {
		cmdName := cmd.GetName()
		if _, found := cmdByName[cmdName]; found {
			errs = append(errs, CommandNameNotUniqueError(cmdName))
		}
		errs = append(errs, cmd.GetDeclarationErrors()...)

		if cmd.GetDeclaredAction() == nil && len(cmd.GetDeclaredSubCommands()) == 0 {
			errs = append(errs, ActionInvalidError("command '"+cmdName+"' has no action neither sub-commands to execute"))
		}
		cmdByName[cmdName] = cmd
	}
	return errs
}
