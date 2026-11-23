package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/catalog/internal/di"
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
	runCmd.Flags().StringP("profile", "p", string(config.ProfileLocal), "set config profile")

	viper.AutomaticEnv()
	_ = viper.BindEnv("profile", "APP_CONFIG_PROFILE")
	_ = viper.BindPFlag("profile", runCmd.Flags().Lookup("profile"))

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
		profile := viper.GetString("profile")

		preConfig := &config.PreConfig{
			ServiceName: serviceName,
			Profile:     config.Profile(profile),
		}
		di.NewApp(preConfig, meta).Run()
	},
}
