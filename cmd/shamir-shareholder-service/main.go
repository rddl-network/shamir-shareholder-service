package main

import (
	"bytes"
	"log"
	"os"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/rddl-network/shamir-shareholder-service/config"
	"github.com/rddl-network/shamir-shareholder-service/service"
	"github.com/spf13/viper"
)

func loadConfig(path string) (cfg *config.Config, err error) {
	v := viper.New()
	v.AddConfigPath(path)
	v.SetConfigName("app")
	v.SetConfigType("toml")
	v.AutomaticEnv()

	err = v.ReadInConfig()
	if err == nil {
		cfg = config.GetConfig()
		cfg.ServiceHost = v.GetString("service-host")
		cfg.ServicePort = v.GetInt("service-port")
		cfg.DBPath = v.GetString("db-path")
		return
	}
	log.Println("no config file found")

	tmpl := template.New("appConfigFileTemplate")
	configTemplate, err := tmpl.Parse(config.DefaultConfigTemplate)
	if err != nil {
		return
	}

	var buffer bytes.Buffer
	if err = configTemplate.Execute(&buffer, config.GetConfig()); err != nil {
		return
	}

	if err = v.ReadConfig(&buffer); err != nil {
		return
	}
	if err = v.SafeWriteConfig(); err != nil {
		return
	}

	log.Println("default config file created. please adapt it and restart the application. exiting...")
	os.Exit(0)
	return
}

func main() {
	cfg, err := loadConfig("./")
	if err != nil {
		log.Fatalf("fatal error loading config file: %s", err)
	}

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	db, err := service.InitDB(cfg.DBPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	service := service.NewShamirService(router, db)
	err = service.Run()
	if err != nil {
		log.Fatalf("fatal error spinning up service: %s", err)
	}
}
