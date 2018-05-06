package command

import (
	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.ParsedCommand = (*parsed)(nil)
)

func NewParsed(declaration common.CommandDeclaration) common.ParsedCommand {
	return &parsed{
		CommandDeclaration: declaration,
		foundFlags:  make(map[string]common.Flag),
	}
}

type parsed struct {
	common.CommandDeclaration
	foundFlags      map[string]common.Flag
	foundSubCommand common.ParsedCommand
}

func (c *parsed) SubCommand(command common.ParsedCommand) {
	c.foundSubCommand = command
}

func (c *parsed) GetSubCommand() common.ParsedCommand {
	return c.foundSubCommand
}

func (c *parsed) Flags(flags []common.Flag) {
	for _, f := range flags {
		c.foundFlags[f.GetName()] = f
	}
}

func (c *parsed) GetFlags() map[string]common.Flag {
	return c.foundFlags
}
