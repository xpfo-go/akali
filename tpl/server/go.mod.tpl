module <xpfo{ .ModulePath }xpfo>

go <xpfo{ .GoVersion }xpfo>

require (
	github.com/gin-gonic/gin v1.9.1
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/jinzhu/configor v1.2.2
	github.com/spf13/cobra v1.8.0
	github.com/xpfo-go/logs v0.0.0-20240503134654-f538ed2b19af
	go.uber.org/automaxprocs v1.5.3
<xpfo{ if .EnableMySQL }xpfo>	github.com/go-sql-driver/mysql v1.6.0
	github.com/jmoiron/sqlx v1.3.5
	github.com/DATA-DOG/go-sqlmock v1.5.2
	github.com/stretchr/testify v1.8.3
<xpfo{ end }xpfo><xpfo{ if .EnableRedis }xpfo>	github.com/redis/go-redis/v9 v9.5.1
<xpfo{ end }xpfo><xpfo{ if .EnableSwagger }xpfo>	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.3
<xpfo{ end }xpfo><xpfo{ if .EnableMetrics }xpfo>	github.com/prometheus/client_golang v1.19.0
<xpfo{ end }xpfo><xpfo{ if and .EnableMySQL .EnableMetrics }xpfo>	github.com/dlmiddlecote/sqlstats v1.0.2
<xpfo{ end }xpfo>)
