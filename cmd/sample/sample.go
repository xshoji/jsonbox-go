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

type User struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Age       int    `json:"age,omitempty"`
	Language  string `json:"language,omitempty"`
	CreatedOn string `json:"_createdOn,omitempty"`
}

func main() {
	boxId := os.Getenv("BOX_ID")
	if boxId == "" {
		log.Fatal("Environment variable \"BOX_ID\" is not defined.")
	}
	collection := "users"
	client := jsonboxgo.NewClient(baseUrl, boxId)

	user := User{
		Name:     "taro_" + randomString(),
		Age:      randomNumber(),
		Language: "JP",
	}

	// Create
	result := client.Create(collection, user)
	fmt.Println(">>> Create")
	fmt.Println(string(result))
	fmt.Println("")

	// Bind struct
	var createdUser User
	json.Unmarshal(result, &createdUser)

	// ReadAll
	result = client.ReadAll("users")
	fmt.Println(">>> ReadAll")
	fmt.Println(string(result))
	fmt.Println("")

	// Read by recordId
	result, _ = client.Read(collection, createdUser.Id)
	fmt.Println(">>> Read")
	fmt.Println(string(result))
	fmt.Println("")

	// Update
	createdUser.Name = "updated_" + randomString()
	createdUser.Age = randomNumber()
	result, _ = client.Update("users", createdUser.Id, createdUser)
	fmt.Println(">>> Update")
	fmt.Println(string(result))
	fmt.Println("")

	// Read by recordId
	result, _ = client.Read(collection, createdUser.Id)
	fmt.Println(">>> Read (Updated)")
	fmt.Println(string(result))
	fmt.Println("")

	// Delete
	result, _ = client.Delete("users", createdUser.Id)
	fmt.Println(">>> Delete")
	fmt.Println(string(result))
	fmt.Println("")

	// Read by recordId
	_, found := client.Read("users", createdUser.Id)
	fmt.Println(">>> Read (Deleted)")
	fmt.Println(found)
	fmt.Println("")
}

func randomString() string {
	seed := strconv.FormatInt(time.Now().UnixNano(), 10)
	shaBytes := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(shaBytes[:])
}

func randomNumber() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100-1) + 1
}
