package records

import (
	"database/sql"
	"ela/foundation/app/data"
)

func RetrievePackageRows(packageId string, columns []string) (*sql.Rows, error) {
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
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}

	var row *sql.Rows
	if packageId != "" {
		row, err = stmt.Query(packageId)
	} else {
		row, err = stmt.Query()
	}
	if err != nil {
		return nil, err
	}

	return row, nil
}

func RetrievePackage(id string) ([]*data.PackageConfig, error) {
	row, err := RetrievePackageRows(id, []string{"id, source, name"})
	if err != nil {
		return nil, err
	}
	defer row.Close()
	results := make([]*data.PackageConfig, 0, 10)
	if row.Next() {
		pk := data.DefaultPackage()
		row.Scan(&pk.PackageId, &pk.Source, &pk.InstallLocation)
		results = append(results, pk)
	}
	return results, nil
}

func AddPackage(pkData *data.PackageConfig) error {
	query := `
		replace into 
		packages(id, location, build, version, has_service, has_activity, name, desc, source) 
		values(?,?,?,?,?,?,?,?,?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	hasService := 0
	hasActivity := 0
	if len(pkData.Services) > 0 {
		hasService = 1
	}
	res, err := stmt.Exec(
		pkData.PackageId,
		pkData.InstallLocation,
		pkData.Build,
		pkData.Version,
		hasService,
		hasActivity,
		pkData.Name,
		pkData.Description,
		pkData.Source,
	)
	if err != nil {
		return err
	}
	println(res)
	return nil
}
