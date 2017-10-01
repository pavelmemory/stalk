package command

import (
	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.Parsed = (*parsed)(nil)
)

func NewParsed(declaration common.Declaration) common.Parsed {
	return &parsed{
		Declaration: declaration,
		foundFlags:  make(map[string]common.Flag),
	}
}

type parsed struct {
	common.Declaration
	foundFlags      map[string]common.Flag
	foundSubCommand common.Parsed
}

func (c *parsed) SubCommand(command common.Parsed) {
	c.foundSubCommand = command
}

func (c *parsed) GetSubCommand() common.Parsed {
	return c.foundSubCommand
}

func (c *parsed) FoundFlags(flags []common.Flag) {
	for _, f := range flags {
		c.foundFlags[f.GetName()] = f
	}
}

func (c *parsed) GetFoundFlags() map[string]common.Flag {
	return c.foundFlags
}
