/*
Package backtrace offers a simple backtrace for a max of 100 calls
*/
package backtrace

import (
	"fmt"
	"io"
	"net/http"
	"os"
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

	// no checking of len of backtrace, since h should be the minimum of calls
	// that always come before a panic might happen
	btComplete := BackTrace()[int(h):]
	for _, info := range btComplete {
		// bf.WriteString(fmt.)
		fmt.Fprintf(rw, "Function: %s\nFile: %s:%d\n\n", info.Function, info.File, info.Line)
	}
}

type FmtPanicCatcher struct {
	Skip   int
	Max    int
	format string    // expects %s %s %s %s and %d for request method, request url, Function, File and Line
	writer io.Writer // where the panic is written to
}

// "Function: %s\nFile: %s:%d\n\n"

// NewPanicCatcher creates a new FmtPanicCatcher that skips the number of entries given
// by skip and has max entries. If max is < 0, all entries after skip are shown.
func NewPanicCatcher(skip, max int) *FmtPanicCatcher {
	return &FmtPanicCatcher{skip, max, "Function: %s\nFile: %s:%d\n\n", os.Stdout}
}

func (f *FmtPanicCatcher) SetFormat(format string) *FmtPanicCatcher {
	f.format = format
	return f
}

func (f *FmtPanicCatcher) SetWriter(wr io.Writer) *FmtPanicCatcher {
	f.writer = wr
	return f
}

func (f *FmtPanicCatcher) Catch(recovered interface{}, rw http.ResponseWriter, req *http.Request) {
	rw.WriteHeader(http.StatusInternalServerError)
	btComplete := BackTrace()[int(f.Skip):]
	if f.Max > 0 {
		btComplete = btComplete[:f.Max]
	}

	for _, info := range btComplete {
		fmt.Fprintf(f.writer, f.format, info.Function, info.File, info.Line)
	}
}
