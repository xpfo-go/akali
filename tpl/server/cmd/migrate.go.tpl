package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"
	"<xpfo{ .ModulePath }xpfo>/internal/config"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "run database migrations",
}

var migrateUpCmd = &cobra.Command{
	Use:   "up",
	Short: "apply all up migrations",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig(cmd, "config")
		m, err := newMigrator()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "rollback one migration",
	Run: func(cmd *cobra.Command, args []string) {
		initConfig(cmd, "config")
		m, err := newMigrator()
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		if err := m.Steps(-1); err != nil && err != migrate.ErrNoChange {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	migrateCmd.PersistentFlags().StringP("config", "c", "./config.yaml", "config file path")
	migrateCmd.AddCommand(migrateUpCmd, migrateDownCmd)
	rootCmd.AddCommand(migrateCmd)
}

func newMigrator() (*migrate.Migrate, error) {
	cfg := config.Configor.Mysql
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&multiStatements=true&loc=%s&time_zone=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, "UTC", url.QueryEscape("'+00:00'"),
	)
	return migrate.New("file://migrations", "mysql://"+dsn)
}
