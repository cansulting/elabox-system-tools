package util

import (
	"database/sql"
	"os"

	"github.com/cansulting/elabox-system-tools/registry/config"
)

func ExecuteQuery(query string, args ...interface{}) error {
	if err := initialize(); err != nil {
		return err
	}
	stmt, err := Db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func SelectQuery(query string, args ...interface{}) (*sql.Rows, error) {
	if err := initialize(); err != nil {
		return nil, err
	}
	stmt, err := Db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return row, nil
}

func DeleteDB() error {
	if Db != nil {
		if err := Db.Close(); err != nil {
			return err
		}
		Db = nil
	}
	path := config.DB_DIR + "/" + config.DB_NAME
	if _, err := os.Stat(path); err != nil {
		return err
	}
	return os.Remove(path)
}
