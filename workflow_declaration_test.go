package stalk

import (
	"testing"

	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/flag"
)

func TestWorkflow_GetDeclarationErrors_SetupInvalid(t *testing.T) {
	t.Parallel()
	wf := New().Setup(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_CleanupInvalid(t *testing.T) {
	t.Parallel()
	wf := New().Cleanup(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_OnErrorInvalid(t *testing.T) {
	t.Parallel()
	wf := New().OnError(nil)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorActionInvalid)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().GlobalFlags(flag.String("invalid name"))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().GlobalFlags(
		flag.String("same_name"),
		flag.String("same_name"),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_GlobalFlagShortcutDuplication(t *testing.T) {
	t.Parallel()
	wf := New().GlobalFlags(
		flag.String("name1").Shortcut('s'),
		flag.String("name2").Shortcut('s'),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagShortcutNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().Commands(command.New("invalid name"))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorCommandNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().Commands(
		command.New("same_name"),
		command.New("same_name"),
	)
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorCommandNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagNameInvalid(t *testing.T) {
	t.Parallel()
	wf := New().Commands(command.New("cmd").Flags(flag.String("")))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameInvalid)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagNameDuplication(t *testing.T) {
	t.Parallel()
	wf := New().Commands(command.New("cmd").Flags(
		flag.String("same_name"),
		flag.String("same_name"),
	))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagNameNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandFlagShortcutDuplication(t *testing.T) {
	t.Parallel()
	wf := New().Commands(command.New("cmd").Flags(
		flag.String("name1").Shortcut('s'),
		flag.String("name2").Shortcut('s'),
	))
	actual := wf.GetDeclarationErrors()
	assertErrorCode(t, actual, common.ErrorFlagShortcutNotUnique)
}

func TestWorkflow_GetDeclarationErrors_CommandExecuteInvalid(t *testing.T) {
	t.Parallel()
	wf := New().Commands(command.New("cmd").Execute(nil))
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
