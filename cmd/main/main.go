package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/xshoji/jsonbox-go/jsonboxgo"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const baseUrl = "https://jsonbox.io/"

func main() {
	boxId := os.Getenv("BOX_ID")
	if boxId == "" {
		log.Fatal("Environment variable \"BOX_ID\" is not defined.")
	}
	collection := "users"
	client := jsonboxgo.NewClient(baseUrl, boxId)

	anonymousStruct := struct {
		Name     string
		Age      int
		Language string
	}{
		Name:     "taro_" + createRandomString(),
		Age:      createRandomNumber(),
		Language: "JP",
	}

	// Create
	result := client.Create(collection, anonymousStruct)
	fmt.Println(">>> Create")
	fmt.Println(result)
	fmt.Println("")

	// ReadAll
	result = client.ReadAll("users")
	fmt.Println(">>> ReadAll")
	fmt.Println(result)
	fmt.Println("")

	// Read by recordId
	var jsonObject []interface{}
	json.Unmarshal([]byte(result), &jsonObject)
	recordId := jsonObject[0].(map[string]interface{})["_id"].(string)
	result = client.Read("users", recordId)
	fmt.Println(">>> Read")
	fmt.Println(result)
	fmt.Println("")

	// Update
	anonymousStruct.Name = "updated_" + createRandomString()
	anonymousStruct.Age = createRandomNumber()
	result = client.Update("users", recordId, anonymousStruct)
	fmt.Println(">>> Update")
	fmt.Println(result)
	fmt.Println("")

	// Delete
	result = client.Delete("users", recordId)
	fmt.Println(">>> Delete")
	fmt.Println(result)
	fmt.Println("")

	// Read by recordId
	result = client.Read("users", recordId)
	fmt.Println(">>> Read")
	fmt.Println(result)
	fmt.Println("")
}

func createRandomString() string {
	seed := strconv.FormatInt(time.Now().UnixNano(), 10)
	shaBytes := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(shaBytes[:])
}

func createRandomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100-1) + 1
}
