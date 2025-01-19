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

	cfg := models.Config{}
	cfg.DBConfig.Host = os.Getenv("host")
	cfg.DBConfig.User = os.Getenv("user")
	cfg.DBConfig.Password = os.Getenv("password")

	err := viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	fmt.Println("Config loaded successfully", cfg)

	database := db.Connection(cfg.DBConfig.Host, cfg.DBConfig.User, cfg.DBConfig.Password)

	err = db.CreateTables(database)
	if err != nil {
		panic(err)
	}
	app := new(App)
	app.Config = cfg

	app.DB = database
	app.TimeslotService = service.NewTimeslotService(database)
	app.UserService = service.NewUserService(database)
	app.EventService = service.NewEventService(database)
	return app, nil
}
