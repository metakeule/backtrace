/*
Package backtrace offers a simple backtrace for a max of 100 calls
*/
package backtrace

import (
	"fmt"
	"net/http"
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

// the number is the number of skipped entries in the backtrace
type HTTPPanicCatcher int

func (h HTTPPanicCatcher) Catch(recovered interface{}, rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(rw, "500 Server panicked!\n\n\n")

	btComplete := BackTrace()[int(h):]
	for _, info := range btComplete {
		// bf.WriteString(fmt.)
		fmt.Fprintf(rw, "Function: %s\nFile: %s:%d\n\n", info.Function, info.File, info.Line)
	}
}
