package flag

import (
	"strconv"

	"github.com/pavelmemory/stalk/common"
)

var (
	_ common.Flag         = (*impl)(nil)
	_ common.ParsedString = (*impl)(nil)
	_ common.ParsedInt    = (*impl)(nil)
	_ common.ParsedBool   = (*impl)(nil)
	_ common.ParsedFloat  = (*impl)(nil)
)

type impl struct {
	name       string
	shortcut   string
	required   bool
	proceed    func(value string) error
	value      interface{}
	signal     bool
	hasDefault bool
	stringer   func(flag common.Flag) string
}

func (f *impl) Name(value string) common.Flag {
	f.name = value
	return f
}

func (f *impl) GetName() string {
	return f.name
}

func (f *impl) Shortcut(value string) common.Flag {
	f.shortcut = value
	return f
}

func (f *impl) GetShortcut() string {
	return f.shortcut
}

func (f *impl) Required(value bool) common.Flag {
	f.required = value
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
	f.stringer = stringer
	return f
}

func (f *impl) GetStringer() func(flag common.Flag) string {
	return f.stringer
}

func (f *impl) String() string {
	if f.stringer == nil {
		return common.DefaultFlagStringer(f)
	}
	return f.stringer(f)
}

func Int(name string) common.Flag {
	f := &impl{name: name}
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
	f := &impl{name: name}
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
	return &impl{
		name:   name,
		signal: true,
		proceed: func(value string) error {
			return common.ErrorNotImplemented
		},
	}
}

func SignalSetByDefaultTo(name string) common.Flag {
	return setDefault(Signal(name), nil)
}

func Float(name string) common.Flag {
	f := &impl{
		name: name,
	}
	f.proceed = func(value string) (err error) {
		f.value, err = strconv.ParseFloat(value, 64)
		return
	}
	return f
}

func FloatWithDefault(name string, value float64) common.Flag {
	return setDefault(Float(name), value)
}

func setDefault(f common.Flag, value interface{}) common.Flag {
	fi := f.(*impl)
	fi.hasDefault = true
	fi.value = value
	return f
}
