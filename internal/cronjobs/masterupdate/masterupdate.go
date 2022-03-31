package masterupdate

import (
	"fmt"
	"go-heartbeat/global"
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/pkg/mysql"
	"time"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func MasterUpdate() error {
	time.Sleep(3 * time.Second)
	times := time.Now().UnixNano()
	log.Printf("master update: %v", times)

	MasterInfo := &mysql.DBInfo{
		DBType:             global.HeartbeatSetting.MasterConnectSetting.DBType,
		Host:               global.HeartbeatSetting.MasterConnectSetting.Host,
		UserName:           global.HeartbeatSetting.MasterConnectSetting.Name,
		Password:           global.HeartbeatSetting.MasterConnectSetting.Password,
		Charset:            global.HeartbeatSetting.MasterConnectSetting.Charset,
		DatabaseName:       global.HeartbeatSetting.MasterConnectSetting.DbName,
		MaxIdleConnections: global.HeartbeatSetting.MasterConnectSetting.MaxIdleConnections,
		MaxOpenConnections: global.HeartbeatSetting.MasterConnectSetting.MaxOpenConnections,
	}

	connect := mysql.NewDBModel(MasterInfo)
	err := connect.Connect()
	if err != nil {
		return err
	}

	return nil
}

func masterupdate(con *mysql.DBModel, tblName string) error {
	serverId, err := query.GetServerId(con)
	if err != nil {
		return err
	}
	binlogFile, position, err := query.GetPosition(con)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET `ts`=\"%d\" ,`file`=\"%s\",`position`=\"%s\"  WHERE `server_id`= %d", tblName, time.Now().UnixNano(), binlogFile, position, serverId)

	affenctedRows, err := con.RunExec(query)
	if err != nil {
		return errors.Wrapf(err, "query sql error %s", query)
	}
	if affenctedRows != 1 {
		return errors.Errorf("update timestamp affenced rows error, affend rows: %d", affenctedRows)
	}

	return nil
}
