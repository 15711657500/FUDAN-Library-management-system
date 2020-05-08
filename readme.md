# FUDAN Library Management System

## Requirements

```
go 1.14
github.com/go-sql-driver/mysql v1.5.0
github.com/jmoiron/sqlx v1.2.0
github.com/modood/table v0.0.0-20200225102042-88de94bb9876
```

## Folder Structure
```
src/   go source files
+-- library.go main func
+-- library_test.go for tests
+-- file.go read data from csv file
+-- rent.go funcs about borrowing, querying
+-- book.go funcs about managing books
+-- users.go funs about manaing users

data/  csv files for inserting information into database
```

## Setup

First, open "src/library.go" and modify the mysql account configuration.

```go
const (
	USER     = "root"
	Password = "xxx"
	DBName   = "ass3"
)
```

For tests, run command

```
go test
```

To get the library management system with command line interface, run command

```
make
./library
```