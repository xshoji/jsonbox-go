package jsonboxgo

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

type RoundTripFunc func(req *http.Request) *http.Response

func (f RoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

func NewTestClient(fn RoundTripFunc) *http.Client {
	return &http.Client{
		Transport: fn,
	}
}

func CreateNewTestClient(respondedBody string) *http.Client {
	// > Unit Testing http client in Go ｜ hassansin
	// > http://hassansin.github.io/Unit-Testing-http-client-in-Go
	return NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: 200,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(respondedBody)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

func TestNewClient(t *testing.T) {
	// test cases
	testCases := map[string]struct {
		TestCase        string
		InputBaseUrl    string
		InputBoxId      string
		ExpectedUrlFull string
	}{
		"Slash exists on last of InputBaseUrl.": {
			InputBaseUrl:    "https://test.com/",
			InputBoxId:      "box_xxxxx",
			ExpectedUrlFull: "https://test.com/box_xxxxx",
		},
		"Slash exists on first of InputBoxId.": {
			InputBaseUrl:    "https://test.com",
			InputBoxId:      "/box_xxxxx",
			ExpectedUrlFull: "https://test.com/box_xxxxx",
		},
		"Slash not exists.": {
			InputBaseUrl:    "https://test.com",
			InputBoxId:      "box_xxxxx",
			ExpectedUrlFull: "https://test.com/box_xxxxx",
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			client := NewClient(param.InputBaseUrl, param.InputBoxId, nil)
			t.Logf("Case:%v\n", param.TestCase)
			defaultClient := client.(defaultClient)
			actual := defaultClient.baseUrlFull
			expected := param.ExpectedUrlFull
			if actual != expected {
				t.Errorf("  Failed: actual -> %v(%T), expected -> %v(%T)\n", actual, actual, expected, expected)
			}
		})
	}
}

func TestCreate(t *testing.T) {
	// test cases
	InputBaseUrl := "https://test.com"
	InputBoxId := "box_test"
	testCases := map[string]struct {
		InputCollection    string
		InputObject        interface{}
		InputRespondedBody string
		ExpectedUrlFull    string
	}{
		"Success case.": {
			InputCollection: "users",
			InputObject: struct {
				Name string `json:"name,omitempty"`
			}{
				Name: `taro`,
			},
			InputRespondedBody: `{"_id":"aaaa","name":"taro","_createdOn":"2020-04-26T16:26:13.935Z"}`,
			ExpectedUrlFull:    "https://test.com/box_xxxxx",
		},
		"Success case 2.": {
			InputCollection: "/users",
			InputObject: struct {
				Name string `json:"name,omitempty"`
			}{
				Name: `taro`,
			},
			InputRespondedBody: `{"_id":"aaaa","name":"taro","_createdOn":"2020-04-26T16:26:13.935Z"}`,
			ExpectedUrlFull:    "https://test.com/box_xxxxx",
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			mockHttpClient := CreateNewTestClient(param.InputRespondedBody)
			client := NewClient(InputBaseUrl, InputBoxId, mockHttpClient)
			result := client.Create(param.InputCollection, param.InputObject)
			actual := string(result)
			expected := param.InputRespondedBody
			if actual != expected {
				t.Errorf("  Failed: actual -> %v(%T), expected -> %v(%T)\n", actual, actual, expected, expected)
			}
		})
	}
}
