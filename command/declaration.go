package command

import (
	"strings"

	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.Declaration = (*declaration)(nil)

	DefaultCommandStringerProvider = func(declaration common.Declaration) string {
		name := declaration.GetName()
		flags := ""
		if len(declaration.GetFlags()) != 0 {
			var flgs []string
			for _, flg := range declaration.GetFlags() {
				flgs = append(flgs, flg.String())
			}
			flags = " " + strings.Join(flgs, " ")
		}
		subcommands := ""
		if len(declaration.GetSubCommands()) != 0 {
			var subcmds []string
			for _, subcmd := range declaration.GetSubCommands() {
				subcmds = append(subcmds, subcmd.GetName())
			}
			subcommands = "[" + strings.Join(subcmds, "|") + "]"
		}
		return name + flags + subcommands
	}
)

func New(name string) common.Declaration {
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
	declaredSubCommands []common.Declaration
	action              func(ctx common.Runtime) error
	before              func(ctx common.Runtime) error
	after               func(ctx common.Runtime, err error)
	onError             func(ctx common.Runtime, err error)
	stringerProvider    func(declaration common.Declaration) string
	description         string
	declErrs            []error
}

func (c *declaration) GetName() string {
	return c.name
}

func (c *declaration) Flags(flags ...common.Flag) common.Declaration {
	c.declErrs = append(c.declErrs, common.ValidateFlagDeclarations(flags)...)
	c.declaredFlags = flags
	return c
}

func (c *declaration) GetFlags() []common.Flag {
	return c.declaredFlags
}

func (c *declaration) SubCommands(commands ...common.Declaration) common.Declaration {
	c.declErrs = append(c.declErrs, common.ValidateCommandDeclarations(commands)...)
	c.declaredSubCommands = commands
	return c
}

func (c *declaration) GetSubCommands() []common.Declaration {
	return c.declaredSubCommands
}

func (c *declaration) Execute(action func(ctx common.Runtime) error) common.Declaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': Execute"))
	}
	c.action = action
	return c
}

func (c *declaration) GetExecution() func(ctx common.Runtime) error {
	return c.action
}

func (c *declaration) Before(action func(ctx common.Runtime) error) common.Declaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': Before"))
	}
	c.before = action
	return c
}

func (c *declaration) GetBefore() func(ctx common.Runtime) error {
	return c.before
}

func (c *declaration) After(action func(ctx common.Runtime, err error)) common.Declaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': After"))
	}
	c.after = action
	return c
}

func (c *declaration) GetAfter() func(ctx common.Runtime, err error) {
	return c.after
}

func (c *declaration) OnError(action func(ctx common.Runtime, err error)) common.Declaration {
	if action == nil {
		c.declErrs = append(c.declErrs, common.ActionInvalidError("action is 'nil': OnError"))
	}
	c.onError = action
	return c
}

func (c *declaration) GetOnError() func(ctx common.Runtime, err error) {
	return c.onError
}

func (c *declaration) StringerProvider(stringer func(declaration common.Declaration) string) common.Declaration {
	c.stringerProvider = stringer
	return c
}

func (c *declaration) GetStringerProvider() func(declaration common.Declaration) string {
	return c.stringerProvider
}

func (c *declaration) String() string {
	if c.stringerProvider == nil {
		return DefaultCommandStringerProvider(c)
	}
	return c.stringerProvider(c)
}

func (c *declaration) Description(value string) common.Declaration {
	c.description = value
	return c
}

func (c *declaration) GetDescription() string {
	return c.description
}

func (c *declaration) GetDeclarationErrors() []error {
	return c.declErrs
}
