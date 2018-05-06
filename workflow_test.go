package stalk

import (
	"fmt"
	"testing"

	"errors"
	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
)

func TestWorkflow_Run(t *testing.T) {
	createCmd := command.New("create").
		WithFlags(flag.String("name").Required(true).WithShortcut('n')).
		WithBefore(func(ctx common.Runtime) error {
			fmt.Println("before create command")
			return nil
		}).
		WithAction(func(ctx common.Runtime) error {
			fmt.Println("create command")
			fmt.Println("ctx.HasGlobalFlag(\"verbose\")", ctx.HasGlobalFlag("verbose"))
			fmt.Println("ctx.StringFlag(\"name\")", ctx.StringFlag("name"))
			fmt.Println("ctx.HasFlag(\"name\")", ctx.HasFlag("name"))
			fmt.Println("ctx.GetArgs()", ctx.GetArgs())
			fmt.Println("ctx.GlobalFlags()", ctx.GlobalFlags())
			fmt.Println("ctx.Flags()", ctx.Flags())
			v, f := ctx.Get("1")
			fmt.Println("ctx.Get(\"1\")", v, f)
			return nil
		})

	showCmd := command.New("show").
		WithFlags(flag.String("name").Required(true).WithShortcut('n')).
		WithAction(func(ctx common.Runtime) error {
			fmt.Println("show command")
			fmt.Println(ctx.StringGlobalFlag("example"))
			fmt.Println(ctx.StringGlobalFlag("help"))
			fmt.Println(ctx.StringFlag("name"))
			return nil
		})

	aws := command.New("aws").
		WithSubCommands(createCmd, showCmd).
		WithAction(func(ctx common.Runtime) error {
			ctx.Set("1", "something")
			return nil
		})

	app := New().
		WithCommands(aws).
		WithGlobalFlags(
			flag.String("example").WithShortcut('e'),
			flag.Signal("verbose").WithShortcut('v'),
			flag.Signal("help").WithShortcut('h'))
	err := app.Run([]string{"--verbose", "--example", "don'tbrelieve", "aws", "create", "--name", "tattoo", "valhalla", "and", "dumb"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestWorkflow_CurrentCommand(t *testing.T) {
	commandName := "c"
	childCommandName := "cc"
	emptyError := errors.New("")

	expectedChecks := 0
	checkedTimes := 0

	doCheckOnError := func(expectedCmdName string) func(common.Runtime, error) {
		expectedChecks++
		return func(ctx common.Runtime, err error) {
			cmdName := ctx.CurrentCommand().GetName()
			if cmdName != expectedCmdName {
				t.Errorf("%d: | expected: %s | got: %s", checkedTimes, expectedCmdName, cmdName)
			}
			checkedTimes++
		}
	}

	doCheckAndReturn := func(cmdName string, err error) func(ctx common.Runtime) error {
		return func(ctx common.Runtime) error {
			doCheckOnError(cmdName)(ctx, err)
			return err
		}
	}

	err := New().
		WithCommands(
			command.New(commandName).
				WithAction(doCheckAndReturn(commandName, nil)).
				WithOnError(doCheckOnError(commandName)).
				WithBefore(doCheckAndReturn(commandName, nil)).
				WithAfter(doCheckOnError(commandName)).
				WithSubCommands(
					command.New(childCommandName).
						WithAction(doCheckAndReturn(childCommandName, emptyError)).
						WithAfter(doCheckOnError(childCommandName)).
						WithOnError(doCheckOnError(childCommandName))),
		).
		WithOnError(doCheckOnError(commandName)).
		Run([]string{commandName, childCommandName})
	if err == nil {
		t.Error("error expected")
	}
	if checkedTimes > expectedChecks {
		t.Error("too many checks triggered")
	}
	if checkedTimes < expectedChecks {
		t.Error("not all checks triggered")
	}
}
