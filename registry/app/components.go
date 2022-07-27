// Copyright 2021 The Elabox Authors
// This file is part of the elabox-system-tools library.

// The elabox-system-tools library is under open source LGPL license.
// If you simply compile or link an LGPL-licensed library with your own code,
// you can release your application under any license you want, even a proprietary license.
// But if you modify the library or copy parts of it into your code,
// you’ll have to release your application under similar terms as the LGPL.
// Please check license description @ https://www.gnu.org/licenses/lgpl-3.0.txt

package app

import (
	"github.com/cansulting/elabox-system-tools/foundation/app/data"
	"github.com/cansulting/elabox-system-tools/foundation/errors"
	"github.com/cansulting/elabox-system-tools/foundation/logger"
	"github.com/cansulting/elabox-system-tools/registry/util"
)

func retrievePackagesFor(action string, table string) ([]string, error) {
	query := `select packageId from ` + table + ` where action = ?`
	row, err := util.SelectQuery(query, action)
	if err != nil {
		logger.GetInstance().Error().Err(err).Caller().Msg("Failed to retrieve packages for " + action)
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

// retrieve all packages for specific activity
// @return: []string the list of package ids
func RetrievePackagesForActivity(action string) ([]string, error) {
	return retrievePackagesFor(action, "activities")
}

func RetrieveActivities(pkid string) ([]string, error) {
	query := `select action from activities where packageId = ?`
	row, err := util.SelectQuery(query, pkid)
	if err != nil {
		logger.GetInstance().Error().Err(err).Caller().Msg("Failed to retrieve activities for " + pkid)
		return nil, errors.SystemNew("records.RetrieveAction failed to retrieve packages for "+pkid, err)
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

func removeActivities(pkId string) error {
	query := "delete from activities where packageId = ?"
	//args = append(args, pkg.PackageId, activity)
	if err := util.ExecuteQuery(query, pkId); err != nil {
		return errors.SystemNew("records.RemoveActivities failed for "+pkId, err)
	}
	return nil
}

// register activities for specific package
func registerActivities(pkg *data.PackageConfig) error {
	if err := removeActivities(pkg.PackageId); err != nil {
		return err
	}
	query := ""
	//args := make([]interface{}, 0, 6)
	for _, activity := range pkg.ActivityGroup.Activities {
		query = query + "insert into activities(packageId, action) values(?,?);"
		//args = append(args, pkg.PackageId, activity)
		if err := util.ExecuteQuery(query, pkg.PackageId, activity); err != nil {
			return errors.SystemNew("records.AddActivities failed to add activities for "+pkg.PackageId, err)
		}
	}
	return nil
}
func setEnableService(pk string, status bool) error {
	query := "insert or replace into service_status(packageId,status) values(?,?);"
	if err := util.ExecuteQuery(query, pk, status); err != nil {
		return errors.SystemNew("error in updating status of "+pk, err)
	}
	return nil
}

/*
func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return retrieveAction(action, "broadcast_actions")
}*/
