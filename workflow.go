package stalk

import (
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
)

// tree-structure representation of commands with flags and trigger-callbacks
type Workflow interface {
	// sets supported set of flags available to each command - global flags
	GlobalFlags(flags ...common.Flag) Workflow
	// returns set of flags available to each command - global flags
	GetGlobalFlags() []common.Flag
	// sets supported set of commands
	Commands(command ...common.Declaration) Workflow
	// returns supported set of commands
	GetCommands() []common.Declaration
	// sets function that will be executed only once before first command
	Setup(func(ctx common.Runtime) error) Workflow
	// returns function that will be executed only once before first command
	GetSetup() func(ctx common.Runtime) error
	// parses provided slice of strings that represents commands, flags, flag values and command arguments
	// and runs founded commands with founded or default flag values
	Run(cmd []string) error
	// after 'Run' execution returns errors found in declarations of global flags, commands and command flags
	GetDeclarationErrors() []error
	// sets function that will be executed only once after last command
	Cleanup(func(ctx common.Runtime, err error)) Workflow
	// returns function that will be executed only once after last command
	GetCleanup() func(ctx common.Runtime, err error)
	// sets function that will be executed only once if any command ends with error
	// this function won't be applied to declaration errors of type 'common.Error'
	OnError(func(ctx common.Runtime, err error)) Workflow
	// returns function that will be executed only once if any command ends with error
	GetOnError() func(ctx common.Runtime, err error)
	// sets flag declaration that will be used as a 'help' signal
	// only 'name' and 'shortcut' are valuable
	// if nil provided default 'help' will be unsupported
	HelpFlag(help common.Flag) Workflow
	// returns flag declaration that will be used as a 'help' signal
	// default help flag has name 'help' and shortcut 'h'
	GetHelpFlag() common.Flag
}

// creates new workflow that needs to be tuned with flags and commands
func New() Workflow {
	return &workflow{
		helpFlag: flag.Signal("help").Shortcut('h'),
	}
}

var _ Workflow = (*workflow)(nil)

type workflow struct {
	flags    []common.Flag
	commands []common.Declaration
	setup    func(ctx common.Runtime) error
	cleanup  func(ctx common.Runtime, err error)
	onError  func(ctx common.Runtime, err error)
	declErrs []error
	helpFlag common.Flag
}

func (w *workflow) Run(cmd []string) (err error) {
	if len(w.declErrs) != 0 {
		return common.DeclarationErrors(w.declErrs)
	}

	if len(cmd) == 0 {
		return nil
	}

	var runCtx common.Runtime
	defer func() {
		if err != nil && len(w.declErrs) == 0 {
			if onError := w.GetOnError(); onError != nil {
				onError(runCtx, err)
			}
		}

		cleanup := w.GetCleanup()
		if cleanup != nil {
			cleanup(runCtx, err)
		}
	}()

	runCtx, err = parse(w, cmd)
	if err != nil {
		return
	}

	if setup := w.GetSetup(); setup != nil {
		if err = setup(runCtx); err != nil {
			return
		}
	}

	err = runCtx.Run()
	return
}

func (w *workflow) Setup(action func(ctx common.Runtime) error) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': Setup"))
	}
	w.setup = action
	return w
}

func (w *workflow) GetSetup() func(ctx common.Runtime) error {
	return w.setup
}

func (w *workflow) Cleanup(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': Cleanup"))
	}
	w.cleanup = action
	return w
}

func (w *workflow) GetCleanup() func(ctx common.Runtime, err error) {
	return w.cleanup
}

func (w *workflow) OnError(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		w.declErrs = append(w.declErrs, common.ActionInvalidError("action is 'nil': OnError"))
	}
	w.onError = action
	return w
}

func (w *workflow) GetOnError() func(ctx common.Runtime, err error) {
	return w.onError
}

func (w *workflow) GlobalFlags(flags ...common.Flag) Workflow {
	w.declErrs = append(w.declErrs, common.ValidateFlagDeclarations(flags)...)
	w.flags = flags
	return w
}

func (w *workflow) GetGlobalFlags() []common.Flag {
	return w.flags
}

func (w *workflow) Commands(commands ...common.Declaration) Workflow {
	w.declErrs = append(w.declErrs, common.ValidateCommandDeclarations(commands)...)
	w.commands = commands
	return w
}

func (w *workflow) GetCommands() []common.Declaration {
	return w.commands
}

func (w *workflow) GetDeclarationErrors() []error {
	return w.declErrs
}

func (w *workflow) HelpFlag(helpFlag common.Flag) Workflow {
	w.helpFlag = helpFlag
	return w
}

func (w *workflow) GetHelpFlag() common.Flag {
	return w.helpFlag
}