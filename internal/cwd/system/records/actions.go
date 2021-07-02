package records

func retrieveAction(action string, table string) ([]string, error) {
	query := `select packageId from ` + table + ` where action = ?`
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query(action)
	if err != nil {
		return nil, err
	}
	packages := make([]string, 0, 5)
	length := 0
	defer row.Close()
	for row.Next() {
		row.Scan(&packages[length])
		length++
	}
	return packages, nil
}

func RetrievePackagesWithActivity(action string) ([]string, error) {
	return retrieveAction(action, "activity_actions")
}

func RetrievePackagesWithBroadcast(action string) ([]string, error) {
	return retrieveAction(action, "broadcast_actions")
}
