package engine

import (
	"fmt"
	"reflect"
	"testing"
)

func testConf(c Conf, fname string, expected interface{}, t *testing.T) {
	e, _ := Basic()
	err := e.SetConf(c)
	if err != nil {
		t.Errorf(fmt.Sprintf("error with configuration setting %+v: %+v", fname, err))
	}
	val := reflect.ValueOf(e).Elem()
	f := val.FieldByName(fname)
	if f.Interface() != expected {
		t.Errorf(fmt.Sprintf("%s: %+v should equal %v\n", fname, f.Interface(), expected))
	}
}

func TestConf(t *testing.T) {
	testConf(ServePanic(false), "ServePanic", false, t)
	testConf(RedirectTrailingSlash(false), "RedirectTrailingSlash", false, t)
	testConf(RedirectFixedPath(false), "RedirectFixedPath", false, t)
	testConf(HTMLStatus(true), "HTMLStatus", true, t)
	testConf(LoggingOn(false), "LoggingOn", false, t)
	testConf(MaxFormMemory(500), "MaxFormMemory", int64(500), t)
}
