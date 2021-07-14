package records

import (
	"ela/foundation/app/data"
	"ela/foundation/errors"
)

func retrievePackagesFor(action string, table string) ([]string, error) {
	query := `select packageId from ` + table + ` where action = ?`
	row, err := selectQuery(query, action)
	if err != nil {
		return nil, errors.SystemNew("records.RetrieveAction failed to retrieve packages for "+action, err)
	}
	packages := make([]string, 0, 5)
	length := 0
	defer row.Close()
	for row.Next() {
		packages = append(packages, "")
		row.Scan(&packages[length])
		length++
	}
	return packages, nil
}

func RetrievePackagesForActivity(action string) ([]string, error) {
	return retrievePackagesFor(action, "activities")
}

func RemoveActivities(pkId string) error {
	query := "delete from activities where packageId = ?"
	//args = append(args, pkg.PackageId, activity)
	if err := executeQuery(query, pkId); err != nil {
		return errors.SystemNew("records.RemoveActivities failed for "+pkId, err)
	}
	return nil
}

func AddActivities(pkg *data.PackageConfig) error {
	if err := RemoveActivities(pkg.PackageId); err != nil {
		return err
	}
	query := ""
	//args := make([]interface{}, 0, 6)
	for _, activity := range pkg.Activities {
		query = query + "insert into activities(packageId, action) values(?,?);"
		//args = append(args, pkg.PackageId, activity)
		if err := executeQuery(query, pkg.PackageId, activity); err != nil {
			return errors.SystemNew("records.AddActivities failed to add activities for "+pkg.PackageId, err)
		}
	}
	return nil
}

/*
func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return retrieveAction(action, "broadcast_actions")
}*/
