package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"task-tracker-service/internal/config"
)

type ClientDB struct {
	config *config.Config
	logger *logrus.Logger
	db     *sql.DB
}

func NewStore(config *config.Config, logger *logrus.Logger) *ClientDB {
	return &ClientDB{
		config: config,
		logger: logger,
	}
}

func (s *ClientDB) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}
	err = db.Ping()
	if err != nil {
		return err
	}
	s.db = db
	s.logger.Info("Connection to db successfully")
	return nil
}

func (s *ClientDB) Close() {
	s.db.Close()
}
