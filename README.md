# ZBX1

#### What is this ?

RESTful web service to store files

#### How to use it ?

https://zbx1-dashboard.herokuapp.com

#### Set up

```bash
# Clone the repository
git clone https://github.com/metiago/zbx1

# Enter in directory
cd zbx1

# Install required modules
go mod tidy

# Create RSA keys
openssl genrsa -out app.rsa 1024
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

#### Running locally
```bash
./local_run.sh
```

#### Running on docker compose
```bash
./compose_run.sh
```

#### Running integration tests
```bash
cd tests
go test -race -v ./...
```