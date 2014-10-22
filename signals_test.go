package engine

import (
	"fmt"
	"testing"
)

func testSignal(method string, t *testing.T) {
	var sent bool = false

	e, _ := New()

	testsignalq := func() {
		go func() {
			for msg := range e.Signals {
				fmt.Printf("test: %s\n", msg)
			}
		}()
	}

	testsignalq()

	testqueue := func(s string) {
		if s != "SENT" {
			t.Errorf("Read signal is not `SENT`")
		}
	}

	e.Queues["testqueue"] = testqueue

	e.Handle("/test_signal_sent", method, func(c *Ctx) {
		sent = true
		for i := 0; i < 10; i++ {
			e.Send("testqueue", "SENT")
		}
	})

	PerformRequest(e, method, "/test_signal_sent")

	if sent == false {
		t.Errorf("Signal handler was not invoked.")
	}

}

func TestSignal(t *testing.T) {
	testSignal("POST", t)
	testSignal("DELETE", t)
	testSignal("PUT", t)
	testSignal("PATCH", t)
	testSignal("OPTIONS", t)
	testSignal("HEAD", t)
}
