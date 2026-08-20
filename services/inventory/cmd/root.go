package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/codejsha/shared-library-go/pkg/config"

	"github.com/codejsha/bookstore-microservices/inventory/internal/di"
)

const serviceName = "inventory"

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
	runCmd.Flags().StringP("label", "l", "", "set config label")
	runCmd.Flags().StringP("config-server", "c", "configserver:http://localhost:8888", "set cloud config import")

	viper.AutomaticEnv()
	_ = viper.BindEnv("profile", "APP_CONFIG_PROFILE")
	_ = viper.BindEnv("label", "APP_CONFIG_LABEL")
	_ = viper.BindEnv("config-server", "APP_CONFIG_SERVER")
	_ = viper.BindPFlag("profile", runCmd.Flags().Lookup("profile"))
	_ = viper.BindPFlag("label", runCmd.Flags().Lookup("label"))
	_ = viper.BindPFlag("config-server", runCmd.Flags().Lookup("config-server"))

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
		label := viper.GetString("label")
		configServer := viper.GetString("config-server")

		preConfig := &config.PreConfig{
			ServiceName:      serviceName,
			Profile:          config.Profile(profile),
			Label:            label,
			ConfigServerAddr: configServer,
		}
		di.NewApp(preConfig, meta).Run()
	},
}
