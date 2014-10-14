package engine

import (
	"fmt"
	"log"
	"os"
	"reflect"

	"testing"
)

type testitem struct {
	c        Conf
	fname    string
	expected interface{}
}

func testConf(testitems []*testitem, t *testing.T) {
	e, _ := New()
	var c []Conf
	for _, tt := range testitems {
		c = append(c, tt.c)
	}
	err := e.SetConf(c...)
	if err != nil {
		t.Errorf(fmt.Sprintf("Engine returned configuration error: %+v", err))
	}
	val := reflect.ValueOf(e).Elem()
	for _, ci := range testitems {
		f := val.FieldByName(ci.fname)
		if f.Interface() != ci.expected {
			t.Errorf(fmt.Sprintf("engine.%s is %+v, but should be %v\n", ci.fname, f.Interface(), ci.expected))
		}
	}
}

func TestConf(t *testing.T) {
	l := log.New(os.Stdout, "[TEST]", 0)
	tc := []*testitem{
		&testitem{ServePanic(false), "ServePanic", false},
		&testitem{RedirectTrailingSlash(false), "RedirectTrailingSlash", false},
		&testitem{RedirectFixedPath(false), "RedirectFixedPath", false},
		&testitem{HTMLStatus(true), "HTMLStatus", true},
		&testitem{LoggingOn(true), "LoggingOn", true},
		&testitem{Logger(l), "Logger", l},
		&testitem{MaxFormMemory(500), "MaxFormMemory", int64(500)},
	}
	testConf(tc, t)
}
