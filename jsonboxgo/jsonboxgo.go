//go:generate mockgen -source=$GOFILE -destination=mock_$GOFILE -package=$GOPACKAGE
package jsonboxgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client interface {
	Create(string, interface{}) []byte
	Read(string, string) ([]byte, bool)
	ReadAll(string) []byte
	Update(string, string, interface{}) ([]byte, bool)
	Delete(string, string) ([]byte, bool)
}

type DefaultClient struct {
	baseUrl     string
	boxId       string
	baseUrlFull string
	httpClient  *http.Client
}

// Create new jsonbox-go client
func NewClient(baseUrl string, boxId string, httpClient *http.Client) Client {
	client := DefaultClient{
		baseUrl:     baseUrl,
		boxId:       boxId,
		baseUrlFull: handleSuffix(baseUrl) + handleSuffixAndPrefix(boxId),
		httpClient:  httpClient,
	}
	return client
}

// Create
func (c DefaultClient) Create(collection string, object interface{}) []byte {
	resp, err := c.doRequest("POST", collection, "", object)
	if err != nil {
		log.Fatal("Create failed. | ", err)
	}
	return readAsBytes(resp)
}

// Read all
func (c DefaultClient) ReadAll(collection string) []byte {
	resp, err := http.Get(c.baseUrlFull + handleSuffixAndPrefix(collection))
	if err != nil {
		log.Fatal("ReadAll failed. | ", err)
	}
	return readAsBytes(resp)
}

// Read one
func (c DefaultClient) Read(collection string, recordId string) (respondedBody []byte, found bool) {
	resp, err := c.doRequest("GET", collection, recordId, nil)
	if err != nil {
		log.Fatal("Read failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	bytes := readAsBytes(resp)
	// list type json object is unexpected.
	var listObject []interface{}
	if json.Unmarshal(bytes, &listObject) == nil {
		return nil, false
	}
	return bytes, true
}

// Update
func (c DefaultClient) Update(collection string, recordId string, object interface{}) (respondedBody []byte, updated bool) {
	resp, err := c.doRequest("PUT", collection, recordId, object)
	if err != nil {
		log.Fatal("Update failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return c.Read(collection, recordId)
}

// Delete
func (c DefaultClient) Delete(collection string, recordId string) (respondedBody []byte, deleted bool) {
	resp, err := c.doRequest("DELETE", collection, recordId, nil)
	if err != nil {
		log.Fatal("Delete failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return readAsBytes(resp), true
}

func (c DefaultClient) doRequest(httpMethod string, collection string, recordId string, object interface{}) (*http.Response, error) {
	var body io.Reader = nil
	if object != nil {
		requestBody := toJsonString(object)
		body = strings.NewReader(requestBody)
	}
	req, err := http.NewRequest(httpMethod, c.baseUrlFull+handleSuffixAndPrefix(collection)+handleSuffixAndPrefix(recordId), body)
	if err != nil {
		log.Fatal(`http.NewRequest("`+httpMethod+`") failed. | `, err)
	}
	req.Header.Set("Content-Type", "application/json")
	return c.httpClient.Do(req)
}

func readAsBytes(resp *http.Response) []byte {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll() failed. | ", err)
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Panic("resp.Body.Close() failed. | ", err)
		}
	}()
	return body
}

// Convert to json string
func toJsonString(v interface{}) string {
	result, _ := json.Marshal(v)
	return string(result)
}

// Adjust suffix
func handleSuffix(char string) string {
	if strings.HasSuffix(char, "/") {
		return strings.TrimRight(char, "/")
	}
	return char
}

// Adjust prefix
func handlePrefix(char string) string {
	if !strings.HasPrefix(char, "/") {
		return "/" + char
	}
	return char
}

// Adjust suffix and prefix
func handleSuffixAndPrefix(char string) string {
	return handlePrefix(handleSuffix(char))
}
