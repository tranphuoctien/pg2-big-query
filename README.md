# Postgres stream to bigquery
Listen and streaming data to big query

## Getting Started

Provide package for logger

## Use

Big Query API [Quick start](https://cloud.google.com/bigquery/docs/quickstarts/quickstart-client-libraries) to get json file

After that then can replace content into folder bgs with file ***sk.json***

Then you can change your schema on line 100 in ***example/main.go*** to match what you need.

```
go run cmd/main.go -tables .* -connect "postgresql://postgres:pass@localhost/database?sslmode=disable" -v true
```

### Prerequisites

Golang version 1.13 above

```
go version
```
To check your localy version

Runing
```
go test
```


## Versioning

This project is the current develop

## Authors

* **Tien TP** - *Initial work* - [TienTP](https://g.ghn.vn/tientp)

See also the list of [contributors](https://g.ghn.vn/logistic/bi/streaming/pg2-big-query/master) who participated in this project.
