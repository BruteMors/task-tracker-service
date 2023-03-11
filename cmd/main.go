package main

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/sirupsen/logrus"
	"log"
	"task-tracker-service/internal/config"
	"task-tracker-service/internal/controller/http_api"
	"task-tracker-service/internal/controller/stdio"
	"task-tracker-service/internal/database"
	"task-tracker-service/internal/domain/service"
	"task-tracker-service/internal/store"
)

var inputMethod string
var configPath string
var storageType string

func init() {
	flag.StringVar(&inputMethod, "input-method", "web", "choose an interface for entering a task (command line (cmd) or web api (web)")
	flag.StringVar(&configPath, "config-path", "internal/config/apiserver.toml", "path to config file (.toml or .env file)")
	flag.StringVar(&storageType, "storage-type", "database", "choose database or local storage type (database or local)")
}

func main() {
	flag.Parse()
	switch inputMethod {
	case "cmd":
		fmt.Println("start cmd app...")
		var str string
		_, err := fmt.Scan(&str)
		if err != nil {
			return
		}
		if str == "StartApp" {
			fmt.Println("App is started")
			cfg := config.NewConfig()
			_, err := toml.DecodeFile(configPath, &cfg)
			if err != nil {
				log.Println("can not find path to config, app will use default confs:", err)
			}
			logger := logrus.New()
			logger.Info("Init database")
			var storage service.DayToDoStorage
			if storageType == "local" {
				logger.Info("Starting local storage")
				storage = database.NewDayToDoStorage()
			}
			if storageType == "database" {
				logger.Info("Starting database storage")
				logger.Info("Init client database")
				client := store.NewStore(cfg, logger)
				err := client.Open()
				if err != nil {
					return
				}
				defer client.Close()
				logger.Info("Init dayToDoRepository")
				storage = store.NewDayToDoRepository(client)
			}

			srvc := service.NewDayToDoService(storage)
			handler := stdio.NewDayToDoHandler(srvc)

			for {
				cmd, err := handler.GetCommand()
				if err != nil {
					fmt.Println(err)
					continue
				}
				isExit, err := handler.ParseCommand(cmd)
				if err != nil {
					fmt.Println(err)
				}
				if isExit {
					fmt.Println("App is terminated")
					return
				}
			}
		}
	case "web":
		fmt.Println("start web app...")
		cfg := config.NewConfig()
		_, err := toml.DecodeFile(configPath, &cfg)
		if err != nil {
			log.Println("can not find path to config, app will use default confs:", err)
		}
		api := http_api.NewApiServer(cfg, storageType)

		if err := api.Start(); err != nil {
			log.Fatal(err)
		}
	}
}
