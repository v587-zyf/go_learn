package test

import "strings"

func Split(str string, sep string) []string {
	var ret []string = make([]string, 0, strings.Count(str, sep)+1)
	index := strings.Index(str, sep)
	for index >= 0 {
		s := str[:index]
		if len(s) == 0 {
			continue
		}
		ret = append(ret, s)
		str = str[index+len(sep):]
		index = strings.Index(str, sep)
	}
	ret = append(ret, str)

	return ret
}

func Fib(n int) int {
	if n < 2 {
		return n
	}
	return Fib(n-1) * Fib(n-2)
}
