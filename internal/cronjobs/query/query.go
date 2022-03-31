package query

import (
	"database/sql"
	"go-heartbeat/pkg/mysql"

	"github.com/pkg/errors"

	_ "github.com/go-sql-driver/mysql"
)

type QueryServerId struct {
	VariableName string
	Value        int
}

type QueryMasterStatus struct {
	BinlogFile string
	Position   string
}

func GetServerId(con *mysql.DBModel) (int, error) {
	resServerId, err := getServerId(con)
	if err != nil {
		return 0, err
	}
	return resServerId.Value, nil
}

func GetPosition(con *mysql.DBModel) (string, string, error) {
	resMasterStatus, err := getMasterStatus(con)
	if err != nil {
		return "", "", err
	}
	return resMasterStatus.BinlogFile, resMasterStatus.Position, nil
}

func getServerId(con *mysql.DBModel) (*QueryServerId, error) {
	var serverid QueryServerId
	query := "SHOW VARIABLES LIKE \"server_id\";"
	row := con.DBEngine.QueryRow(query)
	if err := row.Scan(&serverid.VariableName, &serverid.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "query sql error %s", query)
	}
	return &serverid, nil
}

func getMasterStatus(con *mysql.DBModel) (*QueryMasterStatus, error) {
	var masterStatus QueryMasterStatus
	query := "SHOW MASTER STATUS;"
	row := con.DBEngine.QueryRow(query)
	if err := row.Scan(&masterStatus.BinlogFile, &masterStatus.Position); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrapf(err, "query sql error %s", query)
	}
	return &masterStatus, nil

}
