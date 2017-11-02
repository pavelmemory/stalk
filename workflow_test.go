package stalk

import (
	"fmt"
	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
	"log"
	"testing"
)

func TestWorkflow_Run(t *testing.T) {
	createCmd := command.New("create")
	createCmd.Flags(flag.String("name").Required(true).Shortcut('n'))
	createCmd.Before(
		func(ctx common.Runtime) error {
			fmt.Println("before create command")
			return nil
		})
	createCmd.Execute(
		func(ctx common.Runtime) error {
			fmt.Println("create command")
			fmt.Println("ctx.HasGlobalFlag(\"verbose\")", ctx.HasGlobalFlag("verbose"))
			fmt.Println("ctx.StringFlag(\"name\")", ctx.StringFlag("name"))
			fmt.Println("ctx.HasFlag(\"name\")", ctx.HasFlag("name"))
			fmt.Println("ctx.GetArgs()", ctx.GetArgs())
			fmt.Println("ctx.FoundGlobalFlags()", ctx.FoundGlobalFlags())
			fmt.Println("ctx.FoundFlags()", ctx.FoundFlags())
			v, f := ctx.Get("1")
			fmt.Println("ctx.Get(\"1\")", v, f)
			return nil
		})

	showCmd := command.New("show")
	showCmd.Flags(flag.String("name").Required(true).Shortcut('n'))
	showCmd.Execute(
		func(ctx common.Runtime) error {
			fmt.Println("show command")
			fmt.Println(ctx.StringGlobalFlag("example"))
			fmt.Println(ctx.StringGlobalFlag("help"))
			fmt.Println(ctx.StringFlag("name"))
			return nil
		})

	aws := command.New("aws").SubCommands(createCmd, showCmd).Execute(func(ctx common.Runtime) error {
		ctx.Set("1", "something")
		return nil
	})

	app := New().Commands(aws).GlobalFlags(
		flag.String("example").Shortcut('e'),
		flag.Signal("verbose").Shortcut('v'),
		flag.Signal("help").Shortcut('h'),
	)
	err := app.Run([]string{"--verbose", "--example", "don'tbrelieve", "aws", "create", "--name", "tattoo", "valhalla", "and", "dumb"})
	if err != nil {
		log.Fatal(err)
	}
}
