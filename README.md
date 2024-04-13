# Instructions
## Set up
1. git clone repository
2. run `go run main.go`

## Endpoints
Server will be running on `localhost:3000`
1. `CreateAccount` -> `POST /accounts`

```
body: {
    "account_id": 2,
    "initial_balance": "3000.23"
}
```

2. `QueryAccount` -> `GET /accounts/{account_id}`
3. `CreateTransaction` -> `POST /transactions`

```
body: {
    "source_account_id": 1,
    "destination_account_id": 3,
    "amount": "49.9"
}
```

## Assumptions
1. Assume edge cases fully covered, so no unit tests, hence did not ensure all dependencies are mockable
2. Assume unnecessary need for logging, so no logging implementation
3. Assume no need for metrics, so no metrics implementation
4. Assume transaction log should only be recorded in DB upon success, failed transactions don't have any record
5. Assume no need for data schema relation for small scale API, so used MongoDB for flexibility
6. Assume account updates and transaction logging is atomic during transaction creation, so provided session based db writes and rollback
7. Assume no security loopholes, so no auth provided on request, and db URI (free account so doesn't matter) is exposed in repo as well
