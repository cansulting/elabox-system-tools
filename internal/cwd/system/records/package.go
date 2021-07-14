package records

import (
	"database/sql"
	"ela/foundation/app/data"
	"ela/foundation/errors"
)

// retrieve all packages
func RetrievePackagesRaw(packageId string, columns []string) (*sql.Rows, error) {
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
		row, err = selectQuery(query, packageId)
	} else {
		row, err = selectQuery(query)
	}
	if err != nil {
		return nil, errors.SystemNew("records.RetrievePackagesRaw failed to retrieve packages ", err)
	}

	return row, nil
}

// retrieve all packages
func RetrievePackages(id string) ([]*data.PackageConfig, error) {
	row, err := RetrievePackagesRaw(id, []string{"id, source, name, location"})
	if err != nil {
		return nil, err
	}
	defer row.Close()
	results := make([]*data.PackageConfig, 0, 10)
	if row.Next() {
		pk := data.DefaultPackage()
		row.Scan(&pk.PackageId, &pk.Source, &pk.Name, &pk.InstallLocation)
		results = append(results, pk)
	}
	return results, nil
}

// add package data to db
func AddPackage(pkData *data.PackageConfig) error {
	query := `
		replace into 
		packages(id, location, build, version, name, desc, source) 
		values(?,?,?,?,?,?,?)`
	err := executeQuery(
		query,
		pkData.PackageId,
		pkData.InstallLocation,
		pkData.Build,
		pkData.Version,
		pkData.Name,
		pkData.Description,
		pkData.Source,
	)
	if err != nil {
		return errors.SystemNew("records.AddPackage Failed to add "+pkData.PackageId, err)
	}
	return nil
}
