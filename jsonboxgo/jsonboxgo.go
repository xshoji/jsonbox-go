package jsonboxgo

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Client interface {
	Create(string, interface{}) []byte
	Read(string, string) ([]byte, bool)
	ReadAll(string) []byte
	ReadByQuery(string, string) []byte
	Update(string, string, interface{}) ([]byte, bool)
	Delete(string, string) ([]byte, bool)
}

type DefaultClient struct {
	baseUrl     string
	boxId       string
	baseUrlFull string
	httpClient  *http.Client
}

// Create new jsonbox-go Client
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
	resp, err := c.doRequest("POST", collection, "", "", object)
	if err != nil {
		log.Fatal("Create failed. | ", err)
	}
	return readAsBytes(resp)
}

// Read all
func (c DefaultClient) ReadAll(collection string) []byte {
	resp, err := c.doRequest("GET", collection, "", "", nil)
	if err != nil {
		log.Fatal("ReadAll failed. | ", err)
	}
	return readAsBytes(resp)
}

// Read all
func (c DefaultClient) ReadByQuery(collection string, query string) []byte {
	resp, err := c.doRequest("GET", collection, "", query, nil)
	if err != nil {
		log.Fatal("ReadByQuery failed. | ", err)
	}
	return readAsBytes(resp)
}

// Read one
func (c DefaultClient) Read(collection string, recordId string) (respondedBody []byte, found bool) {
	resp, err := c.doRequest("GET", collection, recordId, "", nil)
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
	resp, err := c.doRequest("PUT", collection, recordId, "", object)
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
	resp, err := c.doRequest("DELETE", collection, recordId, "", nil)
	if err != nil {
		log.Fatal("Delete failed. | ", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, false
	}
	return readAsBytes(resp), true
}

func (c DefaultClient) doRequest(httpMethod string, collection string, recordId string, query string, object interface{}) (*http.Response, error) {
	var body io.Reader = nil
	if object != nil {
		requestBody := toJsonString(object)
		body = strings.NewReader(requestBody)
	}
	req, err := http.NewRequest(httpMethod, c.baseUrlFull+handleSuffixAndPrefix(collection)+handleSuffixAndPrefix(recordId)+query, body)
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

type QueryBuilder interface {
	Limit(int) QueryBuilder
	Offset(int) QueryBuilder
	SortAsc(string) QueryBuilder
	SortDesc(string) QueryBuilder
	AddEqual(string, string) QueryBuilder
	AddGreaterThan(string, string) QueryBuilder
	AddGreaterThanOrEqual(string, string) QueryBuilder
	AddLessThan(string, string) QueryBuilder
	AddLessThanOrEqual(string, string) QueryBuilder
	Build() string
}

type DefaultQueryBuilder struct {
	queries []string
	filters []string
}

// Create new jsonbox-go QueryBuilder
func NewQueryBuilder() QueryBuilder {
	builder := &DefaultQueryBuilder{
		queries: make([]string, 0),
		filters: make([]string, 0),
	}
	return builder
}

func (d *DefaultQueryBuilder) Limit(limit int) QueryBuilder {
	d.queries = append(d.queries, `limit=`+strconv.Itoa(limit))
	return d
}

func (d *DefaultQueryBuilder) Offset(offset int) QueryBuilder {
	d.queries = append(d.queries, `offset=`+strconv.Itoa(offset))
	return d
}

func (d *DefaultQueryBuilder) SortAsc(sort string) QueryBuilder {
	d.queries = append(d.queries, `sort=`+sort)
	return d
}

func (d *DefaultQueryBuilder) SortDesc(sort string) QueryBuilder {
	d.queries = append(d.queries, `sort=-`+sort)
	return d
}

func (d *DefaultQueryBuilder) AddGreaterThan(field string, value string) QueryBuilder {
	return d.addFilter(field, ":>", value)
}

func (d *DefaultQueryBuilder) AddLessThan(field string, value string) QueryBuilder {
	return d.addFilter(field, ":<", value)
}

func (d *DefaultQueryBuilder) AddGreaterThanOrEqual(field string, value string) QueryBuilder {
	return d.addFilter(field, ":>=", value)
}

func (d *DefaultQueryBuilder) AddLessThanOrEqual(field string, value string) QueryBuilder {
	return d.addFilter(field, ":<=", value)
}

func (d *DefaultQueryBuilder) AddEqual(field string, value string) QueryBuilder {
	return d.addFilter(field, ":=", value)
}

func (d *DefaultQueryBuilder) addFilter(name string, operator string, value string) QueryBuilder {
	filterQuery := ""
	if len(d.filters) == 0 {
		filterQuery += "q="
	}
	d.filters = append(d.filters, filterQuery+name+operator+value)
	return d
}

func (d *DefaultQueryBuilder) Build() string {
	query := ""
	if len(d.queries) > 0 {
		query += strings.Join(d.queries, `&`)
	}
	if len(d.filters) > 0 {
		if query != "" {
			query += "&"
		}
		query += strings.Join(d.filters, `,`)
	}
	return "?" + query
}
