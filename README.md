# ZBX1

#### What is this ?

Coming Soon...

#### How to use it ?

Coming Soon...

#### Set up

1. Install dependencies
```bash
go mod tidy
```
2. To create RSA keys
```bash
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

#### Running
```bash
go run main.go
```

#### Running with docker compose
```bash
docker-compose up -d --build --force-recreate
```

#### Running tests
```bash
cd tests
go test -race -v ./...
```