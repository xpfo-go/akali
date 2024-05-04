# akali — A CLI tool for building Go applications.

## Features
- **Gin**: https://github.com/gin-gonic/gin
- **Sqlx**: https://github.com/jmoiron/sqlx
- **Zap**: https://github.com/uber-go/zap
- **Swaggo**:  https://github.com/swaggo/swag
- More...

### Create a New Project

You can create a new Go project with the following command:

```bash
akali create [ProjectName]
```

## Directory Structure
```
.
├── bin
├── cmd
│   ├── admin.go
│   ├── init.go
│   └── version.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── api
│   │   ├── basic
│   │   ├   └── router.go
│   │   └── router.go
├   ├── config
│   │   └── config.go
│   ├── controller
│   │   └── basic
│   │       ├── basic.go
│   │       └── health.go
│   ├── database
│   │   ├── dao
│   │   ├── do
│   │   ├── entity
│   │   ├── dbmock.go
│   │   ├── init.go
│   │   ├── mysql.go
│   │   ├── sqlx.go
│   │   ├── sqlx_helper.go
│   │   ├── sqlx_helper_test.go
│   │   ├── sqlx_test.go
│   │   └── utils.go
│   ├── middleware
│   │   └── request_id.go
│   │── server
│   │   ├── router.go
│   │   └── server.go
│   │── service
│   │── task
│   │── util
│   │   ├── consts.go
│   │   ├── request.go
│   │   ├── slice.go
│   │   ├── string.go
│   │   └── uuid.go
│   └── version
│       └── version.go
├── config.yaml
├── main.go
├── makefile
├── go.mod
└── go.sum

```

## License

Akali is released under the MIT License. For more information, see the [LICENSE](LICENSE) file.
