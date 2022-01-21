package cmd

import (
	"os"

	"github.com/sidsa-service/sidsa-service/pkg/config"
	"github.com/sidsa-service/sidsa-service/pkg/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	logger  = &log.Logger{}
	rootCmd = &cobra.Command{}
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "configs/dev.yaml", "config file (default is $HOME/.cobra.yaml)")
	rootCmd.PersistentFlags().Bool("debug", true, "开启debug")
	viper.SetDefault("gin.mode", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.AddCommand(api)
	rootCmd.AddCommand(migrate)
	rootCmd.AddCommand(realtimeData)
}

func Execute() error {
	return rootCmd.Execute()

}

func initConfig() {
	// 配置初始化
	config.MustInit(os.Stdout, cfgFile)
	// 日志
	logger = log.New(viper.GetString("service_name"), viper.GetString("log.file"))
}
