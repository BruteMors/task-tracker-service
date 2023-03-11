package http_api

import (
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"task-tracker-service/internal/config"
	"task-tracker-service/internal/database"
	"task-tracker-service/internal/domain/service"
	"task-tracker-service/internal/store"
)

type ApiServer struct {
	router      *mux.Router
	logger      *logrus.Logger
	config      *config.Config
	storageType string
}

func NewApiServer(config *config.Config, storageType string) *ApiServer {
	return &ApiServer{
		router:      mux.NewRouter(),
		logger:      logrus.New(),
		config:      config,
		storageType: storageType,
	}
}

func (s *ApiServer) Start() error {
	err := s.ConfigureLogger()
	if err != nil {
		return err
	}
	s.logger.Info("Logger initialization from the configuration file")
	s.logger.Info("Init database")
	var db service.DayToDoStorage
	if s.storageType == "local" {
		s.logger.Info("Starting local storage")
		db = database.NewDayToDoStorage()
	}
	if s.storageType == "database" {
		s.logger.Info("Starting database storage")
		s.logger.Info("Init client database")
		client := store.NewStore(s.config, s.logger)
		err := client.Open()
		if err != nil {
			return err
		}
		defer client.Close()
		s.logger.Info("Init dayToDoRepository")
		db = store.NewDayToDoRepository(client)
	}
	s.logger.Info("Init service")
	srvc := service.NewDayToDoService(db)
	s.logger.Info("Init handler")
	hanlder := NewHandler(srvc, s.config, s.logger, s.router)
	err = hanlder.InitRouts()
	if err != nil {
		return err
	}
	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *ApiServer) ConfigureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}
	s.logger.SetLevel(level)
	return nil
}
