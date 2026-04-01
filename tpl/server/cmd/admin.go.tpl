package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/xpfo-go/logs"
	"<xpfo{ .ModulePath }xpfo>/internal/api"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
	"<xpfo{ .ModulePath }xpfo>/internal/server"

<xpfo{ if .EnableMySQL }xpfo>
	_ "github.com/go-sql-driver/mysql"
<xpfo{ end }xpfo>
)

func init() {
	rootCmd.Flags().StringP("config", "c", "./config.yaml", "config file (default is config.yml;required)")
}

var rootCmd = &cobra.Command{
	Use:   "<xpfo{ .ProjectName }xpfo>",
	Short: "<xpfo{ .ProjectName }xpfo> Backend",
	Long:  "<xpfo{ .ProjectName }xpfo> Backend, Code by Go",
	RunE: func(cmd *cobra.Command, args []string) error {
		return Start(cmd)
	},
}

func Start(cmd *cobra.Command) error {
	// 1. init
	// 初始化配置文件
	if err := initConfig(cmd, "config"); err != nil {
		return err
	}
	// 初始化日志
	if err := initLogs(); err != nil {
		return err
	}
	// 初始化数据库
	if err := initDatabase(); err != nil {
		return err
	}
	// 初始化缓存
	if err := initRedis(); err != nil {
		return err
	}

	// 2. watch the signal
	ctx, cancelFunc := context.WithCancel(context.Background())

	// 3. start the server
	httpServer := server.NewServer(config.Configor, api.NewRouter)
	go httpServer.Run(ctx)

	interrupt(cancelFunc)
	return nil
}

func Execute() error {
	return rootCmd.Execute()
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
