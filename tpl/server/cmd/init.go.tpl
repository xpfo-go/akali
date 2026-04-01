package cmd

import (
	"fmt"

	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
<xpfo{ if .EnableMySQL }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/database"
<xpfo{ end }xpfo>
<xpfo{ if .EnableRedis }xpfo>
	"<xpfo{ .ModulePath }xpfo>/internal/cache"
<xpfo{ end }xpfo>
)

func initConfig(cmd *cobra.Command, key string) error {
	configFilePath, err := cmd.Flags().GetString(key)
	if err != nil {
		return fmt.Errorf("read %s flag failed: %w", key, err)
	}
	if err := configor.Load(config.Configor, configFilePath); err != nil {
		return fmt.Errorf("load config failed: %w", err)
	}
	return nil
}

func initLogs() error {
	conf := logs.GetLogConf()
	conf.FileName = config.Configor.Log.FileName
	conf.MaxAge = config.Configor.Log.MaxAge
	conf.Level = config.Configor.Log.Level
	logs.InitLogSetting(conf)
	return nil
}

func initDatabase() error {
<xpfo{ if .EnableMySQL }xpfo>
	return database.InitDBClients(&config.Configor.Mysql)
<xpfo{ end }xpfo>
	return nil
}

func initRedis() error {
<xpfo{ if .EnableRedis }xpfo>
	return cache.InitRedis(&config.Configor.Redis)
<xpfo{ end }xpfo>
	return nil
}
