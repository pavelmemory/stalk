package stalk

import (
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
)

// Workflow describes tree-like structure representation of commands with flags and trigger-callbacks
type Workflow interface {
	// WithGlobalFlags sets supported set of flags available to each command - global flags
	WithGlobalFlags(flags ...common.Flag) Workflow
	// GetDeclaredGlobalFlags returns set of flags available to each command - global flags
	GetDeclaredGlobalFlags() []common.Flag
	// WithCommands sets supported set of commands
	WithCommands(command ...common.CommandDeclaration) Workflow
	// GetDeclaredCommands returns supported set of commands
	GetDeclaredCommands() []common.CommandDeclaration
	// WithSetup sets function that will be executed only once before first command
	// You may use it for operations such as open connection to database, etc...
	WithSetup(func(ctx common.Runtime) error) Workflow
	// GetSetup returns function that will be executed only once before first command
	GetDeclaredSetup() func(ctx common.Runtime) error
	// Run parses provided slice of strings that represents commands, flags, flag values and command arguments
	// and runs founded commands with founded or default flag values
	Run(cmd []string) error
	// GetDeclarationErrors returns errors found in declarations of global flags, commands and command flags after 'Run' execution
	GetDeclarationErrors() []error
	// WithCleanup sets function that will be executed only once after last command
	// You may use it for operations such as closing connection to database, etc...
	WithCleanup(func(ctx common.Runtime, err error)) Workflow
	// GetCleanup returns function that will be executed only once after last command
	GetDeclaredCleanup() func(ctx common.Runtime, err error)
	// WithOnError sets function that will be executed only once if any command ends with an error
	// This function won't be applied to declaration errors
	WithOnError(func(ctx common.Runtime, err error)) Workflow
	// GetOnError returns function that will be executed only once if any command ends with an error
	GetDeclaredOnError() func(ctx common.Runtime, err error)
	// WithHelpFlag sets flag declaration that will be used as a 'help' signal
	// Only 'name' and 'shortcut' are valuable
	// if nil provided default 'help' won't be supported
	WithHelpFlag(help common.Flag) Workflow
	// GetHelpFlag returns flag declaration that will be used as a 'help' signal
	// Returns default help flag is user-specific flag was not set
	// Default help flag has name 'help' and shortcut 'h'
	GetDeclaredHelpFlag() common.Flag
}

// creates new workflow that needs to be tuned with flags and commands
func New() Workflow {
	return &workflow{
		helpFlag: flag.Signal("help").WithShortcut('h'),
	}
}

var _ Workflow = (*workflow)(nil)

type workflow struct {
	flags    []common.Flag
	commands []common.CommandDeclaration
	setup    func(ctx common.Runtime) error
	cleanup  func(ctx common.Runtime, err error)
	onError  func(ctx common.Runtime, err error)
	declErrs []error
	helpFlag common.Flag
}

func (w *workflow) Run(cmd []string) (err error) {
	// execution impossible because of invalid declarations
	if len(w.declErrs) != 0 {
		return common.DeclarationErrors(w.declErrs)
	}

	// if no commands provided then we have nothing to execute
	if len(cmd) == 0 {
		return nil
	}

	var runCtx common.Runtime
	runCtx, err = parse(w, cmd)
	if err != nil {
		return
	}

	defer func() {
		// 3. if execution error happens handle it properly first
		if err != nil && len(w.declErrs) == 0 {
			if onError := w.GetDeclaredOnError(); onError != nil {
				onError(runCtx, err)
			}
		}

		// 4. and then run `Cleanup` if defined
		if cleanup := w.GetDeclaredCleanup(); cleanup != nil {
			cleanup(runCtx, err)
		}
	}()

	// 1. execution of user-defined action starts with `Setup` action
	if setup := w.GetDeclaredSetup(); setup != nil {
		if err = setup(runCtx); err != nil {
			return
		}
	}

	// 2. after `Setup` starts execution of all commands found in arguments
	err = runCtx.Run()
	return
}

func (w *workflow) WithSetup(action func(ctx common.Runtime) error) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': Setup"))
	}
	w.setup = action
	return w
}

func (w *workflow) GetDeclaredSetup() func(ctx common.Runtime) error {
	return w.setup
}

func (w *workflow) WithCleanup(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': Cleanup"))
	}
	w.cleanup = action
	return w
}

func (w *workflow) GetDeclaredCleanup() func(ctx common.Runtime, err error) {
	return w.cleanup
}

func (w *workflow) WithOnError(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': OnError"))
	}
	w.onError = action
	return w
}

func (w *workflow) GetDeclaredOnError() func(ctx common.Runtime, err error) {
	return w.onError
}

func (w *workflow) WithGlobalFlags(flags ...common.Flag) Workflow {
	w.declErrs = append(w.declErrs, common.ValidateFlagDeclarations(flags)...)
	w.flags = flags
	return w
}

func (w *workflow) GetDeclaredGlobalFlags() []common.Flag {
	return w.flags
}

func (w *workflow) WithCommands(commands ...common.CommandDeclaration) Workflow {
	w.declErrs = append(w.declErrs, common.ValidateCommandDeclarations(commands)...)
	w.commands = commands
	return w
}

func (w *workflow) GetDeclaredCommands() []common.CommandDeclaration {
	return w.commands
}

func (w *workflow) GetDeclarationErrors() []error {
	return w.declErrs
}

func (w *workflow) WithHelpFlag(helpFlag common.Flag) Workflow {
	w.helpFlag = helpFlag
	return w
}

func (w *workflow) GetDeclaredHelpFlag() common.Flag {
	return w.helpFlag
}