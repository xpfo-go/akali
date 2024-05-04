package cmd

import (
	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ProjectName }xpfo>/internal/config"
	"<xpfo{ .ProjectName }xpfo>/internal/database"
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
	database.InitDBClients(&config.Configor.Mysql)
}
