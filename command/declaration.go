package command

import (
	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.Declaration = (*declaration)(nil)
)

func New(name string) common.Declaration {
	return &declaration{name: name}
}

func (c *declaration) GetName() string {
	return c.name
}

func (c *declaration) Flags(flags ...common.Flag) common.Declaration {
	c.declaredFlags = flags
	return c
}

func (c *declaration) GetFlags() []common.Flag {
	return c.declaredFlags
}

func (c *declaration) SubCommands(commands ...common.Declaration) common.Declaration {
	c.declaredSubCommands = commands
	return c
}

func (c *declaration) GetSubCommands() []common.Declaration {
	return c.declaredSubCommands
}

func (c *declaration) Execute(action func(ctx common.Runtime) error) common.Declaration {
	c.action = action
	return c
}

func (c *declaration) GetExecution() func(ctx common.Runtime) error {
	return c.action
}

func (c *declaration) Before(action func(ctx common.Runtime) error) common.Declaration {
	c.before = action
	return c
}

func (c *declaration) GetBefore() func(ctx common.Runtime) error {
	return c.before
}

func (c *declaration) After(action func(ctx common.Runtime, err error)) common.Declaration {
	c.after = action
	return c
}

func (c *declaration) GetAfter() func(ctx common.Runtime, err error) {
	return c.after
}

func (c *declaration) Stringer(stringer func(declaration common.Declaration) string) common.Declaration {
	c.stringer = stringer
	return c
}

func (c *declaration) GetStringer() func(declaration common.Declaration) string {
	return c.stringer
}

func (c *declaration) String() string {
	if c.stringer == nil {
		return common.DefaultCommandStringer(c)
	}
	return c.stringer(c)
}

type declaration struct {
	name                string
	declaredFlags       []common.Flag
	declaredSubCommands []common.Declaration
	action              func(ctx common.Runtime) error
	before              func(ctx common.Runtime) error
	after               func(ctx common.Runtime, err error)
	stringer            func(declaration common.Declaration) string
}
