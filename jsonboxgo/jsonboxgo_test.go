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

func CreateNewTestClient(statusCode int, respondedBody string) *http.Client {
	// > Unit Testing http client in Go ｜ hassansin
	// > http://hassansin.github.io/Unit-Testing-http-client-in-Go
	return NewTestClient(func(req *http.Request) *http.Response {
		// Test request parameters
		return &http.Response{
			StatusCode: statusCode,
			// Send response to be tested
			Body: ioutil.NopCloser(bytes.NewBufferString(respondedBody)),
			// Must be set to non-nil value or it panics
			Header: make(http.Header),
		}
	})
}

type User struct {
	Id   string `json:"_id,omitempty"`
	Name string `json:"name,omitempty"`
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
			defaultClient := client.(DefaultClient)
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
		InputObject        User
		InputRespondedBody string
		ExpectedUrlFull    string
	}{
		"Success case.": {
			InputCollection:    "users",
			InputObject:        User{Name: `taro`},
			InputRespondedBody: `{"_id":"aaaa","name":"taro","_createdOn":"2020-04-26T16:26:13.935Z"}`,
			ExpectedUrlFull:    "https://test.com/box_test",
		},
		"Success case 2.": {
			InputCollection:    "/users",
			InputObject:        User{Name: `taro`},
			InputRespondedBody: `{"_id":"aaaa","name":"taro","_createdOn":"2020-04-26T16:26:13.935Z"}`,
			ExpectedUrlFull:    "https://test.com/box_test",
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			mockHttpClient := CreateNewTestClient(200, param.InputRespondedBody)
			client := NewClient(InputBaseUrl, InputBoxId, mockHttpClient)
			defaultClient, _ := client.(DefaultClient)
			result := defaultClient.Create(param.InputCollection, param.InputObject)
			actual := string(result)
			expected := param.InputRespondedBody
			if actual != expected {
				t.Errorf("  Failed: actual -> %v(%T), expected -> %v(%T)\n", actual, actual, expected, expected)
			}
			actual = defaultClient.baseUrlFull
			expected = param.ExpectedUrlFull
			if actual != expected {
				t.Errorf("  Failed: actual -> %v(%T), expected -> %v(%T)\n", actual, actual, expected, expected)
			}
		})
	}
}

func TestRead(t *testing.T) {
	// test cases
	InputBaseUrl := "https://test.com"
	InputBoxId := "box_test"
	testCases := map[string]struct {
		InputCollection          string
		InputObject              User
		InputRespondedHttpStatus int
		InputRespondedBody       string
		ExpectedRespondedBody    string
		ExpectedUrlFull          string
		ExpectedFound            bool
	}{
		"Found case.": {
			InputCollection:          "users",
			InputObject:              User{Id: `id001`, Name: `taro`},
			InputRespondedHttpStatus: 200,
			InputRespondedBody:       `{"_id":"id001","name":"taro"}`,
			ExpectedRespondedBody:    `{"_id":"id001","name":"taro"}`,
			ExpectedUrlFull:          "https://test.com/box_test",
			ExpectedFound:            true,
		},
		"Not found case 1.": {
			InputCollection:          "/users",
			InputObject:              User{Id: `id001`, Name: `taro`},
			InputRespondedHttpStatus: 500,
			InputRespondedBody:       `{"message":"Cannot read property '_id' of null"}`,
			ExpectedRespondedBody:    ``,
			ExpectedUrlFull:          "https://test.com/box_test",
			ExpectedFound:            false,
		},
		"Not found case 2.": {
			InputCollection:          "/users",
			InputObject:              User{Id: `id001`, Name: `taro`},
			InputRespondedHttpStatus: 200,
			InputRespondedBody:       `[{"_id":"id002","name":"taro"}]`,
			ExpectedRespondedBody:    ``,
			ExpectedUrlFull:          "https://test.com/box_test",
			ExpectedFound:            false,
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			mockHttpClient := CreateNewTestClient(param.InputRespondedHttpStatus, param.InputRespondedBody)
			client := NewClient(InputBaseUrl, InputBoxId, mockHttpClient)
			defaultClient, _ := client.(DefaultClient)
			result, found := defaultClient.Read(param.InputCollection, param.InputObject.Id)
			actual := string(result)
			expectedBody := param.ExpectedRespondedBody
			if actual != expectedBody {
				t.Errorf("  Failed: actual -> %v(%T), expectedBody -> %v(%T)\n", actual, actual, expectedBody, expectedBody)
			}
			actual = defaultClient.baseUrlFull
			expected := param.ExpectedUrlFull
			if actual != expected {
				t.Errorf("  Failed: actual -> %v(%T), expected -> %v(%T)\n", actual, actual, expected, expected)
			}
			expectedFound := param.ExpectedFound
			if found != expectedFound {
				t.Errorf("  Failed: found -> %v(%T), expectedFound -> %v(%T)\n", found, found, expectedFound, expectedFound)
			}
		})
	}
}
