package jsonboxgo

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Client interface {
	Create(string, interface{}) string
	Read(string, string) string
	ReadAll(string) string
	Update(string, string, interface{}) string
	Delete(string, string) string
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
		baseUrlFull: handleSuffix(baseUrl) + handlePrefix(handleSuffix(boxId)),
	}
	return client
}

func (c clientImpl) Create(collection string, object interface{}) string {
	requestBody := toJsonString(object)
	resp, err := http.Post(c.baseUrlFull+handlePrefix(handleSuffix(collection)), "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal("Create failed.")
	}
	return asString(resp)
}

func (c clientImpl) ReadAll(collection string) string {
	resp, err := http.Get(c.baseUrlFull + handlePrefix(handleSuffix(collection)))
	if err != nil {
		log.Fatal("ReadAll failed.")
	}
	return asString(resp)
}

func (c clientImpl) Update(collection string, recordId string, object interface{}) string {
	requestBody := toJsonString(object)
	resp, err := http.Post(c.baseUrlFull+handlePrefix(handleSuffix(collection))+handlePrefix(handleSuffix(recordId)), "application/json", strings.NewReader(requestBody))
	if err != nil {
		log.Fatal("Update failed.")
	}
	return asString(resp)
}

func (c clientImpl) Delete(collection string, recordId string) string {
	req, err := http.NewRequest("DELETE", c.baseUrlFull+handlePrefix(handleSuffix(collection))+handlePrefix(handleSuffix(recordId)), nil)
	if err != nil {
		log.Fatal("NewRequest(\"DELETE\") failed.")
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("Delete failed.")
	}
	return asString(resp)
}

func (c clientImpl) Read(collection string, recordId string) string {
	resp, err := http.Get(c.baseUrlFull + handlePrefix(handleSuffix(collection)) + handlePrefix(handleSuffix(recordId)))
	if err != nil {
		log.Fatal("Read failed.")
	}
	return asString(resp)
}

func asString(resp *http.Response) string {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("ioutil.ReadAll() failed.")
	}
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			log.Panic("resp.Body.Close() failed.")
		}
	}()
	return string(body)
}

// 値をjson形式の文字列に変換する
func toJsonString(v interface{}) string {
	result, _ := json.Marshal(v)
	return string(result)
}

func handleSuffix(char string) string {
	if strings.HasSuffix(char, "/") {
		return strings.TrimRight(char, "/")
	}
	return char
}

func handlePrefix(char string) string {
	if !strings.HasPrefix(char, "/") {
		return "/" + char
	}
	return char
}
