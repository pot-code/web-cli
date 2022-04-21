package command

import "reflect"

func getReflectedStringSlice(f reflect.Value) []string {
	l := f.Len()
	s := make([]string, l)
	for i := l - 1; i >= 0; i-- {
		s[i] = f.Index(i).String()
	}
	return s
}
