package masterupdate

import (
	"fmt"
	"go-heartbeat/global"
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/pkg/mysql"
	"time"

	"github.com/pkg/errors"
)

// 根据配置文件创建 mysql 连接对象
func MasterNewConnect() (*mysql.DBModel, error) {
	MasterInfo := &mysql.DBInfo{
		DBType:             global.HeartbeatSetting.MasterConnectSetting.DBType,
		Host:               global.HeartbeatSetting.MasterConnectSetting.Host,
		UserName:           global.HeartbeatSetting.MasterConnectSetting.UserName,
		Password:           global.HeartbeatSetting.MasterConnectSetting.Password,
		Charset:            global.HeartbeatSetting.MasterConnectSetting.Charset,
		DatabaseName:       global.HeartbeatSetting.MasterConnectSetting.DbName,
		MaxIdleConnections: global.HeartbeatSetting.MasterConnectSetting.MaxIdleConnections,
		MaxOpenConnections: global.HeartbeatSetting.MasterConnectSetting.MaxOpenConnections,
	}
	con := mysql.NewDBModel(MasterInfo)
	err := con.Connect()
	if err != nil {
		return nil, errors.Wrapf(err, "connect master db: %s failed",
			global.HeartbeatSetting.MasterConnectSetting.Name)
	}
	return con, nil
}

// 获取 server_id，binlog 文件、position 位置，并更新 master 对应行
func masterupdate(con *mysql.DBModel) error {
	binlogFile, position, err := query.GetPosition(con)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET `ts`=\"%d\" ,`file`=\"%s\",`position`=\"%s\"  WHERE `server_id`= %d",
		global.HeartbeatSetting.MasterConnectSetting.TblName,
		time.Now().UnixNano(), binlogFile, position, global.MasterServerId)

	affenctedRows, err := con.RunExec(query)
	if err != nil {
		return errors.Wrapf(err, "query sql error %s", query)
	}
	if affenctedRows != 1 {
		if affenctedRows == 0 {
			return masterinsert(con, global.MasterServerId)
		}
		return errors.Wrapf(
			errors.Errorf("update timestamp affenced rows error, affend rows: %d", affenctedRows),
			"query sql error %s", query)
	}

	return nil
}

// 如果 update 更改行为 0，则 insert 写入新数据
func masterinsert(con *mysql.DBModel, serverId int) error {
	query := fmt.Sprintf("INSERT INTO `%s`(`ts`, `server_id`, `file`, `position`, `relay_master_log_file`, `exec_master_log_pos`) VALUES (\"%d\",%d,'', 0, '', 0);",
		global.HeartbeatSetting.MasterConnectSetting.TblName,
		time.Now().UnixNano(), serverId)

	affenctedRows, err := con.RunExec(query)
	if err != nil {
		return errors.Wrapf(err, "query sql error %s", query)
	}
	if affenctedRows != 1 {
		return errors.Wrapf(
			errors.Errorf("insert timestamp affenced rows error, affend rows: %d", affenctedRows),
			"query sql error %s", query)
	}
	return nil
}
