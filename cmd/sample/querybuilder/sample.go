package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/xshoji/jsonbox-go/jsonboxgo"
	"log"
	"math/rand"
	"net/http"
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
	client := jsonboxgo.NewClient(baseUrl, boxId, http.DefaultClient)

	createUser := func() User {
		return User{
			Name:     "taro_" + randomString(),
			Age:      randomNumber(),
			Language: "JP",
		}
	}

	// Create
	_ = client.Create(collection, createUser())
	_ = client.Create(collection, createUser())
	_ = client.Create(collection, createUser())

	// ReadByQuery
	result := client.ReadByQuery(
		collection,
		jsonboxgo.
			NewQueryBuilder().
			Limit(3).
			Offset(1).
			SortAsc("age").
			AddGreaterThanOrEqual("age", "40"),
	)
	fmt.Println(string(result))
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
