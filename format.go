package expect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
)

func f(i interface{}) (string, string) {
	if i == nil {
		return "nil", " "
	}
	switch i.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return fmt.Sprintf("%v", i), " "
	}

	kind := reflect.TypeOf(i).Kind()

	switch kind {
	case reflect.Ptr:
		// de-reference ptr and call formater again
		return f(reflect.ValueOf(i).Elem().Interface())

	case reflect.Array, reflect.Map, reflect.Slice, reflect.Struct:
		y, err := yaml.Marshal(i)
		if y[len(y)-1] == '\n' {
			y = y[:len(y)-1]
		}
		lines := []string{}
		for _, line := range strings.Split(string(y), "\n") {
			lines = append(lines, "  > "+line)
		}

		if err == nil {
			return strings.Join(lines, "\n"), "\n"
		}
	}

	return fmt.Sprintf("'%v'", i), " "
}
