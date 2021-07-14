package records

import "database/sql"

func executeQuery(query string, args ...interface{}) error {
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(args...)
	if err != nil {
		return err
	}
	return nil
}

func selectQuery(query string, args ...interface{}) (*sql.Rows, error) {
	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	row, err := stmt.Query(args...)
	if err != nil {
		return nil, err
	}
	return row, nil
}
