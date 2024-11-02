package sql_ops

import "database/sql"

func initSQLConnection(conn_string string) (*sql.DB, error) {
	db, err := sql.Open("mysql", conn_string) // Creating sql.DB object(!) via which we will initiate connections to database
	if err != nil {
		return nil, err
	}

	err = db.Ping() // Checking if connection to database is even possible (it can be inaccessible via ethernet or credentials are wrong)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, err
}
