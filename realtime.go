package cmd

import (
	"fmt"
	"github.com/sidsa-service/sidsa-service/interval/app/realtime"
	"github.com/sidsa-service/sidsa-service/pkg/database"
	"github.com/sidsa-service/sidsa-service/pkg/iot/ws"
	"github.com/sidsa-service/sidsa-service/pkg/mqtt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	realtimeData = &cobra.Command{
		Use:   "realtime",
		Short: "实时数据",
		RunE: func(cmd *cobra.Command, args []string) error {
			mqtt.Init("sidsa-service-realtime", viper.GetString("mqtt.addr")+":"+viper.GetString("mqtt.port"))

			err := database.InitMysql(
				logger,
				viper.GetString("db.host"),
				viper.GetInt("db.port"),
				viper.GetString("db.username"),
				viper.GetString("db.password"),
				viper.GetString("db.dbname"),
				viper.GetBool("db.logMode"),
			)
			if err != nil {
				return fmt.Errorf("mysql conn fail:%w", err)
			}
			err = database.InitClickhouse(viper.GetString("db.clickhouse.addr"), viper.GetString("db.clickhouse.dbname"), viper.GetString("db.clickhouse.user"), viper.GetString("db.clickhouse.password"))
			if err != nil {
				return fmt.Errorf("clickhouse conn fail:%w", err)
			}

			ws.Init()
			realtime.InitWsPool()
			realtime.InitTransfer(mqtt.Client)
			err = realtime.RunDataStoreJob()
			if err != nil {
				return fmt.Errorf("data store job run fail:%w", err)
			}

			return realtime.Run(":8991")
		}}
)
