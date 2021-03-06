package context

import (
	"github.com/pavelmemory/stalk/common"
	"sync"
)

type runtimeContext struct {
	sync.RWMutex
	globalFlags    map[string]common.Flag
	rootCommand    common.ParsedCommand
	currentCommand common.ParsedCommand
	args           []string
	storage        map[interface{}]interface{}
}

func NewRuntimeContext(globalFlags []common.Flag, parsedCommand common.ParsedCommand, args []string) common.Runtime {
	rc := &runtimeContext{
		globalFlags: make(map[string]common.Flag),
		storage:     make(map[interface{}]interface{}),
		rootCommand: parsedCommand,
		args:        args,
	}
	for _, globalFlag := range globalFlags {
		rc.globalFlags[globalFlag.GetName()] = globalFlag
	}
	return rc
}

func (rc *runtimeContext) Run() error {
	return rc.runCommand(rc.rootCommand)
}

func (rc *runtimeContext) runCommand(parsedCommand common.ParsedCommand) (err error) {
	if parsedCommand == nil {
		return
	}
	rc.currentCommand = parsedCommand

	defer func() {
		rc.currentCommand = parsedCommand
		if err != nil {
			if onError := parsedCommand.GetDeclaredOnError(); onError != nil {
				onError(rc, err)
			}

			if afterAction := parsedCommand.GetDeclaredAfter(); afterAction != nil {
				afterAction(rc, err)
			}
		}
	}()

	if beforeAction := parsedCommand.GetDeclaredBefore(); beforeAction != nil {
		if err := beforeAction(rc); err != nil {
			return err
		}
	}

	if action := parsedCommand.GetDeclaredAction(); action != nil {
		err = action(rc)
	}

	if err == nil {
		parsedSubCommand := parsedCommand.GetSubCommand()
		err = rc.runCommand(parsedSubCommand)
	}

	return
}

func (rc *runtimeContext) GetArgs() []string {
	return rc.args
}

func (rc *runtimeContext) CurrentCommand() common.ParsedCommand {
	return rc.currentCommand
}

func (rc *runtimeContext) Flags() []string {
	var flagNames []string
	for flagName := range rc.currentCommand.GetFlags() {
		flagNames = append(flagNames, flagName)
	}
	return flagNames
}

func (rc *runtimeContext) GlobalFlags() []string {
	var flagNames []string
	for flagName := range rc.globalFlags {
		flagNames = append(flagNames, flagName)
	}
	return flagNames
}

func (rc *runtimeContext) Set(key interface{}, value interface{}) (oldValue interface{}, overridden bool) {
	rc.Lock()
	oldValue, overridden = rc.storage[key]
	rc.storage[key] = value
	rc.Unlock()
	return
}

func (rc *runtimeContext) Get(key interface{}) (value interface{}, found bool) {
	rc.RLock()
	value, found = rc.storage[key]
	rc.RUnlock()
	return
}

func (rc *runtimeContext) HasFlag(name string) (found bool) {
	_, found = rc.currentCommand.GetFlags()[name]
	return
}

func (rc *runtimeContext) HasGlobalFlag(name string) (found bool) {
	_, found = rc.globalFlags[name]
	return
}

func (rc *runtimeContext) StringFlag(name string) string {
	return stringFromMap(name, rc.currentCommand.GetFlags())
}

func (rc *runtimeContext) StringGlobalFlag(name string) string {
	return stringFromMap(name, rc.globalFlags)
}
func (rc *runtimeContext) IntFlag(name string) int64 {
	return intFromMap(name, rc.currentCommand.GetFlags())
}

func (rc *runtimeContext) IntGlobalFlag(name string) int64 {
	return intFromMap(name, rc.globalFlags)
}

func (rc *runtimeContext) FloatFlag(name string) float64 {
	return floatFromMap(name, rc.currentCommand.GetFlags())
}

func (rc *runtimeContext) FloatGlobalFlag(name string) float64 {
	return floatFromMap(name, rc.globalFlags)
}

func (rc *runtimeContext) BoolFlag(name string) bool {
	return boolFromMap(name, rc.currentCommand.GetFlags())
}

func (rc *runtimeContext) BoolGlobalFlag(name string) bool {
	return boolFromMap(name, rc.globalFlags)
}

func (rc *runtimeContext) CustomFlag(name string) interface{} {
	if f, found := rc.currentCommand.GetFlags()[name]; found {
		return f.(common.Custom).Value()
	}
	return nil
}

func (rc *runtimeContext) CustomGlobalFlag(name string) interface{} {
	if f, found := rc.globalFlags[name]; found {
		return f.(common.Custom).Value()
	}
	return nil
}

func stringFromMap(name string, m map[string]common.Flag) string {
	if f, found := m[name]; found {
		return f.(common.ParsedString).StringValue()
	}
	return ""
}

func intFromMap(name string, m map[string]common.Flag) int64 {
	if f, found := m[name]; found {
		return f.(common.ParsedInt).IntValue()
	}
	return 0
}

func boolFromMap(name string, m map[string]common.Flag) bool {
	if f, found := m[name]; found {
		return f.(common.ParsedBool).BoolValue()
	}
	return false
}

func floatFromMap(name string, m map[string]common.Flag) float64 {
	if f, found := m[name]; found {
		return f.(common.ParsedFloat).FloatValue()
	}
	return float64(0)
}
