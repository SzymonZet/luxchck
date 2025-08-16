package lux

import (
	"bytes"
	"log"
	"szymonzet/luxchck/tester"
	"testing"
)

func TestGetFullUrl(t *testing.T) {
	var gotBuffer bytes.Buffer
	log.SetOutput(&gotBuffer)

	testCases := []struct {
		testName      string
		inputEndpoint string
		wantBuffer    string
		wantUrl       string
	}{
		{"no slashes", "aaa", "", "https://portalpacjenta.luxmed.pl/aaa"},
		{"one slash", "/aaa", "", "https://portalpacjenta.luxmed.pl/aaa"},
		{"two slashes", "//aaa", "", "https://portalpacjenta.luxmed.pl/aaa"},
	}

	for _, test := range testCases {
		t.Run(test.testName, func(t *testing.T) {
			got := getFullUrl(test.inputEndpoint)
			tester.Assert(t, got, test.wantUrl)

			// the idea was for function to fail when parsing url
			// the problem was, I couldn't find a case in which it actually fails
			tester.AssertContains(t, gotBuffer.String(), test.wantBuffer)
		})
	}
}
