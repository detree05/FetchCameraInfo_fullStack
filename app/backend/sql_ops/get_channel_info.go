package sql_ops

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type ChannelInfo struct {
	Name        string
	ExtId       string
	ChannelId   string
	ControlHost string
}

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

func queryChannelInfo(db *sql.DB, cameraName string) (map[string][]ChannelInfo, int, error) {
	// what the fuck it returns:
	// {
	//   "10.0.0.1": [
	//     {
	//        "name": "",
	//        "channel_id": "",
	//        "ext_id": "",
	//        "control_host": ""
	//     },
	//     {
	//        "name": "",
	//        "channel_id": "",
	//        "ext_id": "",
	//        "control_host": ""
	//     },
	//   ],
	// } -- you get the idea
	var (
		channel       ChannelInfo
		channelsCount int
	)
	channels := make(map[string][]ChannelInfo)

	rows, err := db.Query(`select cam.name, cam.channel_id ext_id, chan.channel_id, ist3.control_host from cctv.camera cam
						   join cctv_video.channel chan on cam.channel_id = chan.ext_id
						   join cctv_video.istream3 ist3 on chan.server_id = ist3.id
						   where cam.name like ?
						   order by ist3.control_host`, fmt.Sprintf("%%%s%%", cameraName))
	if err != nil {
		return nil, 0, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&channel.Name, &channel.ExtId, &channel.ChannelId, &channel.ControlHost)
		if err != nil {
			return nil, 0, err
		}

		channelsCount += 1
		channels[channel.ControlHost] = append(channels[channel.ControlHost], channel)
	}

	err = rows.Err()
	if err != nil {
		return nil, 0, err
	}

	return channels, channelsCount, err
}

func GetChannelInfo(projectName, cameraName string) (map[string][]ChannelInfo, int, error) {
	username, password, address, err := getProjectInfo(projectName)
	if err != nil {
		return nil, 0, err
	}

	var connString string = fmt.Sprintf("%s:%s@tcp(%s)/", username, password, address) // "user:password@tcp(127.0.0.1:3306)/hello" where 'hello' is database name

	db, err := initSQLConnection(connString)
	if err != nil {
		return nil, 0, err
	}

	defer db.Close()

	channels, channelsCount, err := queryChannelInfo(db, cameraName)
	if err != nil {
		return nil, 0, err
	}

	return channels, channelsCount, err
}
