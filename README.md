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
│   └── init.go
├── config
│   └── config.go
├── docs
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal
│   ├── api
│   │   ├── basic
│   │   ├   └── router.go
│   │   └── router.go
│   ├── controller
│   │   └── basic
│   │       ├── basic.go
│   │       └── health.go
│   ├── database
│   │   ├── dao
│   │   ├── do
│   │   ├── entity
│   │   ├── init.go
│   │   └── mysql.go
│   └── version
│       └── version.go
├── pkg
│   ├── limiter
│   │   ├── dcs_leaky_bucket.go
│   │   ├── dcs_token_bucket.go
│   │   ├── leaky_bucket.go
│   │   └── token_bucket.go
│   ├── logs
│   │   └── logs.go
│   ├── retry
│   │   └── retry.go
│   └── server
│       ├── router.go
│       └── server.go
├── config.yaml
├── main.go
├── makefile
├── go.mod
└── go.sum

```

## License

Akali is released under the MIT License. For more information, see the [LICENSE](LICENSE) file.
