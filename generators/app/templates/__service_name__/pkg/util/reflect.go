<%=licenseText%>
package util

import (
	"reflect"
)

// GetMethods uses the reflect package to get the method names on defined on
// in.
func GetMethods(in interface{}) []string {
	if in == nil {
		return []string{}
	}

	t := reflect.TypeOf(in)
	if t.Kind() != reflect.Ptr {
		t = reflect.PtrTo(t)
	}

	numMethods := t.NumMethod()
	methods := make([]string, numMethods)
	for i := 0; i < numMethods; i++ {
		methods[i] = t.Method(i).Name
	}

	return methods
}
