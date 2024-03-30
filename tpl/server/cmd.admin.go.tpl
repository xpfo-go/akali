package cmd

import (
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"{{ .ProjectName }}/config"
	"{{ .ProjectName }}/internal/api"
	"{{ .ProjectName }}/pkg/logs"
	"{{ .ProjectName }}/pkg/server"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
)

func init() {
	rootCmd.Flags().StringP("config", "c", "./config.yaml", "config file (default is config.yml;required)")
}

var rootCmd = &cobra.Command{
	Use:   "{{ .ProjectName }}",
	Short: "{{ .ProjectName }} Backend",
	Long:  "{{ .ProjectName }} Backend, Code by Go",
	Run: func(cmd *cobra.Command, args []string) {
		Start(cmd)
	},
}

func Start(cmd *cobra.Command) {
	// 1. init
	// 初始化配置文件
	initConfig(cmd, "config")
	// 初始化日志
	initLogs()
	// 初始化数据库
	initDatabase()

	// 2. watch the signal
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 3. start the server
	httpServer := server.NewServer(config.Configor, api.NewRouter)
	go httpServer.Run(ctx)

	interrupt(cancelFunc)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func interrupt(onSignal func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	for s := range c {
		logs.Infof("Caught signal %s. Exiting.", s)
		onSignal()
		close(c)
	}
}
