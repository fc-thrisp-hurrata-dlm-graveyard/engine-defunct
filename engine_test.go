package engine

import (
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
)

func PerformRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func methodNotMethod(method string) string {
	methods := []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS", "HEAD"}
	newmethod := methods[rand.Intn(len(methods))]
	if newmethod == method {
		methodNotMethod(newmethod)
	}
	return newmethod
}

func testRouteOK(method string, t *testing.T) {
	passed := false
	e := New()

	e.Handle("/test", method, func(c *Ctx) { passed = true })

	w := PerformRequest(e, method, "/test")

	if passed == false {
		t.Errorf(method + " route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

func TestRouteOK(t *testing.T) {
	testRouteOK("POST", t)
	testRouteOK("DELETE", t)
	testRouteOK("PATCH", t)
	testRouteOK("PUT", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("HEAD", t)
}

func testGroupOK(method string, t *testing.T) {
	passed := false
	e := New()

	e.Handle("/test_group", method, func(c *Ctx) { passed = true })

	w := PerformRequest(e, method, "/test_group")

	if passed == false {
		t.Errorf(method + " group route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

func TestGroupOK(t *testing.T) {
	testRouteOK("POST", t)
	testRouteOK("DELETE", t)
	testRouteOK("PATCH", t)
	testRouteOK("PUT", t)
	testRouteOK("OPTIONS", t)
	testRouteOK("HEAD", t)
}

func testSubGroupOK(method string, t *testing.T) {
	passed := false
	e := New()
	g := e.New("/test_group")
	g.Handle("/test_group_subgroup", method, func(c *Ctx) { passed = true })

	w := PerformRequest(e, method, "/test_group/test_group_subgroup")

	if passed == false {
		t.Errorf(method + " group route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

func TestSubGroupOK(t *testing.T) {
	testSubGroupOK("POST", t)
	testSubGroupOK("DELETE", t)
	testSubGroupOK("PATCH", t)
	testSubGroupOK("PUT", t)
	testSubGroupOK("OPTIONS", t)
	testSubGroupOK("HEAD", t)
}

func testRouteNotOK(method string, t *testing.T) {
	passed := false
	e := New()
	othermethod := methodNotMethod(method)
	e.Handle("/test_2", othermethod, func(c *Ctx) { passed = true })
	w := PerformRequest(e, method, "/test")

	if passed == true {
		t.Errorf(method + " route handler was invoked, when it should not")
	}
	if w.Code != http.StatusNotFound {
		t.Errorf("Status code should be %v, was %d. Location: %s", http.StatusNotFound, w.Code, w.HeaderMap.Get("Location"))
	}
}

func TestRouteNotOK(t *testing.T) {
	testRouteNotOK("POST", t)
	testRouteNotOK("DELETE", t)
	testRouteNotOK("PATCH", t)
	testRouteNotOK("PUT", t)
	testRouteNotOK("OPTIONS", t)
	testRouteNotOK("HEAD", t)
}
