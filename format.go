package expect

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/ghodss/yaml"
)

func f(i interface{}) (string, string) {
	switch i.(type) {
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64, float32, float64:
		return fmt.Sprintf("%v", i), " "
	}

	switch reflect.TypeOf(i).Kind() {
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
