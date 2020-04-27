## jsonbox-go

jsonbox-go is the wrapper library in order to use jsonbox as easy.

> jsonbox.io ｜ A Free HTTP based JSON storage  
> https://jsonbox.io/

## Usage

```
go get "github.com/xshoji/jsonbox-go"
```

#### Create new client

```
baseUrl := "https://jsonbox.io/"
boxId := "box_xxxxxxxxxx"
client := jsonboxgo.NewClient(baseUrl, boxId)
```

#### Create record

```
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

```
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

```
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

```
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

```
collection := "users"
result, deleted := client.Delete(collection, user.Id)
fmt.Println(string(result))
// {"message":"Record removed."}
```

## Test

```
go test -v ./...
```

## Sample

```
go run cmd/sample/sample.go
```
