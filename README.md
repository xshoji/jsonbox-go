## jsonbox-go

jsonbox-go is the wrapper library in order to use jsonbox as easy.

> jsonbox.io ｜ A Free HTTP based JSON storage  
> https://jsonbox.io/

## Usage

#### Create new client

```
baseUrl := "https://jsonbox.io/"
boxId := "box_xxxxxxxxxx"
client := jsonboxgo.NewClient(baseUrl, boxId)
```

#### Create object

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
```

#### Read all records

```
result := client.ReadAll(collection)
```

#### Read one by recordId

```
collection := "users"
result, found := client.Read(collection, user.Id)
```

#### Update

```
user.Name = "updated"
user.Age = 24
collection := "users"
result, updated := client.Update(collection, user.Id, user)
```

#### Delete

```
result, deleted := client.Delete(collection, createdUser.Id)
```

## Demo

```
go run cmd/sample/sample.go
```

