package cmd

import (
	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
<xpfo{ if .EnableMySQL }xpfo>	"<xpfo{ .ModulePath }xpfo>/internal/database"
<xpfo{ end }xpfo><xpfo{ if .EnableRedis }xpfo>	"<xpfo{ .ModulePath }xpfo>/internal/cache"
<xpfo{ end }xpfo>
)

func initConfig(cmd *cobra.Command, key string) {
	configFilePath, err := cmd.Flags().GetString(key)

	if err != nil {
		panic(err.Error())
	}
	if err := configor.Load(config.Configor, configFilePath); err != nil {
		panic(err.Error())
	}
}

func initLogs() {
	conf := logs.GetLogConf()
	conf.FileName = config.Configor.Log.FileName
	conf.MaxAge = config.Configor.Log.MaxAge
	conf.Level = config.Configor.Log.Level
	logs.InitLogSetting(conf)
}

func initDatabase() {
<xpfo{ if .EnableMySQL }xpfo>
	database.InitDBClients(&config.Configor.Mysql)
<xpfo{ end }xpfo>
}

func initRedis() {
<xpfo{ if .EnableRedis }xpfo>
	cache.InitRedis(&config.Configor.Redis)
<xpfo{ end }xpfo>
}
