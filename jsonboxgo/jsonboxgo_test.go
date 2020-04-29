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

func CreateNewTestClient(statusCode int, respondedBody string, statusCodeSecond int, respondedBodySecond string) *http.Client {
	// > Unit Testing http client in Go ｜ hassansin
	// > http://hassansin.github.io/Unit-Testing-http-client-in-Go
	return NewTestClient(func() func(*http.Request) *http.Response {
		counter := 0
		// Test request parameters
		return func(req *http.Request) *http.Response {
			createResponse := func(code int, body string) *http.Response {
				return &http.Response{
					StatusCode: code,
					// Send response to be tested
					Body: ioutil.NopCloser(bytes.NewBufferString(body)),
					// Must be set to non-nil value or it panics
					Header: make(http.Header),
				}
			}
			if counter == 0 {
				counter++
				return createResponse(statusCode, respondedBody)
			}
			// counter > 0
			return createResponse(statusCodeSecond, respondedBodySecond)
		}
	}())
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
			mockHttpClient := CreateNewTestClient(200, param.InputRespondedBody, 0, ``)
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
			mockHttpClient := CreateNewTestClient(param.InputRespondedHttpStatus, param.InputRespondedBody, 0, ``)
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

func TestReadAll(t *testing.T) {
	// test cases
	InputBaseUrl := "https://test.com"
	InputBoxId := "box_test"
	testCases := map[string]struct {
		InputCollection          string
		InputRespondedHttpStatus int
		InputRespondedBody       string
		ExpectedRespondedBody    string
		ExpectedUrlFull          string
	}{
		"Found case.": {
			InputCollection:          "users",
			InputRespondedHttpStatus: 200,
			InputRespondedBody:       `[{"_id":"id001","name":"taro"}]`,
			ExpectedRespondedBody:    `[{"_id":"id001","name":"taro"}]`,
			ExpectedUrlFull:          "https://test.com/box_test",
		},
		"Not found case.": {
			InputCollection:          "/users",
			InputRespondedHttpStatus: 200,
			InputRespondedBody:       `[]`,
			ExpectedRespondedBody:    `[]`,
			ExpectedUrlFull:          "https://test.com/box_test",
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			mockHttpClient := CreateNewTestClient(param.InputRespondedHttpStatus, param.InputRespondedBody, 0, ``)
			client := NewClient(InputBaseUrl, InputBoxId, mockHttpClient)
			defaultClient, _ := client.(DefaultClient)
			result := defaultClient.ReadAll(param.InputCollection)
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
		})
	}
}

func TestUpdate(t *testing.T) {
	// test cases
	InputBaseUrl := "https://test.com"
	InputBoxId := "box_test"
	testCases := map[string]struct {
		InputCollection                string
		InputObject                    User
		InputRespondedHttpStatus       int
		InputRespondedBody             string
		InputRespondedHttpStatusSecond int
		InputRespondedBodySecond       string
		ExpectedRespondedBody          string
		ExpectedUrlFull                string
		ExpectedUpdated                bool
	}{
		"Updated case.": {
			InputCollection:                "users",
			InputObject:                    User{Id: `id001`, Name: `taro`},
			InputRespondedHttpStatus:       200,
			InputRespondedBody:             `{"message":"Record updated."}`,
			InputRespondedHttpStatusSecond: 200,
			InputRespondedBodySecond:       `{"_id":"id001","name":"taro"}`,
			ExpectedRespondedBody:          `{"_id":"id001","name":"taro"}`,
			ExpectedUrlFull:                "https://test.com/box_test",
			ExpectedUpdated:                true,
		},
		"Not found case.": {
			InputCollection:                "users",
			InputObject:                    User{Id: `id001`, Name: `taro`},
			InputRespondedHttpStatus:       400,
			InputRespondedBody:             `{"message":"Invalid record Id"}`,
			InputRespondedHttpStatusSecond: 0,
			InputRespondedBodySecond:       ``,
			ExpectedRespondedBody:          ``,
			ExpectedUrlFull:                "https://test.com/box_test",
			ExpectedUpdated:                false,
		},
	}

	// run
	for testCase, param := range testCases {
		t.Run(testCase, func(t *testing.T) {
			mockHttpClient := CreateNewTestClient(
				param.InputRespondedHttpStatus,
				param.InputRespondedBody,
				param.InputRespondedHttpStatusSecond,
				param.InputRespondedBodySecond,
			)
			client := NewClient(InputBaseUrl, InputBoxId, mockHttpClient)
			defaultClient, _ := client.(DefaultClient)
			result, updated := defaultClient.Update(param.InputCollection, param.InputObject.Id, param.InputObject)
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
			expectedUpdated := param.ExpectedUpdated
			if updated != expectedUpdated {
				t.Errorf("  Failed: updated -> %v(%T), expectedUpdated -> %v(%T)\n", updated, updated, expectedUpdated, expectedUpdated)
			}
		})
	}
}
