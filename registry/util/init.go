package util

import (
	"database/sql"
	"ela/foundation/constants"
	"ela/foundation/errors"
	"ela/foundation/path"
	"ela/internal/cwd/system/global"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var Db *sql.DB

func initialize() error {
	if Db != nil {
		return nil
	}
	srcDir := path.GetSystemAppData(constants.SYSTEM_SERVICE_ID)
	src := srcDir + "/" + global.DB_NAME
	log.Println("Registry:", "DB Initialize @", src)
	if err := os.MkdirAll(srcDir, 0700); err != nil {
		return err
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
		source varchar(400) not null
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
