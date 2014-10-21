package engine

import (
	"bytes"
	"testing"
)

func testSignal(method string, t *testing.T) {
	var sent bool = false

	e, _ := Basic()

	go func() {
		for {
			select {
			case sig := <-e.Signals:
				if string(sig) != "SENT" {
					t.Errorf("Read signal is not `SENT`")
				}
			}
		}
	}()

	e.Handle("/test_signal_sent", method, func(c *Ctx) {
		sent = true
		e.Send("", "SENT")
	})

	PerformRequest(e, method, "/test_signal_sent")

	if sent == false {
		t.Errorf("Signal handler was not invoked.")
	}

}

func testSignalTrue(method string, t *testing.T) {
	var sent bool = false

	e, _ := Basic()

	go func() {
		for {
			select {
			case sig := <-e.Signals:
				if bytes.Compare([]byte("true"), sig) != 0 {
					t.Errorf("Read signal is not `true`")
				}
			}
		}
	}()

	e.Handle("/test_signal_true", method, func(c *Ctx) {
		sent = true
		for i := 0; i < 100; i++ {
			ts := []byte("true")
			e.Signals <- ts
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
