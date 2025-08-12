package erroring

import (
	"SzymonZet/LuxmedCheck/tester"
	"bytes"
	"errors"
	"log"
	"testing"
)

func TestQuitOnError(t *testing.T) {
	var gotBuffer bytes.Buffer
	log.SetOutput(&gotBuffer)

	testCases := []struct {
		testName    string
		inputErr    error
		inputMsg    string
		wantBuffer  string
		wantIsPanic bool
	}{
		{"no error", nil, "Message", "", false},
		{"error", errors.New("crash!"), "app crashed", "ERR | FATAL | app crashed | crash!", true},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			defer func() {
				gotIsPanic := false
				if r := recover(); r != nil {
					gotIsPanic = true
				}
				tester.Assert(t, test.wantIsPanic, gotIsPanic)
				tester.AssertContains(t, gotBuffer.String(), test.wantBuffer)

			}()
			QuitIfError(test.inputErr, test.inputMsg)
		})
	}
}
