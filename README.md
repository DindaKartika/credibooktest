# CredibookTest
This simple web service is created for the completion of Credibook Technical Assesment.

## Databases
All datas in this web service saved in:
- MySQL; used as the main database in this service

## APIs
### Login
This endpoint used for user login
##### Endpoint
```sh
[POST] /login
```
##### Request body
```sh
{
	"username": "admin",
	"password": "admin"
}
```

### User
#### Get User
This endpoint used for getting all user list. This API can only be accessed by Admin
##### Endpoint
```sh
[GET] /users
```
##### Auth
User auth for this endpoint is bearer token using token got from the login API
##### Response Success
```sh
{
  "total_records": 3,
  "data": [
    {
      "id": 1,
      "created_at": "2021-06-11T15:49:19Z",
      "updated_at": "2021-06-11T15:49:19Z",
      "deleted_at": null,
      "username": "admin",
      "is_admin": true
    }
  ],
  "current_page": 1,
  "total_pages": 1
}
```

#### Create New User
This endpoint used for creating new user
##### Endpoint
```sh
[GET] /user
```
##### Request body
```sh
{
	"username": "User",
	"password": "My Pass"
}
```
##### Response Success
```sh
{
  "New user created"
}
```

### Transaction
#### API Get All Transaction
This endpoint used to get all transactions saved in the database.
##### Endpoint
```sh
[GET] /transaction
```
##### Auth
User auth for this endpoint is bearer token using token got from the login API
##### Filters
This API provided 3 filters in the query param :
- page
- perpage
- type; will show transaction with types as in the filter
- min_amount; will show transaction in with amount higher than or equal as in the filter
- max_amount; will show transaction in with amount lower than or equal as in the filter
- order_by; can be used to sort the data, please add ' DESC' if it's sorted descending

##### Response success
```sh
{
  "total_records": 3,
  "data": [
    {
      "id": 1,
      "created_at": "2021-06-14T15:21:42Z",
      "updated_at": "2021-06-14T15:21:42Z",
      "deleted_at": null,
      "user_id": 2,
      "amount": 100000,
      "notes": "Kulakan",
      "type": "expense"
    },
    {
      "id": 2,
      "created_at": "2021-06-14T15:22:17Z",
      "updated_at": "2021-06-14T15:22:17Z",
      "deleted_at": null,
      "user_id": 3,
      "amount": 100000,
      "notes": "Hasil penjualan",
      "type": "income"
    }
  ],
  "current_page": 1,
  "total_pages": 1
}
```

#### API Add New Transaction
This endpoint used to create new transaction and save it to the database.
##### Endpoint
```sh
[POST] /transaction
```
##### Auth
User auth for this endpoint is bearer token using token got from the login API
##### Request body
```sh
{
	"amount": 10000, //required
	"notes": "Groceries",
	"type": "expense"
}
```
##### Response success
```sh
{
  "id": 4,
  "created_at": "2021-06-14T23:18:14.65788+07:00",
  "updated_at": "2021-06-14T23:18:14.65788+07:00",
  "deleted_at": null,
  "user_id": 2,
  "amount": 10000,
  "notes": "Groceries",
  "type": "expense"
}
```

#### API Update Transaction
This endpoint used to update existing transaction and save it to the database.
##### Endpoint
```sh
[PUT] /transaction/:id
```
##### Auth
User auth for this endpoint is bearer token using token got from the login API
##### Request body
```sh
{
	"amount": 150000, //required
	"notes": "Kulakan",
	"type": "Expense"
}
```
##### Response success
```sh
{
  "id": 4,
  "created_at": "2021-06-14T23:18:14.65788+07:00",
  "updated_at": "2021-06-14T23:18:14.65788+07:00",
  "deleted_at": null,
  "user_id": 2,
  "amount": 150000,
  "notes": "Kulakan",
  "type": "expense"
}
```

#### API Delete Transaction
This endpoint used to delete transaction and save it to the database.
##### Endpoint
```sh
[DELETE] /transaction/:id
```
##### Auth
User auth for this endpoint is bearer token using token got from the login API
##### Response success
```sh
"transaction has deleted successfully"
```

## Testing
In this service, there's testing provided. To do the testing, you can use command below:
```sh
go test ./... -coverprofile=coverage.out
```