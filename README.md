# Finance Manager (v2)

Different from [finance-manager](https://github.com/nu1lspaxe/finance-manager), the finance-manager-v2 is designed with microservice pattern and combined the utilities with Kafka and gRPC.


## Bank System
> Supporting protocols : HTTP/1.1

The service is responsible for mocking a digital bank in real world. The app (finance-manager), in the version 1 stage, will retrieve data from bank_system on a daily basis. 

### Cronjob
Bank system will run 3 cornjobs once the server starting in order to simulate real world user-bank interaction. 

- Add a new user with a new account, then run deposit transaction
- Add a new account to existing users, then run deposit transaction
- Add a withdrawal transaction to existing accounts

### Users

#### CreateUser
- `[POST] /users`

#### GetUserByID 
- `[GET] /users/{id}`

#### GetUserAccounts
- `[GET] /users/{id}/accounts`

#### GetAllUsers
- `[GET] /users`

#### UpdateUser
- `[PUT] /users/{id}`

### Accounts

#### CreateAccount
- `[POST] /accounts`

#### GetAccountByID
- `[GET] /accounts/{id}`

#### GetAccountBalance
- `[GET] /accounts/{id}/balance`

#### GetAllAccounts
- `[GET] /accounts`

### Transactions

#### GetTransactionByID
- `[GET] /transactions/{id}`


## Services

### Users

> Supporting protocols : HTTP/1.1, GRPC (HTTP/2)

#### CreateUser
- `[POST] /v1/user` 

#### GetUser
- `[GET] /v1/users/{id}`

#### GetAllUsers
- `[GET] /v1/users` 

#### UpdateUser
- `[Patch] /v1/users/{id}` 

#### DeleteUser
- `[Delete] /v1/users/{id}` 

### Records

> Supporting protocols : HTTP/1.1, GRPC (HTTP/2)

#### CreateRecord
- `[POST] /v1/records`

#### GetRecord
- `[GET] /v1/records/{id}`

#### GetUserRecordsWithFilters
- `[GET] /v1/records/user`
- parameters: 
  - user_id
  - record_type
  - start_time
  - end_time

#### UpdateRecord
- `[PATCH] /v1/records/{id}`

#### DeleteRecord
- `[DELETE] /v1/records/{id}`

### Bank Records

The workflow will basically show like below :

```
FM_Users --> BS : Request user binding accounts in finance-manager from bank-system
FM_Users <-- BS : Retrieve binding accounts data to update account balance
FM_Users --> Kafka (update_record) : Publish account_id to topic (update_record)
Kafka (update_record) --> FM_Records_Bank : Subscribe topic (update_record) to update bank records
FM_Records_Bank --> BS : Request account transactions
FM_Records_Bank <-- BS : Retrieve transactions then update records
```

### Exporter
