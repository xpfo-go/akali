package cmd

import (
	"github.com/jinzhu/configor"
	"github.com/spf13/cobra"
	"{{ .ProjectName }}/config"
	"{{ .ProjectName }}/internal/database"
	"{{ .ProjectName }}/pkg/logs"
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
	logs.InitLogSetting(&logs.LogConfig{Level: config.Configor.Log.Level, MaxAge: 21, AppName: "admin"})
}

func initDatabase() {
	database.InitDBClients(config.Configor)
}
