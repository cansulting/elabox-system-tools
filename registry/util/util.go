// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

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

func Count(table string, where string, args ...interface{}) (int, error) {
	if err := initialize(); err != nil {
		return 0, err
	}
	query := "select count(*) from " + table
	if where != "" {
		query += " where " + where
	}
	row, err := Db.Query(query, args...)
	if err != nil {
		return 0, err
	}
	defer row.Close()
	row.Next()
	var count int
	row.Scan(&count)
	return count, nil
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
