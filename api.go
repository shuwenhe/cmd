package cmd

import (
	"fmt"
	"github.com/sidsa-service/sidsa-service/interval/app/apis/router"
	"github.com/sidsa-service/sidsa-service/interval/pkg/model"
	"github.com/sidsa-service/sidsa-service/pkg/database"
	"github.com/sidsa-service/sidsa-service/pkg/iot/ws"
	"github.com/sidsa-service/sidsa-service/pkg/mqtt"
	"github.com/sidsa-service/sidsa-service/pkg/netease"
	"github.com/sidsa-service/sidsa-service/pkg/sms"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"net/http"
)

var (
	api = &cobra.Command{
		Use:   "api",
		Short: "api 服务",
		RunE: func(cmd *cobra.Command, args []string) error {
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
				return err
			}

			err = database.InitClickhouse(viper.GetString("db.clickhouse.addr"), viper.GetString("db.clickhouse.dbname"), viper.GetString("db.clickhouse.user"), viper.GetString("db.clickhouse.password"))
			if err != nil {
				return fmt.Errorf("clickhouse conn fail:%w", err)
			}

			model.Migration(logger, database.Mysql)
			mqtt.Init("sidsa-service-api", viper.GetString("mqtt.addr")+":"+viper.GetString("mqtt.port"))
			sms.MustInit(viper.GetString("aliyun.sms.ak"), viper.GetString("aliyun.sms.sk"))
			ws.Init()

			imClient := netease.NewClient(viper.GetString("netease.appKey"), viper.GetString("netease.appSecret"), viper.GetString("netease.httpProxy"))
			handler, err := router.New(logger, imClient)
			if err != nil {
				return err
			}
			return http.ListenAndServe(viper.GetString("addr"), handler)
		},
	}
)
