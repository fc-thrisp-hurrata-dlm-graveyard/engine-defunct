package engine

import (
	"errors"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"golang.org/x/net/context"
)

func PerformRequest(h http.Handler, method string, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, req)
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
	e, _ := Basic()

	e.Take("/test", method, func(c context.Context) {
		passed = true
	})

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
	e, _ := Basic()

	e.Take("/test_group", method, func(c context.Context) { passed = true })

	w := PerformRequest(e, method, "/test_group")

	if passed == false {
		t.Errorf(method + " group route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}
}

func TestGroupOK(t *testing.T) {
	testGroupOK("POST", t)
	testGroupOK("DELETE", t)
	testGroupOK("PATCH", t)
	testGroupOK("PUT", t)
	testGroupOK("OPTIONS", t)
	testGroupOK("HEAD", t)
}

func testSubGroupOK(method string, t *testing.T) {
	passed := false
	e, _ := Basic()
	g := e.New("/test_group")
	g.Take("/test_group_subgroup", method, func(c context.Context) { passed = true })

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
	e, _ := Basic()
	othermethod := methodNotMethod(method)
	e.Take("/test_not_ok", othermethod, func(c context.Context) { passed = true })
	w := PerformRequest(e, method, "/test_not_ok")

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

type mockResponseWriter struct{}

func (m *mockResponseWriter) Header() (h http.Header) {
	return http.Header{}
}

func (m *mockResponseWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (m *mockResponseWriter) WriteString(s string) (n int, err error) {
	return len(s), nil
}

func (m *mockResponseWriter) WriteHeader(int) {}

func TestRouter(t *testing.T) {
	t.Parallel()
	engine, _ := Basic()

	routed := false

	engine.Manage("GET", "/user/:name", func(c context.Context) {
		routed = true
		want := Params{Param{"name", "gopher"}}
		curr := currentCtx(c)
		if !reflect.DeepEqual(curr.Params, want) {
			t.Fatalf("wrong wildcard values: want %v, got %v", want, curr.Params)
		}
	})

	w := new(mockResponseWriter)

	req, _ := http.NewRequest("GET", "/user/gopher", nil)
	engine.ServeHTTP(w, req)

	if !routed {
		t.Fatal("routing failed")
	}
}

type handlerStruct struct {
	handeled *bool
}

func (h handlerStruct) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	*h.handeled = true
}

func TestRouterAPI(t *testing.T) {
	var get, post, put, patch, delete, handler, handlerFunc bool

	httpHandler := handlerStruct{&handler}

	router, _ := Basic()
	router.Manage("GET", "/GET", func(c context.Context) {
		get = true
	})
	router.Manage("POST", "/POST", func(c context.Context) {
		post = true
	})
	router.Manage("PUT", "/PUT", func(c context.Context) {
		put = true
	})
	router.Manage("PATCH", "/PATCH", func(c context.Context) {
		patch = true
	})
	router.Manage("DELETE", "/DELETE", func(c context.Context) {
		delete = true
	})
	router.Handler("GET", "/Handler", httpHandler)
	router.HandlerFunc("GET", "/HandlerFunc", func(w http.ResponseWriter, r *http.Request) {
		handlerFunc = true
	})

	w := new(mockResponseWriter)

	r, _ := http.NewRequest("GET", "/GET", nil)
	router.ServeHTTP(w, r)
	if !get {
		t.Error("routing GET failed")
	}

	r, _ = http.NewRequest("POST", "/POST", nil)
	router.ServeHTTP(w, r)
	if !post {
		t.Error("routing POST failed")
	}

	r, _ = http.NewRequest("PUT", "/PUT", nil)
	router.ServeHTTP(w, r)
	if !put {
		t.Error("routing PUT failed")
	}

	r, _ = http.NewRequest("PATCH", "/PATCH", nil)
	router.ServeHTTP(w, r)
	if !patch {
		t.Error("routing PATCH failed")
	}

	r, _ = http.NewRequest("DELETE", "/DELETE", nil)
	router.ServeHTTP(w, r)
	if !delete {
		t.Error("routing DELETE failed")
	}

	r, _ = http.NewRequest("GET", "/Handler", nil)
	router.ServeHTTP(w, r)
	if !handler {
		t.Error("routing Handler failed")
	}

	r, _ = http.NewRequest("GET", "/HandlerFunc", nil)
	router.ServeHTTP(w, r)
	if !handlerFunc {
		t.Error("routing HandlerFunc failed")
	}
}

func TestRouterRoot(t *testing.T) {
	router, _ := Basic()
	recv := catchPanic(func() {
		router.Manage("GET", "noSlashRoot", nil)
	})
	if recv == nil {
		t.Fatal("registering path not beginning with '/' did not panic")
	}
}

func TestRouterLookup(t *testing.T) {
	routed := false
	wantHandle := func(_ context.Context) {
		routed = true
	}
	wantParams := Params{Param{"name", "gopher"}}

	router, _ := Basic()

	// try empty router first
	handle, _, tsr := router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}

	// insert route and try again
	router.Manage("GET", "/user/:name", wantHandle)

	handle, params, tsr := router.Lookup("GET", "/user/gopher")
	if handle == nil {
		t.Fatal("Got no handle!")
	} else {
		handle(nil)
		if !routed {
			t.Fatal("Routing failed!")
		}
	}

	if !reflect.DeepEqual(params, wantParams) {
		t.Fatalf("Wrong parameter values: want %v, got %v", wantParams, params)
	}

	handle, _, tsr = router.Lookup("GET", "/user/gopher/")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if !tsr {
		t.Error("Got no TSR recommendation!")
	}

	handle, _, tsr = router.Lookup("GET", "/nope")
	if handle != nil {
		t.Fatalf("Got handle for unregistered pattern: %v", handle)
	}
	if tsr {
		t.Error("Got wrong TSR recommendation!")
	}
}

type mockFileSystem struct {
	opened bool
}

func (mfs *mockFileSystem) Open(name string) (http.File, error) {
	mfs.opened = true
	return nil, errors.New("this is just a mock")
}

func TestRouterServeFiles(t *testing.T) {
	router, _ := Basic()
	mfs := &mockFileSystem{}

	recv := catchPanic(func() {
		router.ServeFiles("/noFilepath", mfs)
	})
	if recv == nil {
		t.Fatal("registering path not ending with '*filepath' did not panic")
	}

	router.ServeFiles("/*filepath", mfs)
	w := new(mockResponseWriter)
	r, _ := http.NewRequest("GET", "/favicon.ico", nil)
	router.ServeHTTP(w, r)
	if !mfs.opened {
		t.Error("serving file failed")
	}
}

func testMiddleware(method string, t *testing.T) {
	passed := false

	e, _ := Basic()

	e.Take("/test_middleware", method, func(c context.Context) {})

	e.Middleware(func(http.ResponseWriter, *http.Request) { passed = true })

	w := PerformRequest(e, method, "/test_middleware")

	if passed == false {
		t.Errorf(method + " group route handler was not invoked.")
	}
	if w.Code != http.StatusOK {
		t.Errorf("Status code should be %v, was %d", http.StatusOK, w.Code)
	}

}

func TestMiddleware(t *testing.T) {
	testMiddleware("POST", t)
	testMiddleware("DELETE", t)
	testMiddleware("PATCH", t)
	testMiddleware("PUT", t)
	testMiddleware("OPTIONS", t)
	testMiddleware("HEAD", t)
}
