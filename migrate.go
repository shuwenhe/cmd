package cmd

import (
	"fmt"
	"github.com/sidsa-service/sidsa-service/interval/pkg/model"
	"github.com/sidsa-service/sidsa-service/pkg/database"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	migrate = &cobra.Command{
		Use:   "migrate",
		Short: "迁移数据，创建 table",
		RunE: func(cmd *cobra.Command, args []string) (err error) {

			err = database.InitMysql(
				logger,
				viper.GetString("db.host"),
				viper.GetInt("db.port"),
				viper.GetString("db.username"),
				viper.GetString("db.password"),
				viper.GetString("db.dbname"),
				viper.GetBool("db.logMode"),
			)
			//err = database.InitClickhouse(viper.GetString("db.clickhouse.addr"), viper.GetString("db.clickhouse.dbname"), viper.GetString("db.clickhouse.user"), viper.GetString("db.clickhouse.password"))
			if err != nil {
				return fmt.Errorf("mysql conn fail:%w", err)
			}

			return database.Mysql.AutoMigrate(
				model.Photo{},
			)
		}}
)
