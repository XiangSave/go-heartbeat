package masterupdate

import (
	"fmt"
	"go-heartbeat/global"
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/pkg/mysql"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

var MasterUpdateVars struct {
	Con *mysql.DBModel
}

// 根据配置文件创建 mysql 连接对象
func MasterNewConnect() error {
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
	MasterUpdateVars.Con = mysql.NewDBModel(MasterInfo)
	err := MasterUpdateVars.Con.Connect()
	if err != nil {
		return errors.Wrapf(err, "connect master db: %s failed",
			global.HeartbeatSetting.MasterConnectSetting.Name)
	}
	return nil
}

// 创建无传参函数，供 cron 调用
func MasterUpdate() {
	// con, err := MasterNewConnect()
	// if err != nil {
	// 	log.Errorf("%+v", err)
	// }
	err := masterupdate(MasterUpdateVars.Con)
	if err != nil {
		log.Errorf("%+v", err)
	}
}

// 获取 server_id，binlog 文件、position 位置，并更新 master 对应行
func masterupdate(con *mysql.DBModel) error {
	serverId, err := query.GetServerId(con)
	if err != nil {
		return err
	}
	binlogFile, position, err := query.GetPosition(con)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE %s SET `ts`=\"%d\" ,`file`=\"%s\",`position`=\"%s\"  WHERE `server_id`= %d",
		global.HeartbeatSetting.MasterConnectSetting.TblName,
		time.Now().UnixNano(), binlogFile, position, serverId)

	affenctedRows, err := con.RunExec(query)
	if err != nil {
		return errors.Wrapf(err, "query sql error %s", query)
	}
	if affenctedRows != 1 {
		if affenctedRows == 0 {
			return masterinsert(con, serverId)
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
