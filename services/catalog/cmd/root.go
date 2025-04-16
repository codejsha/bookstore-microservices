package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/codejsha/bookstore-microservices/catalog/internal/di"
	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
)

const serviceName = "catalog"

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(migrateCmd)

	runCmd.Flags().StringP("env", "e", string(object.Local), "set environment")
	migrateCmd.Flags().StringP("env", "e", string(object.Local), "set environment")

	logrus.SetFormatter(&logrus.JSONFormatter{DisableHTMLEscape: true})
	logrus.SetOutput(os.Stdout)
}

var rootCmd = &cobra.Command{
	Use:   fmt.Sprintf("%s-app", serviceName),
	Short: fmt.Sprintf("%s application", serviceName),
	Long:  fmt.Sprintf("%s application of bookstore microservices", serviceName),
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: fmt.Sprintf("run %s application", serviceName),
	Long:  fmt.Sprintf("run %s application of bookstore microservices", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("%s\n", metadata())
		env, _ := cmd.Flags().GetString("env")
		server := di.InitializeServer(object.Env(env), meta)
		server.Run()
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: fmt.Sprintf("migrate %s database schema", serviceName),
	Long:  fmt.Sprintf("migrate %s database schema of bookstore microservices", serviceName),
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("%s\n", metadata())
		env, _ := cmd.Flags().GetString("env")
		manager := di.InitializeMigrationManager(object.Env(env))
		manager.Migrate()
	},
}
