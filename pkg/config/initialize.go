package config

import (
	"database/sql"
	"errors"
	"log/slog"
)

func Initialise() (*Configuration, *slog.Logger, *sql.DB, error) {
	conf := loadConf()
	logger := setupLog(conf.Log)
	db, err := setupDb(conf.Db)

	return conf, logger, db, errors.Join(err, conf.Nats.validate(), conf.App.validate())
}
