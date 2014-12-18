package backtrace

import (
	// "fmt"
	"strings"
	"testing"
)

func TestBackTrace(t *testing.T) {
	trace := Filter(BackTrace(), func(index int, fp FootPrint) bool {
		if strings.HasSuffix(fp.File, "runtime/proc.c") {
			return false
		}
		return true
	})

	if len(trace) != 3 {
		t.Errorf("len of trace should be 3, but is: %d", len(trace))
	}

	traceFunction0 := ".TestBackTrace"

	idx := strings.LastIndex(trace[0].Function, ".")

	if idx == -1 {
		t.Errorf("trace[0].Function must have .but is %#v", trace[0].Function)
	}

	if trace[0].Function[idx:] != traceFunction0 {
		t.Errorf("trace[0].Function[idx:] should be %#v but is %#v", traceFunction0, trace[0].Function[idx:])
	}

	if trace[0].Line != 10 {
		t.Errorf("trace[0].Line should be %d but is %d", 10, trace[0].Line)
	}

	endFile := "backtrace_test.go"
	if !strings.HasSuffix(trace[0].File, endFile) {
		t.Errorf("trace[0].File should end with %#v, but is %#v", endFile, trace[0].File)
	}

	traceFunction1 := "testing.tRunner"
	if trace[1].Function != traceFunction1 {
		t.Errorf("trace[1].Function should be %#v but is %#v", traceFunction1, trace[1].Function)
	}
}
