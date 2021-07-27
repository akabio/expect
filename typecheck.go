package expect

import "reflect"

func sameType(a, b interface{}) bool {
	at := reflect.TypeOf(a)
	bt := reflect.TypeOf(b)
	return at == bt
}

func typeName(a interface{}) string {
	if a == nil {
		return "<nil>"
	}
	return reflect.TypeOf(a).String()
}
