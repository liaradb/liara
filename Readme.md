# LiaraDB

**_Event Sourcing database_**

Event-native database, to power Event Sourced workflows for Microservices with Domain Driven Design.

## Run

```bash
cd ./services/liaradb
go run ./cmd
```

## Make targets

Go to `./services/liaradb`.

| Name        | Command            | Notes                                     |
| ----------- | ------------------ | ----------------------------------------- |
| Test        | `make test`        | Run tests                                 |
| Cover       | `make cover`       | Run tests with code coverage              |
| Report      | `make report`      | Run tests with coverage report            |
| Report file | `make report-file` | Run tests and emit coverage report file   |
| Build       | `make build`       | Build executable                          |
| Doc         | `make doc`         | View documentation site                   |
| Clean       | `make clean`       | Removes build and test directories        |
| Clean build | `make clean-build` | Removes build directory                   |
| Clean test  | `make clean-test`  | Removes test directory                    |
| Clean data  | `make clean-data`  | Removes data directory                    |
| Clean all   | `make clean-all`   | Removes build, test, and data directories |

## Update gRPC

```bash
cd ./modules/protos
make
```
