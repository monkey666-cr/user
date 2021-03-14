# User Service

This is the User service

Generated with

```
micro new user
```

## Usage

Generate the proto code

```
make proto
```

Run the service

```
micro run .
```

## go 交叉编译

```shell
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o *.go
```

## proto 生成代码

```shell
protoc -I ./ --go_out=./ --micro_out=./ ./*.proto
```

## docker 编译

```shell
docker build -t user:latest .
```