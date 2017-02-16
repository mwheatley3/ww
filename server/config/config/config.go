package config

import (
	"encoding"
	"encoding/json"
	"fmt"
	"github.com/serenize/snaker"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
)

var (
	unmarshaler = reflect.TypeOf((*encoding.TextUnmarshaler)(nil)).Elem()
	stringMap   = reflect.MapOf(reflect.TypeOf(""), reflect.TypeOf(""))
)

// A Provider generates config values as strings in the form of k=v
type Provider func() ([]string, error)

// Load decodes config values provided by providers
// into the struct pointed at by config
func Load(config interface{}, separator string, providers ...Provider) error {
	for _, p := range providers {
		d, err := p()

		if err != nil {
			return err
		}

		load(d, config, separator)
	}

	return nil
}

// FromFile provides config values from a file
func FromFile(path string) Provider {
	return func() ([]string, error) {
		b, err := ioutil.ReadFile(path)

		if err != nil {
			return nil, err
		}

		return strings.Split(string(b), "\n"), nil
	}
}

// FromEnv provides config values from the env
func FromEnv() Provider {
	return func() ([]string, error) {
		return os.Environ(), nil
	}
}

// FromStrings provides a static set of config values
func FromStrings(d []string) Provider {
	return func() ([]string, error) {
		return d, nil
	}
}

func load(data []string, config interface{}, separator string) {
	if separator == "" {
		separator = "__"
	}

	for _, s := range data {
		parts := strings.SplitN(s, "=", 2)
		if len(parts) != 2 {
			continue
		}

		keys := strings.Split(strings.TrimSpace(parts[0]), separator)
		v := reflect.ValueOf(config)
		ctxt := walk(v, keys)

		val := strings.TrimSpace(parts[1])

		if len(val) > 0 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
			val = val[1 : len(val)-1]
		}

		if ctxt.IsValid() {
			set(ctxt, val)
		}
	}
}

func set(ctxt reflect.Value, val string) {
	if ctxt.CanSet() {
		// value implements Settable
		if ctxt.Type().Implements(unmarshaler) {
			ctxt.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val))
			return
		}

		// ptr to value implements Settable
		if ctxt.CanAddr() {
			ptr := ctxt.Addr()

			if ptr.Type().Implements(unmarshaler) {
				ptr.Interface().(encoding.TextUnmarshaler).UnmarshalText([]byte(val))
				return
			}
		}

		k := ctxt.Kind()
		switch {
		case k == reflect.Bool:
			var b bool
			if val == "" {
				b = false
			} else {
				var err error
				b, err = strconv.ParseBool(val)
				if err != nil {
					b = false
				}
			}
			ctxt.SetBool(b)
		case k == reflect.String:
			ctxt.SetString(val)
		case k == reflect.Int || k == reflect.Int8 || k == reflect.Int16 || k == reflect.Int32 || k == reflect.Int64:
			if v, err := strconv.ParseInt(val, 10, 64); err == nil {
				ctxt.SetInt(v)
			}
		case k == reflect.Uint || k == reflect.Uint8 || k == reflect.Uint16 || k == reflect.Uint32 || k == reflect.Uint64:
			if v, err := strconv.ParseUint(val, 10, 64); err == nil {
				ctxt.SetUint(v)
			}
		case k == reflect.Float32 || k == reflect.Float64:
			if v, err := strconv.ParseFloat(val, 64); err == nil {
				ctxt.SetFloat(v)
			}
		case k == reflect.Ptr:
			set(ctxt.Elem(), val)
		case k == reflect.Map && ctxt.Type() == stringMap:
			m := make(map[string]string, 0)
			if err := json.Unmarshal([]byte(val), &m); err == nil {
				ctxt.Set(reflect.ValueOf(m))
			}
		}
	}
}

// String pretty prints a config struct
func String(conf interface{}) string {
	v := reflect.ValueOf(conf)
	str := str(v, 1)

	return "Config:\n" + str
}

func str(v reflect.Value, level int) string {
	v = deref(v)

	if v.Kind() == reflect.Struct {
		t := v.Type()
		strs := []string{}

		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			v2 := deref(v.Field(i))

			var str2 string

			if v2.Kind() == reflect.Struct {
				if f.Anonymous {
					str2 = str(v2, level)
				} else {
					str2 = "\n" + str(v2, level+1)
				}
			} else if v2.IsValid() {
				str2 = " " + fmt.Sprintf("%v", v2.Interface())
			}

			prefix := strings.Repeat("  ", level)
			name := f.Name + ":"

			if f.Anonymous {
				prefix = ""
				name = ""
			}

			strs = append(strs, prefix+name+str2)
		}

		return strings.Join(strs, "\n")
	}

	return ""
}

// walk the value and turn any nil
// ptrs into ptrs to a zero value of
// the correct type
func zero(v reflect.Value) {
	if v.Kind() == reflect.Ptr {
		if !v.IsNil() {
			v = v.Elem()
		} else {
			val := reflect.New(v.Type().Elem())
			v.Set(val)
			v = val
		}

		zero(v)
		return
	}

	if v.Kind() != reflect.Struct {
		return
	}

	num := v.NumField()

	for i := 0; i < num; i++ {
		f := v.Field(i)
		zero(f)
	}
}

func walk(v reflect.Value, keys []string) reflect.Value {
	zero(v)

	for _, k := range keys {
		v = deref(v).FieldByNameFunc(matchField(k))

		if !v.IsValid() {
			break
		}
	}

	return v
}

func deref(v reflect.Value) reflect.Value {
	if v.Kind() == reflect.Interface {
		return deref(v.Elem())
	}

	return reflect.Indirect(v)
}

func matchField(field string) func(string) bool {
	field = strings.ToLower(field)

	return func(name string) bool {
		return snaker.CamelToSnake(name) == field
	}
}
