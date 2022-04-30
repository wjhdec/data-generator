package cmd

import (
	"context"
	"datagen/internal/pg"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "data-gen",
	Short: "create data in postgresql",
	Run: func(cmd *cobra.Command, args []string) {
		connCfg, err := pgxpool.ParseConfig(viper.GetString("database.dsn"))
		if err != nil {
			log.Panic(err)
		}
		connCfg.MaxConns = viper.GetInt32("database.max_conns")
		pool, err := pgxpool.ConnectConfig(context.Background(), connCfg)
		if err != nil {
			log.Panic(err)
		}
		defer pool.Close()
		gen := pg.NewSqlGen(pool, viper.GetViper())
		gen.CreateSql()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is ./config.toml)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigName(cfgFile)
	} else {
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
	}
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Panic(err)
	}
}
