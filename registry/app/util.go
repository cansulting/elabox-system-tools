package app

import (
	"database/sql"

	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

func convertRawToPackageConfig(rows *sql.Rows) []*data.PackageConfig {
	results := make([]*data.PackageConfig, 0, 10)
	for rows.Next() {
		pk := data.DefaultPackage()
		rows.Scan(&pk.PackageId, &pk.Source, &pk.Version, &pk.Name, &pk.InstallLocation, &pk.Nodejs)
		results = append(results, pk)
	}
	return results
}

// retrieve all packages
func retrievePackagesRaw(packageId string, columns []string) (*sql.Rows, error) {
	columnsStr := "*"
	if len(columns) > 0 {
		columnsStr = ""
		for index, column := range columns {
			columnsStr += column
			if index < len(columns)-1 {
				columnsStr += ","
			}
		}
	}
	query := `select ` + columnsStr + ` from packages`
	if packageId != "" {
		query += ` where id = ?`
	}

	var row *sql.Rows
	var err error
	if packageId != "" {
		row, err = util.SelectQuery(query, packageId)
	} else {
		row, err = util.SelectQuery(query)
	}
	if err != nil {
		return nil, errors.SystemNew("records.RetrievePackagesRaw failed to retrieve packages ", err)
	}

	return row, nil
}

func retrievePackagesWhere(columns []string, where string) (*sql.Rows, error) {
	columnsStr := "*"
	if len(columns) > 0 {
		columnsStr = ""
		for index, column := range columns {
			columnsStr += column
			if index < len(columns)-1 {
				columnsStr += ","
			}
		}
	}
	query := `select ` + columnsStr + ` from packages`
	if where != "" {
		query += ` where ` + where
	}

	var row *sql.Rows
	var err error
	row, err = util.SelectQuery(query)

	if err != nil {
		return nil, errors.SystemNew("records.RetrievePackagesRaw failed to retrieve packages ", err)
	}

	return row, nil
}

func CloseDB() error {
	return util.Db.Close()
}

/*
func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return records.RetrievePackagesWithBroadcast(action)
}*/

func retrievePackageSource(packageId string) (string, error) {
	pk, err := RetrievePackage(packageId)
	if err != nil {
		return "", err
	}
	if pk != nil {
		return pk.Source, nil
	}
	return "", nil
}
