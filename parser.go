package stalk

import (
	"strings"

	"github.com/pavelmemory/stalk/command"
	"github.com/pavelmemory/stalk/common"
	"github.com/pavelmemory/stalk/context"
)

func parse(workflow Workflow, args []string) (common.Runtime, error) {
	nextStart, parsedGlobalFlags, err := parseFlags(workflow.GetGlobalFlags(), args, 0)
	if err != nil {
		return nil, err
	}

	argsStart, parsedCommand, err := parseCommands(workflow.GetCommands(), args, nextStart)
	if err != nil {
		return nil, err
	}

	runCtx := context.NewRuntimeContext(parsedGlobalFlags, parsedCommand, args[argsStart:])
	return runCtx, nil
}

func parseCommands(declaredCommands []common.Declaration, parts []string, start int) (int, common.Parsed, error) {
	if start >= len(parts) {
		return start, nil, nil
	}

	expectedCommandDeclarationsByName := make(map[string]common.Declaration)
	for _, declaredCommand := range declaredCommands {
		expectedCommandDeclarationsByName[declaredCommand.GetName()] = declaredCommand
	}
	if foundCommandDeclaration, found := expectedCommandDeclarationsByName[parts[start]]; found {
		parsedCommand := command.NewParsed(foundCommandDeclaration)
		nextStart, commandFlags, err := parseFlags(foundCommandDeclaration.GetFlags(), parts, start+1)
		if err != nil {
			return 0, nil, err
		}

		parsedCommand.FoundFlags(commandFlags)

		if len(foundCommandDeclaration.GetSubCommands()) > 0 {
			var subCmd common.Parsed
			nextStart, subCmd, err = parseCommands(foundCommandDeclaration.GetSubCommands(), parts, nextStart)
			if err != nil {
				return 0, nil, err
			}
			parsedCommand.SubCommand(subCmd)
		}
		return nextStart, parsedCommand, nil
	}
	return start, nil, common.NotImplementedError("command: '" + parts[start] + "'")
}

func parseFlags(expectedFlags []common.Flag, parts []string, start int) (int, []common.Flag, error) {
	expectedFlagsByName := make(map[string]common.Flag)
	expectedFlagsByShortcut := make(map[rune]common.Flag)
	requiredFlagsByName := make(map[string]common.Flag)
	for _, flag := range expectedFlags {
		expectedFlagsByName[flag.GetName()] = flag
		if flag.GetShortcut() != common.ShortcutNotProvided {
			expectedFlagsByShortcut[flag.GetShortcut()] = flag
		}
		if flag.IsRequired() {
			requiredFlagsByName[flag.GetName()] = flag
		}
	}

	var (
		foundFlags      []common.Flag
		lastParsedIndex int
	)

	for lastParsedIndex = start; lastParsedIndex < len(parts); lastParsedIndex++ {
		part := parts[lastParsedIndex]
		if !strings.HasPrefix(part, "-") {
			break
		}

		f, err := getFlag(part, expectedFlagsByName, expectedFlagsByShortcut)
		if err != nil {
			return 0, nil, err
		}
		if f == nil {
			return lastParsedIndex, nil, nil
		}

		if f.IsRequired() {
			delete(requiredFlagsByName, f.GetName())
		}

		delete(expectedFlagsByName, f.GetName())
		delete(expectedFlagsByShortcut, f.GetShortcut())
		if !f.IsSignal() {
			if lastParsedIndex+1 >= len(parts) {
				return 0, nil, common.NotAllRequiredValuesError(f.String())
			}
			lastParsedIndex++
			if err := f.Proceed(parts[lastParsedIndex]); err != nil {
				return 0, nil, err
			}
		}
		foundFlags = append(foundFlags, f)

	}
	if len(requiredFlagsByName) != 0 {
		var flagStrings []string
		for _, requiredFlag := range requiredFlagsByName {
			flagStrings = append(flagStrings, requiredFlag.String())
		}
		return 0, nil, common.NotAllRequiredFlagsError(strings.Join(flagStrings, "\n"))
	}
	for _, f := range expectedFlagsByName {
		if f.HasDefault() {
			foundFlags = append(foundFlags, f)
		}
	}
	return lastParsedIndex, foundFlags, nil
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
