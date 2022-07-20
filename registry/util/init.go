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

	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/foundation/perm"
	"github.com/cansulting/elabox-system-tools/registry/config"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func initialize() error {
	if Db != nil {
		return nil
	}
	srcDir := config.DB_DIR
	src := srcDir + "/" + config.DB_NAME
	logger.GetInstance().Info().Str("category", "registry").Msg("DB Initialize @" + src)
	if _, err := os.Stat(srcDir); err != nil {
		if err := os.MkdirAll(srcDir, perm.PUBLIC_VIEW); err != nil {
			return err
		}
	}
	_db, err := sql.Open("sqlite3", src)
	if err != nil {
		return errors.SystemNew("Initialize DB failed. Unable to open sql.", err)
	}
	Db = _db
	if err := createPackageTable(Db); err != nil {
		return errors.SystemNew("Initialize DB failed. Unable to create package table.", err)
	}
	if err := createActionTable(Db); err != nil {
		return errors.SystemNew("Initialize DB failed. Unable to create action table.", err)
	}
	if err := createServiceStatusTable(Db); err != nil {
		return errors.SystemNew("Initialize DB failed. Unable to create action table.", err)
	}
	//if createBroadcastTable(db); err != nil {
	//	return err
	//}
	return nil
}

func createPackageTable(db *sql.DB) error {
	packageQuery := `create table if not exists packages(
		id varchar(100) not null primary key,
		name varchar(40) not null,
		desc text,
		location varchar(100) not null,
		build smallint not null,
		version varchar(10) not null,
		source varchar(400) not null,
		nodejs tinyint(1) not null,
		exportService tinyint(1) not null,
		program varchar(200)
	)`
	stmt, err := db.Prepare(packageQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

func createActionTable(db *sql.DB) error {
	packageQuery := `create table if not exists activities(
		id integer primary key autoincrement,
		packageId varchar(100) not null,
		action varchar(100) not null
	)`
	stmt, err := db.Prepare(packageQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
func createServiceStatusTable(db *sql.DB) error {
	packageQuery := `create table if not exists service_status(
		packageId varchar(100) primary key not null,
		status integer not null
	)`
	stmt, err := db.Prepare(packageQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}

/*
func createBroadcastTable(db *sql.DB) error {
	packageQuery := `create table if not exists broadcast_actions(
		id int autoincrement primar key,
		packageId varchar(100) not null,
		action varchar(100) not null
	)`
	stmt, err := db.Prepare(packageQuery)
	if err != nil {
		return err
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	return nil
}
*/
