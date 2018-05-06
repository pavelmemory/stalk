package command

import (
	"strings"

	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.CommandDeclaration = (*declaration)(nil)

	DefaultCommandStringer = func(declaration common.CommandDeclaration) string {
		name := declaration.GetName()
		flags := ""
		if len(declaration.GetDeclaredFlags()) != 0 {
			var flgs []string
			for _, flg := range declaration.GetDeclaredFlags() {
				flgs = append(flgs, flg.String())
			}
			flags = " " + strings.Join(flgs, " ")
		}
		subcommands := ""
		if len(declaration.GetDeclaredSubCommands()) != 0 {
			var subcmds []string
			for _, subcmd := range declaration.GetDeclaredSubCommands() {
				subcmds = append(subcmds, subcmd.GetName())
			}
			subcommands = "[" + strings.Join(subcmds, "|") + "]"
		}
		return name + flags + subcommands
	}
)

func New(name string) common.CommandDeclaration {
	name = strings.TrimSpace(name)
	decl := &declaration{name: name}
	if name == "" {
		decl.declErrs = append(decl.declErrs, common.CommandNameInvalidError(common.EmptyNameMessage))
	}
	if len(strings.Fields(name)) > 1 {
		decl.declErrs = append(decl.declErrs, common.CommandNameInvalidError(name))
	}

	return decl
}

type declaration struct {
	name                string
	declaredFlags       []common.Flag
	declaredSubCommands []common.CommandDeclaration
	action              func(ctx common.Runtime) error
	before              func(ctx common.Runtime) error
	after               func(ctx common.Runtime, err error)
	onError             func(ctx common.Runtime, err error)
	stringer            func(declaration common.CommandDeclaration) string
	description         string
	declErrs            []error
}

func (c *declaration) GetName() string {
	return c.name
}

func (c *declaration) WithFlags(flags ...common.Flag) common.CommandDeclaration {
	c.declErrs = append(c.declErrs, common.ValidateFlagDeclarations(flags)...)
	c.declaredFlags = flags
	return c
}

func (c *declaration) GetDeclaredFlags() []common.Flag {
	return c.declaredFlags
}

func (c *declaration) WithSubCommands(commands ...common.CommandDeclaration) common.CommandDeclaration {
	c.declErrs = append(c.declErrs, common.ValidateCommandDeclarations(commands)...)
	c.declaredSubCommands = commands
	return c
}

func (c *declaration) GetDeclaredSubCommands() []common.CommandDeclaration {
	return c.declaredSubCommands
}

func (c *declaration) WithAction(action func(ctx common.Runtime) error) common.CommandDeclaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': Execute"))
	}
	c.action = action
	return c
}

func (c *declaration) GetDeclaredAction() func(ctx common.Runtime) error {
	return c.action
}

func (c *declaration) WithBefore(action func(ctx common.Runtime) error) common.CommandDeclaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': Before"))
	}
	c.before = action
	return c
}

func (c *declaration) GetDeclaredBefore() func(ctx common.Runtime) error {
	return c.before
}

func (c *declaration) WithAfter(action func(ctx common.Runtime, err error)) common.CommandDeclaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': After"))
	}
	c.after = action
	return c
}

func (c *declaration) GetDeclaredAfter() func(ctx common.Runtime, err error) {
	return c.after
}

func (c *declaration) WithOnError(action func(ctx common.Runtime, err error)) common.CommandDeclaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': OnError"))
	}
	c.onError = action
	return c
}

func (c *declaration) GetDeclaredOnError() func(ctx common.Runtime, err error) {
	return c.onError
}

func (c *declaration) WithStringer(stringer func(declaration common.CommandDeclaration) string) common.CommandDeclaration {
	c.stringer = stringer
	return c
}

func (c *declaration) GetDeclaredStringer() func(declaration common.CommandDeclaration) string {
	return c.stringer
}

func (c *declaration) String() string {
	if c.stringer == nil {
		return DefaultCommandStringer(c)
	}
	return c.stringer(c)
}

func (c *declaration) WithDescription(value string) common.CommandDeclaration {
	c.description = value
	return c
}

func (c *declaration) GetDeclaredDescription() string {
	return c.description
}

func (c *declaration) GetDeclarationErrors() []error {
	return c.declErrs
}
