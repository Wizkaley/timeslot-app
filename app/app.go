package app

import (
	"fmt"
	"os"
	"timeslot-app/db"
	"timeslot-app/models"
	"timeslot-app/service"

	"github.com/jackc/pgx"
	"github.com/spf13/viper"
)

type App struct {
	Config          models.Config
	DB              *pgx.Conn
	TimeslotService *service.TimeslotServiceImplementaion
	UserService     *service.UserService
	EventService    *service.EventService
}

var Service *App

func InitApp() (*App, error) {
	// viper.AddConfigPath("/Users/eshankaley/go/src/timeslot-app/")
	// viper.SetConfigName("config")

	// viper.AutomaticEnv()

	// err := viper.ReadInConfig()
	// if err != nil {
	cfg := models.Config{}
	// fmt.Println("Error reading config file, %s", err)
	cfg.DBConfig.Host = os.Getenv("host")
	cfg.DBConfig.User = os.Getenv("user")
	cfg.DBConfig.Password = os.Getenv("password")
	// return nil, err
	// }

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Config loaded successfully", cfg)
	// logger.Logger(cfg.LogConfig.LogLevel)

	database := db.Connection(cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password)

	err = db.CreateTables(database)
	if err != nil {
		panic(err)
	}
	app := new(App)
	app.Config = cfg

	app.DB = database
	app.TimeslotService = service.Init(database)
	app.UserService = service.NewUserService(database)
	app.EventService = service.NewEventService(database)
	return app, nil
}
