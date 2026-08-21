# API Definition

## API Definition (User)

### Create User

#### Request
```http
POST api/v1/users
```

```json
{
  "username": "mohit_dev",
  "password": "************"
}
```

#### Response
```
201 Created
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "username": "mohit_dev"
}
```
**Headers**
```http
Set-Cookie: token=<token_value>; HttpOnly; Secure
```

## API Definition (Authentication)

### Sign in

#### Request
```http
POST api/v1/auth/signIn
```

```json
{
  "username": "mohit_dev",
  "password": "************"
}
```

#### Response
```
200 OK
```

**Headers**
```http
Set-Cookie: token=<token_value>; HttpOnly; Secure
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "username": "mohit_dev"
}
```

### Sign out

#### Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
POST api/v1/auth/signOut
```

#### Response
```
204 No Content
```

**Headers**
```http
Set-Cookie: token=; HttpOnly; Secure; Max-Age=0
```

## API Definition (Expense)

### Create Expense

#### Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
POST api/v1/users/{{userId}}/expenses
```

```json
{
  "description": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00Z"
}
```

#### Response
```
201 Created
```
```http
Location: {{host}}/api/v1/users/{{userId}}/expenses
```
```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "description": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00Z"
}
```

### Get Expenses

#### Get Bulk Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
GET api/v1/users/{{userId}}/expenses?cursor={base64_string_from_previous_result}&limit={limit}&sortField={sortField}&filterField={filterField}&filterValue={filterValue}&sortOrder={sortOrder}
```

#### Response
```
200 OK
```

```json
{
  "expenses": [
    {
      "id": "286d7bbf-e6e0-4bfd-b4e0-906a613193db",
      "description": "Car Repair",
      "amount": 350.5,
      "date": "2024-02-15T08:00:00Z"
    },
    {
      "id": "3f91f017-af32-46b3-9c53-47adb1314c9a",
      "description": "Groceries",
      "amount": 75.3,
      "date": "2024-03-10T08:00:00Z"
    }
  ],
  "cursor": "base64_string"
}
```

#### Get One Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
GET api/v1/users/{{userId}}/expenses/{{id}}
```

#### Response
```
200 OK
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "description": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00Z"
}
```

### Update Expense

#### Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
PUT api/v1/users/{{userId}}/expenses/{{id}}
```

```json
{
  "id": "00000000-0000-0000-0000-000000000000",
  "description": "Groceries",
  "amount": 279.7,
  "date": "2024-06-08T08:00:00Z"
}
```

#### Response
```
204 No Content
```
or
```
201 Created
```
```http
Location: {{host}}/api/v1/users/{{userId}}/expenses/{{id}}
```

### Delete Expense

#### Request
**Headers**
```http
Cookie: token=<token_value>
```

```http
DELETE api/v1/users/{{userId}}/expenses/{{id}}
```

#### Response
```
204 No Content
```