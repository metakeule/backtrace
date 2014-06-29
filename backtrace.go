/*
Package backtrace offers a simple backtrace for a max of 100 calls
*/
package backtrace

import (
	"runtime"
)

type FootPrint struct {
	Line     int
	File     string
	Function string
}

func BackTrace() (trace []FootPrint) {
	for i := 0; i < 100; i++ {
		pc, file, line, _ := runtime.Caller(1 + i)
		if file == "" {
			continue
		}
		if f := runtime.FuncForPC(pc); f != nil {
			trace = append(trace, FootPrint{line, file, f.Name()})
			continue
		}
		trace = append(trace, FootPrint{line, file, ""})
	}
	return
}

func Filter(trace []FootPrint, fn func(index int, fp FootPrint) bool) []FootPrint {
	res := []FootPrint{}
	for i, fp := range trace {
		if fn(i, fp) {
			res = append(res, fp)
		}
	}
	return res
}
