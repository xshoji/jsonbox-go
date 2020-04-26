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

type clientImpl struct {
	baseUrl     string
	boxId       string
	baseUrlFull string
}

func NewClient(baseUrl string, boxId string) Client {
	client := clientImpl{
		baseUrl:     baseUrl,
		boxId:       boxId,
		baseUrlFull: handleSuffix(baseUrl) + handleSuffixAndPrefix(boxId),
	}
	return client
}

func (c clientImpl) Create(collection string, object interface{}) []byte {
	requestBody := toJsonString(object)
	resp, err := http.Post(c.baseUrlFull+handleSuffixAndPrefix(collection), "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal("Create failed. | ", err)
	}
	log.Println(resp.Status)
	return readAsBytes(resp)
}

func (c clientImpl) ReadAll(collection string) []byte {
	resp, err := http.Get(c.baseUrlFull + handleSuffixAndPrefix(collection))
	if err != nil {
		log.Fatal("ReadAll failed. | ", err)
	}
	log.Println(resp.Status)
	return readAsBytes(resp)
}

func (c clientImpl) Update(collection string, recordId string, object interface{}) (respondedBody []byte, updated bool) {
	resp, err := c.doRequest("PUT", collection, recordId, object)
	if err != nil {
		log.Fatal("Update failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return readAsBytes(resp), true
}

func (c clientImpl) Delete(collection string, recordId string) (respondedBody []byte, deleted bool) {
	resp, err := c.doRequest("DELETE", collection, recordId, nil)
	if err != nil {
		log.Fatal("Delete failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return readAsBytes(resp), true
}

func (c clientImpl) Read(collection string, recordId string) (respondedBody []byte, found bool) {
	resp, err := http.Get(c.baseUrlFull + handleSuffixAndPrefix(collection) + handleSuffixAndPrefix(recordId))
	if err != nil {
		log.Fatal("Read failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return readAsBytes(resp), true
}

func (c clientImpl) doRequest(httpMethod string, collection string, recordId string, object interface{}) (*http.Response, error) {
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
	return http.DefaultClient.Do(req)
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
