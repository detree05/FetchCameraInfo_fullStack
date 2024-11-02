package sql_ops

import (
	"fmt"

	"fci-backend.detree05.com/cfg"
	_ "github.com/go-sql-driver/mysql"
)

func getProjectInfo(projectName string) (string, string, string, error) {
	var (
		connString  string
		sqlUsername string
		sqlPassword string
		sqlAddress  string
	)

	connString = fmt.Sprintf("%s:%s@tcp(%s:%s)/projects",
		cfg.Config.Database.Username,
		cfg.Config.Database.Password,
		cfg.Config.Database.Address,
		cfg.Config.Database.Port,
	)

	db, err := initSQLConnection(connString)
	if err != nil {
		return "", "", "", err
	}

	defer db.Close()

	rows, err := db.Query(`select sql_address, sql_username, sql_password from projects_info
						   where project_name like ?`, fmt.Sprintf("%%%s%%", projectName))
	if err != nil {
		return "", "", "", err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&sqlAddress, &sqlUsername, &sqlPassword)
		if err != nil {
			return "", "", "", err
		}
	}

	if rows.Err() != nil {
		return "", "", "", err
	}

	return sqlUsername, sqlPassword, sqlAddress, err
}
