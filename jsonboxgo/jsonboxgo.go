package jsonboxgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client interface {
	Create(string, interface{}) []byte
	Read(string, string) []byte
	ReadAll(string) []byte
	Update(string, string, interface{}) []byte
	Delete(string, string) []byte
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
	return readAsBytes(resp)
}

func (c clientImpl) ReadAll(collection string) []byte {
	resp, err := http.Get(c.baseUrlFull + handleSuffixAndPrefix(collection))
	if err != nil {
		log.Fatal("ReadAll failed. | ", err)
	}
	return readAsBytes(resp)
}

func (c clientImpl) Update(collection string, recordId string, object interface{}) []byte {
	requestBody := toJsonString(object)
	req, err := http.NewRequest("PUT", c.baseUrlFull+handleSuffixAndPrefix(collection)+handleSuffixAndPrefix(recordId), strings.NewReader(requestBody))
	if err != nil {
		log.Fatal(`http.NewRequest("PUT") failed. | `, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Update failed. | ", err)
	}
	return readAsBytes(resp)
}

func (c clientImpl) Delete(collection string, recordId string) []byte {
	req, err := http.NewRequest("DELETE", c.baseUrlFull+handleSuffixAndPrefix(collection)+handleSuffixAndPrefix(recordId), nil)
	if err != nil {
		log.Fatal(`http.NewRequest("DELETE") failed. | `, err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Delete failed. | ", err)
	}
	return readAsBytes(resp)
}

func (c clientImpl) Read(collection string, recordId string) []byte {
	resp, err := http.Get(c.baseUrlFull + handleSuffixAndPrefix(collection) + handleSuffixAndPrefix(recordId))
	if err != nil {
		log.Fatal("Read failed. | ", err)
	}
	return readAsBytes(resp)
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
