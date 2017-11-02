package flag

import (
	"strconv"

	"fmt"
	"github.com/pavelmemory/stalk/common"
	"strings"
	"unicode/utf8"
)

var (
	emptyRune rune

	DefaultFlagStringerProvider = func(flag common.Flag) string {
		name := "[--" + flag.GetName()
		shortcut := "]"
		if flag.GetShortcut() != emptyRune {
			shortcut = "|-" + string(flag.GetShortcut()) + "]"
		}
		if flag.IsSignal() {
			return name + shortcut
		}

		required := ""
		if flag.IsRequired() {
			required = "+"
		}

		fimpl := flag.(*impl)
		if flag.HasDefault() {
			return name + shortcut + required + " <" + fimpl.valueTypeName + ", " + fmt.Sprint(fimpl.value) + ">"
		}
		return name + shortcut + required + " [" + fimpl.valueTypeName + "]"
	}

	DefaultFlagDescriptionProvider = func(flag common.Flag) string {
		return flag.GetDescription()
	}

	_ common.Flag         = (*impl)(nil)
	_ common.ParsedString = (*impl)(nil)
	_ common.ParsedInt    = (*impl)(nil)
	_ common.ParsedBool   = (*impl)(nil)
	_ common.ParsedFloat  = (*impl)(nil)
)

type impl struct {
	name          string
	shortcut      rune
	required      bool
	proceed       func(value string) error
	value         interface{}
	valueTypeName string
	signal        bool
	hasDefault    bool
	stringerProv  func(flag common.Flag) string
	descProv      func(flag common.Flag) string
	description   string
	declErrs      []error
}

func (f *impl) Name(value string) common.Flag {
	value = strings.TrimSpace(value)
	if value == "" {
		f.declErrs = append(f.declErrs, common.FlagNameInvalidError(common.EmptyNameMessage))
	}
	if len(strings.Fields(value)) > 1 {
		f.declErrs = append(f.declErrs, common.FlagNameInvalidError(f.String()))
	}

	f.name = value
	return f
}

func (f *impl) GetName() string {
	return f.name
}

func (f *impl) Shortcut(value rune) common.Flag {
	if !utf8.ValidRune(value) || value == common.ShortcutNotProvided {
		f.declErrs = append(f.declErrs, common.FlagShortcutInvalidError(f.String()))
	}

	f.shortcut = value
	return f
}

func (f *impl) GetShortcut() rune {
	return f.shortcut
}

func (f *impl) Required(value bool) common.Flag {
	f.required = value
	if f.HasDefault() {
		f.declErrs = append(f.declErrs, common.FlagRequiredAndHasDefaultError(f.GetName()))
	}
	if f.IsSignal() {
		f.declErrs = append(f.declErrs, common.FlagSignalAndRequiredError(f.GetName()))
	}
	return f
}

func (f *impl) IsRequired() bool {
	return f.required
}

func (f *impl) HasDefault() bool {
	return f.hasDefault
}

func (f *impl) Proceed(value string) error {
	return f.proceed(value)
}

func (f *impl) IsSignal() bool {
	return f.signal
}

func (f *impl) StringValue() string {
	return f.value.(string)
}

func (f *impl) IntValue() int64 {
	return f.value.(int64)
}

func (f *impl) BoolValue() bool {
	return f.value.(bool)
}

func (f *impl) FloatValue() float64 {
	return f.value.(float64)
}

func (f *impl) Stringer(stringer func(flag common.Flag) string) common.Flag {
	f.stringerProv = stringer
	return f
}

func (f *impl) GetStringer() func(flag common.Flag) string {
	if f.stringerProv == nil {
		return DefaultFlagStringerProvider
	}
	return f.stringerProv
}

func (f *impl) String() string {
	return f.GetStringer()(f)
}

func (f *impl) DescriptionProvider(provider func(flag common.Flag) string) common.Flag {
	f.descProv = provider
	return f
}

func (f *impl) Description(value string) common.Flag {
	f.description = value
	return f
}

func (f *impl) GetDescription() string {
	return f.description
}

func (f *impl) GetDescriptionProvider() func(flag common.Flag) string {
	if f.descProv == nil {
		return DefaultFlagDescriptionProvider
	}
	return f.descProv
}

func (f *impl) GetDeclarationErrors() []error {
	return f.declErrs
}

func Int(name string) common.Flag {
	f := &impl{valueTypeName: "INT"}
	f.Name(name)
	f.proceed = func(value string) (err error) {
		f.value, err = strconv.ParseInt(value, 10, 64)
		return
	}
	return f
}

func IntWithDefault(name string, value int64) common.Flag {
	return setDefault(Int(name), value)
}

func String(name string) common.Flag {
	f := &impl{valueTypeName: "STRING"}
	f.Name(name)
	f.proceed = func(value string) error {
		f.value = value
		return nil
	}
	return f
}

func StringWithDefault(name, value string) common.Flag {
	return setDefault(String(name), value)
}

func Signal(name string) common.Flag {
	f := &impl{signal: true}
	f.Name(name)
	f.proceed = func(value string) error {
		return common.NotImplementedError("signal flag '" + name + "' doesn't expect any value")
	}
	return f
}

func Float(name string) common.Flag {
	f := &impl{valueTypeName: "FLOAT"}
	f.Name(name)
	f.proceed = func(value string) (err error) {
		f.value, err = strconv.ParseFloat(value, 64)
		return
	}
	return f
}

func FloatWithDefault(name string, value float64) common.Flag {
	return setDefault(Float(name), value)
}

func Bool(name string) common.Flag {
	f := &impl{valueTypeName: "BOOL"}
	f.Name(name)
	f.proceed = func(value string) (err error) {
		f.value, err = strconv.ParseBool(value)
		return
	}
	return f
}

func BoolWithDefault(name string, value bool) common.Flag {
	return setDefault(Bool(name), value)
}

func setDefault(f common.Flag, value interface{}) common.Flag {
	fi := f.(*impl)
	fi.hasDefault = true
	fi.value = value
	if fi.IsRequired() {
		fi.declErrs = append(fi.declErrs, common.FlagRequiredAndHasDefaultError(f.GetName()))
	}
	return f
}
