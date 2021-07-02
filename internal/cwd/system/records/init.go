package records

import (
	"database/sql"
	"ela/foundation/constants"
	"ela/foundation/path"
	"ela/internal/cwd/system/global"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Initialize() error {
	src := path.GetSystemAppData(constants.SYSTEM_SERVICE_ID) + "/" + global.DB_NAME
	log.Println("System:Records", "DB Initialize @", src)
	_db, err := sql.Open("sqlite3", src)
	if err != nil {
		return err
	}
	db = _db
	if createPackageTable(db); err != nil {
		return err
	}
	if createActionTable(db); err != nil {
		return err
	}
	if createBroadcastTable(db); err != nil {
		return err
	}
	return nil
}

func createPackageTable(db *sql.DB) error {
	packageQuery := `create table if not exists packages(
		id varchar(100) not null primary key,
		name varchar(100) not null,
		desc text,
		location varchar(20) not null,
		build smallint not null,
		version varchar(10) not null,
		has_service tinyint default 0,
		has_activity tinyint default 0,
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
	packageQuery := `create table if not exists activity_actions(
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
