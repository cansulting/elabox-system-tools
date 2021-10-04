package app

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

func retrievePackagesFor(action string, table string) ([]string, error) {
	query := `select packageId from ` + table + ` where action = ?`
	row, err := util.SelectQuery(query, action)
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

func removeActivities(pkId string) error {
	query := "delete from activities where packageId = ?"
	//args = append(args, pkg.PackageId, activity)
	if err := util.ExecuteQuery(query, pkId); err != nil {
		return errors.SystemNew("records.RemoveActivities failed for "+pkId, err)
	}
	return nil
}

func addActivities(pkg *data.PackageConfig) error {
	if err := removeActivities(pkg.PackageId); err != nil {
		return err
	}
	query := ""
	//args := make([]interface{}, 0, 6)
	for _, activity := range pkg.Activities {
		query = query + "insert into activities(packageId, action) values(?,?);"
		//args = append(args, pkg.PackageId, activity)
		if err := util.ExecuteQuery(query, pkg.PackageId, activity); err != nil {
			return errors.SystemNew("records.AddActivities failed to add activities for "+pkg.PackageId, err)
		}
	}
	return nil
}

/*
func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return retrieveAction(action, "broadcast_actions")
}*/
