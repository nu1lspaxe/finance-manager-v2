# Finance Manager (v2)

Different from [finance-manager](https://github.com/nu1lspaxe/finance-manager), the finance-manager-v2 is designed with microservice pattern and combined the utilities with Kafka and gRPC.

## Todo List
- [ ] Clickhouse CDC from Postgres (Kafka + Debezium)
- [ ] ETL
- [ ] Binance API 
- [ ] Microservice single-sign on for each service database 
  - After a user sign in, server has to store token into corresponding database by the services which are authorized to the user
  - Need to ensure token synchronicity

## CA Certificates

The self-signed certificates are used in all services.

1. We use `mkcert` to generate certificates, please follow the steps:

  ```bash
  # 1. In WSL2, get WSL2 ip address by command `ip a`
  # 2. Find eth0 or other interface, it should contains `inet 172.XXX.XXX.XXX/20`
  # 3. Run the command in powershell, replace <wsl_ip_addr> with 172.XXX.XXX.XXX/20
  # 4. <service_name> stands for bank_syste, records, users, records_bank in this project
  mkcert localhost 127.0.0.1 ::1 <wsl_ip_addr> <service_name>

  ### Explanation
  # - <wsl_ip_addr> used to allow our (frontend) flutter communicate with our services
  # - <service_name> used as domain name to allow our services communicate within docker environment
  ```

2. Two files `localhost+4.pem` and `localhost+4-key.pem`　you will see after running step 1. Try to rename these two file with `cert.pem` and `key.pem` and move them into `<service_name>/certs` 

3. In powershell, run `mkcert -CAROOT` to get ca-root file location. There should have two files: `rootCA.pem` and `rootCA-key.pem` within the given path. Rename `rootCA.pem` into `ca.pem` then copy the file into each `<service_name>/certs`

## gRPC CLI

We use [Evans CLI](https://github.com/ktr0731/evans) for interacting with gRPC server.

```bash
evans -r repl -p <grpc_port> --tls --cacert certs/cert.pem

$ show package 
$ package <package_name>  # set package connection
```

## Bank System

> Supporting protocols : HTTP/1.1

The service is responsible for mocking a digital bank in real world. The app (finance-manager), in the version 1 stage, will retrieve data from bank_system on a daily basis.

### Cronjob

Bank system will run 3 cornjobs once the server starting in order to simulate real world user-bank interaction.

- Add a new user with a new account, then run deposit transaction
- Add a new account to existing users, then run deposit transaction
- Add a withdrawal transaction to existing accounts

### Users

- `[POST] /v1/user`
- `[GET] /users/{id}`
- `[GET] /users/{id}/accounts`
- `[GET] /users`
- `[PUT] /users/{id}`

### Accounts

- `[POST] /accounts`
- `[GET] /accounts/{id}`
- `[GET] /accounts/{id}/balance`
- `[GET] /accounts`

### Transactions

- `[GET] /transactions/{id}`

## Services

### Users

> Supporting protocols : HTTP/1.1, GRPC (HTTP/2)

- `[POST] /users/signin`
- `[POST] /users/logout`
- `[POST] /users/signup`
- `[GET] /v1/users/{id}`
- `[GET] /v1/users`
- `[Patch] /v1/users/{id}`
- `[Delete] /v1/users/{id}`

### Records

> Supporting protocols : HTTP/1.1, GRPC (HTTP/2)

- `[POST] /v1/records`
- `[GET] /v1/records/{id}`
- `[GET] /v1/records/user`
  - user_id
  - record_type
  - start_time
  - end_time
- `[PATCH] /v1/records/{id}`
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

### Analysis