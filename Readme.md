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

| Target        | Notes                                     |
| ------------- | ----------------------------------------- |
| `test`        | Run tests                                 |
| `cover`       | Run tests with code coverage              |
| `report`      | Run tests with coverage report            |
| `report-file` | Run tests and emit coverage report file   |
| `build`       | Build executable                          |
| `doc`         | View documentation site                   |
| `clean`       | Removes build and test directories        |
| `clean-build` | Removes build directory                   |
| `clean-test`  | Removes test directory                    |
| `clean-data`  | Removes data directory                    |
| `clean-all`   | Removes build, test, and data directories |

## Update gRPC

```bash
cd ./modules/protos
make
```
