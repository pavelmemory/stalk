package stalk

import (
	"github.com/pavelmemory/stalk/common"
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
}

// creates new workflow that needs to be tuned with flags and commands
func New() Workflow {
	return &workflow{}
}

var _ Workflow = (*workflow)(nil)

type workflow struct {
	flags    []common.Flag
	commands []common.Declaration
	setup    func(ctx common.Runtime) error
	cleanup  func(ctx common.Runtime, err error)
	onError  func(ctx common.Runtime, err error)
	declErrs []error
}

func (app *workflow) Run(cmd []string) (err error) {
	if len(app.declErrs) != 0 {
		return common.DeclarationErrors(app.declErrs)
	}

	if len(cmd) == 0 {
		return nil
	}

	var runCtx common.Runtime
	defer func() {
		if err != nil && len(app.declErrs) == 0 {
			if onError := app.GetOnError(); onError != nil {
				onError(runCtx, err)
			}
		}

		cleanup := app.GetCleanup()
		if cleanup != nil {
			cleanup(runCtx, err)
		}
	}()

	runCtx, err = parse(app, cmd)
	if err != nil {
		return
	}

	if setup := app.GetSetup(); setup != nil {
		if err = setup(runCtx); err != nil {
			return
		}
	}

	err = runCtx.Run()
	return
}

func (app *workflow) Setup(action func(ctx common.Runtime) error) Workflow {
	if action == nil {
		app.declErrs = append(app.declErrs, common.ActionInvalidError("action is 'nil': Setup"))
	}
	app.setup = action
	return app
}

func (app *workflow) GetSetup() func(ctx common.Runtime) error {
	return app.setup
}

func (app *workflow) Cleanup(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		app.declErrs = append(app.declErrs, common.ActionInvalidError("action is 'nil': Cleanup"))
	}
	app.cleanup = action
	return app
}

func (app *workflow) GetCleanup() func(ctx common.Runtime, err error) {
	return app.cleanup
}

func (app *workflow) OnError(action func(ctx common.Runtime, err error)) Workflow {
	if action == nil {
		app.declErrs = append(app.declErrs, common.ActionInvalidError("action is 'nil': OnError"))
	}
	app.onError = action
	return app
}

func (app *workflow) GetOnError() func(ctx common.Runtime, err error) {
	return app.onError
}

func (app *workflow) GlobalFlags(flags ...common.Flag) Workflow {
	app.declErrs = append(app.declErrs, common.ValidateFlagDeclarations(flags)...)
	app.flags = flags
	return app
}

func (app *workflow) GetGlobalFlags() []common.Flag {
	return app.flags
}

func (app *workflow) Commands(commands ...common.Declaration) Workflow {
	app.declErrs = append(app.declErrs, common.ValidateCommandDeclarations(commands)...)
	app.commands = commands
	return app
}

func (app *workflow) GetCommands() []common.Declaration {
	return app.commands
}

func (app *workflow) GetDeclarationErrors() []error {
	return app.declErrs
}
