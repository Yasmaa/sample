package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type config struct {
	Postgres struct {
		DBMS   string
		USER   string
		PASS   string
		HOST   string
		PORT   string
		DBNAME string
		OPTION string
	}

	Smtp struct {
		USER     string
		PASSWORD string
		FROM     string
		HOST     string
		PORT	 string
	}
	Redis struct {
		HOST     string
		PORT	 string
	}

	Link struct {
		HOST     string
		PORT	 string
	}

	Token string
}

var C config

func LoadConfig() {
	Conf := &C
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath(filepath.Join("config"))

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("load config error")
		fmt.Println(err)
		os.Exit(1)
	}

	if err := viper.Unmarshal(&Conf); err != nil {
		fmt.Println("config Unmarshal error")
		fmt.Println(err)
		os.Exit(1)
	}
}
