package stalk

import (
	"testing"

	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
)

var emptyAction = func(runtime common.Runtime) error {
	return nil
}

func TestWorkflow_GetDeclarationErrors_SetupInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithSetup(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CleanupInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCleanup(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_OnErrorInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithOnError(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithGlobalFlags(flag.String("invalid name"))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().WithGlobalFlags(
		flag.String("same_name"),
		flag.String("same_name"),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagShortcutDuplication(t *testing.T) {
	t.Parallel()
	wf := New().WithGlobalFlags(
		flag.String("name1").WithShortcut('s'),
		flag.String("name2").WithShortcut('s'),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagShortcutNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("invalid name").WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorCommandNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(
		command.New("same_name").WithAction(emptyAction),
		command.New("same_name").WithAction(emptyAction),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorCommandNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithFlags(flag.String("")).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithFlags(
		flag.String("same_name"),
		flag.String("same_name"),
	).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagShortcutDuplication(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithFlags(
		flag.String("name1").WithShortcut('s'),
		flag.String("name2").WithShortcut('s'),
	).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagShortcutNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandExecuteInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithAction(nil).WithSubCommands(command.New("c").WithAction(emptyAction)))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandBeforeInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithBefore(nil).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandAfterInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithAfter(nil).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandOnErrorInvalid(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd").WithOnError(nil).WithAction(emptyAction))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandNoExecuteNoSubCommands(t *testing.T) {
	t.Parallel()
	wf := New().WithCommands(command.New("cmd"))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func assertErrorCode(t *testing.T, actual []error, expected common.ErrorCode) {
	switch len(actual) {
	case 0:
		t.Fatal("error expected")
	case 1:
		if cErr, ok := actual[0].(common.Error); ok {
			if cErr.Cause == expected {
				return
			}
		}
		t.Error("\nexpected:\n", expected, "\nactual:\n", actual)
	default:
		t.Fatal("unexpected", actual)
	}
}
