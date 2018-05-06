package stalk

import (
	"strings"

	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/context"
)

func parse(workflow Workflow, args []string) (common.Runtime, error) {
	nextStart, parsedGlobalFlags, err := parseFlags(workflow.GetDeclaredGlobalFlags(), args, 0)
	if err != nil {
		return nil, err
	}

	argsStart, parsedCommand, err := parseCommands(workflow.GetDeclaredCommands(), args, nextStart)
	if err != nil {
		return nil, err
	}

	runCtx := context.NewRuntimeContext(parsedGlobalFlags, parsedCommand, args[argsStart:])
	return runCtx, nil
}

func parseCommands(declaredCommands []common.CommandDeclaration, parts []string, start int) (int, common.ParsedCommand, error) {
	if start >= len(parts) {
		return start, nil, nil
	}

	expectedCommandDeclarationsByName := make(map[string]common.CommandDeclaration)
	for _, declaredCommand := range declaredCommands {
		expectedCommandDeclarationsByName[declaredCommand.GetName()] = declaredCommand
	}
	if foundCommandDeclaration, found := expectedCommandDeclarationsByName[parts[start]]; found {
		parsedCommand := command.NewParsed(foundCommandDeclaration)
		nextStart, commandFlags, err := parseFlags(foundCommandDeclaration.GetDeclaredFlags(), parts, start+1)
		if err != nil {
			return 0, nil, err
		}

		parsedCommand.Flags(commandFlags)

		if len(foundCommandDeclaration.GetDeclaredSubCommands()) > 0 {
			var subCmd common.ParsedCommand
			nextStart, subCmd, err = parseCommands(foundCommandDeclaration.GetDeclaredSubCommands(), parts, nextStart)
			if err != nil {
				return 0, nil, err
			}
			parsedCommand.SubCommand(subCmd)
		}
		return nextStart, parsedCommand, nil
	}
	return start, nil, common.NotImplementedError("command: '" + parts[start] + "'")
}

func parseFlags(expectedFlags []common.Flag, rawInput []string, start int) (lastParsedIndex int, foundFlags []common.Flag, err error) {
	expectedFlagsByName := make(map[string]common.Flag)
	expectedFlagsByShortcut := make(map[rune]common.Flag)
	requiredFlagsByName := make(map[string]common.Flag)
	for _, flag := range expectedFlags {
		expectedFlagsByName[flag.GetName()] = flag
		if flag.GetShortcut() != common.ShortcutNotProvided {
			expectedFlagsByShortcut[flag.GetShortcut()] = flag
		}
		if flag.IsDeclaredRequired() {
			requiredFlagsByName[flag.GetName()] = flag
		}
	}

	for lastParsedIndex = start; lastParsedIndex < len(rawInput); lastParsedIndex++ {
		part := rawInput[lastParsedIndex]
		if !strings.HasPrefix(part, "-") {
			break
		}

		var flag common.Flag
		flag, err = getFlag(part, expectedFlagsByName, expectedFlagsByShortcut)
		if err != nil || flag == nil {
			return
		}

		if flag.IsDeclaredRequired() {
			delete(requiredFlagsByName, flag.GetName())
		}

		delete(expectedFlagsByName, flag.GetName())
		delete(expectedFlagsByShortcut, flag.GetShortcut())
		if !flag.IsDeclaredSignal() {
			if lastParsedIndex+1 >= len(rawInput) {
				return 0, nil, common.NotAllRequiredValuesError(flag.String())
			}
			lastParsedIndex++
			if err := flag.Parse(rawInput[lastParsedIndex]); err != nil {
				return 0, nil, err
			}
		}
		foundFlags = append(foundFlags, flag)
	}

	if len(requiredFlagsByName) != 0 {
		var flagStrings []string
		for _, requiredFlag := range requiredFlagsByName {
			flagStrings = append(flagStrings, requiredFlag.String())
		}
		return 0, nil, common.NotAllRequiredFlagsError(strings.Join(flagStrings, "\n"))
	}

	for _, flag := range expectedFlagsByName {
		if flag.HasDefault() {
			foundFlags = append(foundFlags, flag)
		}
	}
	return
}

func getFlag(part string, expectedFlagsByName map[string]common.Flag, expectedFlagsByShortcut map[rune]common.Flag) (common.Flag, error) {
	switch {
	case strings.HasPrefix(part, "--"):
		flagName := part[2:]
		if flagName == "" {
			return nil, common.FlagSyntaxError(part)
		}
		if f, found := expectedFlagsByName[flagName]; found {
			return f, nil
		}
		return nil, common.FlagNotSupportedError(part)
	case strings.HasPrefix(part, "-"):
		flagName := part[1:]
		if flagName == "" {
			return nil, common.FlagSyntaxError(part)
		}
		if f, found := expectedFlagsByShortcut[rune(flagName[0])]; found {
			return f, nil
		}
		return nil, common.FlagNotSupportedError(part)
	default:
		return nil, nil
	}
}
