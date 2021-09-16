package expect

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/ghodss/yaml"
)

type presentation int

const (
	compact presentation = iota
	block
)

var presentations = map[presentation]string{
	compact: " ",
	block:   "\n",
}

func format(i interface{}) (string, presentation, bool) {
	if i == nil {
		return "nil", compact, false
	}
	// basic types
	switch t := i.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return fmt.Sprintf("%v", i), compact, false
	case string:
		return fmt.Sprintf("%v", i), compact, true
	case bool:
		return fmt.Sprintf("%v", i), compact, false
	case time.Time:
		return t.Format(time.RFC3339Nano), compact, false
	case error:
		return t.Error(), compact, true
	}

	kind := reflect.TypeOf(i).Kind()

	switch kind {
	case reflect.Ptr:
		// de-reference ptr and call formater again
		if isNil(i) {
			return "<nil>", compact, false
		}

		return format(reflect.ValueOf(i).Elem().Interface())

	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		y, err := yaml.Marshal(i)
		if err == nil {
			if len(y) > 0 && y[len(y)-1] == '\n' {
				y = y[:len(y)-1]
			}

			return string(y), block, false
		}

		return fmt.Sprintf("%v", i), compact, false
	}

	return fmt.Sprintf("%v", i), compact, true
}

func formatBoth(x interface{}, v interface{}) (string, string, presentation) {
	xf, xd, dx := format(x)
	vf, vd, dv := format(v)

	if vd == block || xd == block {
		return xf, vf, block
	}

	if strings.Contains(xf, "\n") || strings.Contains(vf, "\n") {
		return xf, vf, block
	}

	return del(xf, dx), del(vf, dv), compact
}

func formatOne(x interface{}) (string, presentation) {
	xf, xd, dx := format(x)
	if xd == block {
		return xf, block
	}

	if strings.Contains(xf, "\n") {
		return xf, block
	}

	return del(xf, dx), compact
}

func indent(v string, pres presentation) string {
	if pres == block {
		lines := strings.Split(v, "\n")
		for i := range lines {
			lines[i] = "    " + lines[i]
		}

		return strings.Join(lines, "\n")
	}

	return v
}

func del(v string, d bool) string {
	if d {
		return "'" + v + "'"
	}

	return v
}
