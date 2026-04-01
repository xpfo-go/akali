package database

import (
	"fmt"
	"sync"

<xpfo{ if .EnableMetrics }xpfo>
	"github.com/dlmiddlecote/sqlstats"
	"github.com/prometheus/client_golang/prometheus"
<xpfo{ end }xpfo>
	"github.com/jmoiron/sqlx"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
)

var (
	DefaultDBClient *DBClient
)

var (
	defaultDBClientOnce sync.Once
	defaultDBClientErr  error
)

// InitDBClients ...
func InitDBClients(defaultDBConfig *config.MysqlConfig) error {
	if DefaultDBClient == nil {
		defaultDBClientOnce.Do(func() {
			client := NewDBClient(defaultDBConfig)
			if err := client.Connect(); err != nil {
				defaultDBClientErr = fmt.Errorf("connect mysql failed: %w", err)
				return
			}
			DefaultDBClient = client

<xpfo{ if .EnableMetrics }xpfo>
			// https://github.com/dlmiddlecote/sqlstats
			// Create a new collector, the name will be used as a label on the metrics
			collector := sqlstats.NewStatsCollector(defaultDBConfig.Name, client.DB)
			// Register it with Prometheus
			if err := prometheus.Register(collector); err != nil {
				if _, ok := err.(prometheus.AlreadyRegisteredError); !ok {
					defaultDBClientErr = fmt.Errorf("register mysql metrics failed: %w", err)
					return
				}
			}
<xpfo{ end }xpfo>
		})
	}

	return defaultDBClientErr
}

// GetDefaultDBClient 获取默认的DB实例
func GetDefaultDBClient() *DBClient {
	return DefaultDBClient
}

// GenerateDefaultDBTx 生成一个事务链接
func GenerateDefaultDBTx() (*sqlx.Tx, error) {
	return GetDefaultDBClient().DB.Beginx()
}
