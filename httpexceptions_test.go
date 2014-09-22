package engine

import (
	"fmt"
	"net/http"
	"testing"
)

func testCallException(code int, t *testing.T) {
	e := New()
	e.Handle("/test", "GET", func(c *Ctx) { c.Exception(code) })

	w := PerformRequest(e, "GET", "/test")

	if w.Code != code {
		t.Errorf("Status code should be %v, was %d", http.StatusNotFound, w.Code)
	}
}

func TestCallException(t *testing.T) {
	testCallException(404, t)
	testCustomException(418, t)
	testCallException(500, t)
}

func testCustomException(code int, t *testing.T) {
	result := fmt.Sprintf("CUSTOM %d", code)
	e := New()
	e.HttpExceptions[code].Update(func(c *Ctx) { c.rw.Write([]byte(result)) })
	e.Handle("/test", "GET", func(c *Ctx) { c.Exception(code) })

	w := PerformRequest(e, "GET", "/test")

	if w.Body.String() != result {
		t.Errorf("Body should be %s, was but was %s.", result, w.Body.String())
	}
	if w.Code != code {
		t.Errorf("Status code should be %v, was %d", code, w.Code)
	}
}

func TestCustomException(t *testing.T) {
	testCustomException(404, t)
	testCustomException(418, t)
	testCustomException(500, t)
}
