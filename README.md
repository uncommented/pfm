# Prerequisites

```
brew install protobuf

go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

# Environment variables

Setup environment variables with `.env` file.

```shell
# KIS
KIS_APPKEY=<appkey>
KIS_APPSECRET=<appsecret>
KIS_TOKEN=<token>
KIS_TOKEN_EXPIRED="2024-dd-yy HH:MM:SS"

CANO=12345678
ACNT_PRDT_CD=01

# Upbit
UPBIT_ACCESS_KEY=<accesskey>
UPBIT_SECRET_KEY=<secretkey>
```
