package database

import (
	"github.com/dlmiddlecote/sqlstats"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"
	"{{ .ProjectName }}/config"
	"sync"
)

var (
	DefaultDBClient *DBClient
)

var (
	defaultDBClientOnce sync.Once
)

// InitDBClients ...
func InitDBClients(defaultDBConfig *config.Config) {
	if DefaultDBClient == nil {
		defaultDBClientOnce.Do(func() {
			DefaultDBClient = NewDBClient(defaultDBConfig)
			if err := DefaultDBClient.Connect(); err != nil {
				panic(err)
			}

			// https://github.com/dlmiddlecote/sqlstats
			// Create a new collector, the name will be used as a label on the metrics
			collector := sqlstats.NewStatsCollector(defaultDBConfig.Database.Name, DefaultDBClient.DB)
			// Register it with Prometheus
			prometheus.MustRegister(collector)
		})
	}

}

// GetDefaultDBClient 获取默认的DB实例
func GetDefaultDBClient() *DBClient {
	return DefaultDBClient
}

// GenerateDefaultDBTx 生成一个事务链接
func GenerateDefaultDBTx() (*sqlx.Tx, error) {
	return GetDefaultDBClient().DB.Beginx()
}
