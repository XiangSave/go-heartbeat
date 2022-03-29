package masterupdate

import (
	"go-heartbeat/global"
	"go-heartbeat/pkg/mysql"
	"time"

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

func masterupdate(con *mysql.DBInfo, tblName string) error {
	// 查看 server id(仅查看) show variables like 'server_id';

	// 查看 master status(仅查看) show master status

	// con.RunExec("UPDATE %s SET ")

}
