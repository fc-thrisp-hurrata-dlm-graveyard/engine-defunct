package engine

import "testing"

func testSignal(method string, t *testing.T) {
	var sent bool = false

	e, _ := Basic()

	testqueue := func(s string) {
		if s != "SENT" {
			t.Errorf("Read signal is not `SENT`")
		}
	}

	e.Queues["testqueue"] = testqueue

	e.Handle("/test_signal_sent", method, func(c *Ctx) {
		sent = true
		e.Send("testqueue", "SENT")
	})

	PerformRequest(e, method, "/test_signal_sent")

	if sent == false {
		t.Errorf("Signal handler was not invoked.")
	}

}

func testSignalTrue(method string, t *testing.T) {
	var sent bool = false

	e, _ := Basic()

	testtrue := func(s string) {
		if s != "true" {
			t.Errorf("Read signal is not `true`")
		}
	}

	e.Queues["testtrue"] = testtrue

	e.Handle("/test_signal_true", method, func(c *Ctx) {
		sent = true
		for i := 0; i < 100; i++ {
			e.Send("testtrue", "true")
		}
	})

	PerformRequest(e, method, "/test_signal_true")

	if sent == false {
		t.Errorf("Signal handler was not invoked.")
	}
}

func TestSignal(t *testing.T) {
	t.Parallel()
	testSignal("HEAD", t)
	testSignalTrue("HEAD", t)
}
