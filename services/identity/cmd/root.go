package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/codejsha/bookstore-microservices/commonlib-go/pkg/object"
	"github.com/codejsha/bookstore-microservices/identity/internal/di"
)

const serviceName = "identity"

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		_ = fmt.Errorf("failed to execute root command: %v", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().StringP("env", "e", string(object.Local), "set environment")

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
