## jsonbox-go

jsonbox-go is the wrapper library in order to use jsonbox.io as easy.

> jsonbox.io ｜ A Free HTTP based JSON storage  
> https://jsonbox.io/

## Usage

```
go get "github.com/xshoji/jsonbox-go"
```

#### Create new client

```go
baseUrl := "https://jsonbox.io/"
boxId := "box_xxxxxxxxxx"
client := jsonboxgo.NewClient(baseUrl, boxId)
```

## CRUD operation

#### Create record

```go
user := struct {
	Id        string `json:"_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Age       int    `json:"age,omitempty"`
	CreatedOn string `json:"_createdOn,omitempty"`
}{
	Name: "taro",
	Age: 100,
}
collection := "users"
result := client.Create(collection, user)
fmt.Println(string(result))
// {
//   "_id": "5ea6e8c543f5c4001710132b",
//   "name": "taro",
//   "age": 100,
//   "_createdOn": "2020-04-27T14:14:29.843Z"
// }
```

#### Read all records

```go
collection := "users"
result := client.ReadAll(collection)
fmt.Println(string(result))
// [
//   {
//     "_id": "5ea6e8c543f5c4001710132b",
//     "name": "taro",
//     "age": 100,
//     "_createdOn": "2020-04-27T14:14:29.843Z"
//   }
// ]
```

#### Read one by recordId

```go
collection := "users"
result, found := client.Read(collection, user.Id)
fmt.Println(string(result))
// {
//   "_id": "5ea6e8c543f5c4001710132b",
//   "name": "taro",
//   "age": 100,
//   "_createdOn": "2020-04-27T14:14:29.843Z"
// }
```

#### Update

```go
user.Name = "updated"
user.Age = 24
collection := "users"
result, updated := client.Update(collection, user.Id, user)
fmt.Println(string(result))
// {
//   "_id": "5ea6e8c543f5c4001710132b",
//   "name": "updated",
//   "age": 24,
//   "_createdOn": "2020-04-27T14:14:29.843Z",
//   "_updatedOn": "2020-04-27T14:14:31.010Z"
// }
```

#### Delete

```go
collection := "users"
result, deleted := client.Delete(collection, user.Id)
fmt.Println(string(result))
// {"message":"Record removed."}
```

## Read by query operation

```go
collection := "users"
result := client.ReadByQuery(
	collection,
	jsonboxgo.
		NewQueryBuilder().
		Offset(1).
		Limit(3).
		SortAsc("age").
		AddEqual("country", "JP").
		AddGreaterThanOrEqual("age", "40")
)
fmt.Println(string(result))
// [
//   {
//     "_id": "5ea9bc0225ec0a0017640226",
//     "name": "taro_95ad1144ec3fcfb628234f6eeddc8e9fde02126c1ab0b0f163a79a3d8910c666",
//     "country": "JP",
//     "age": 40,
//     "_createdOn": "2020-04-29T17:40:18.986Z"
//   },
//   {
//     "_id": "5ea9bbec25ec0a0017640222",
//     "name": "taro_e3d68d50f4399a0b21df5ae4ca71ad932e3d87f4c274d2e3334350f47b2c887b",
//     "country": "JP",
//     "age": 44,
//     "_createdOn": "2020-04-29T17:39:56.906Z"
//   },
//   {
//     "_id": "5ea9bce825ec0a001764022a",
//     "name": "taro_ccd90e826262cffb6c5e6be6993c7752da008a16f6811ba53f3aad19e9cc54d1",
//     "country": "JP",
//     "age": 44,
//     "_createdOn": "2020-04-29T17:44:08.839Z"
//   }
// ]
```

## Test

```
go test -v ./...
```

## Sample

```
go run cmd/sample/sample.go
```
