package serverinit

import (
	"fmt"
	"go-heartbeat/global"
	"go-heartbeat/pkg/mysql"
	"log"

	"github.com/pkg/errors"
)

// 数据库创建主从监控表
func MasterCreateTable(con *mysql.DBModel) error {
	log.Println(global.HeartbeatSetting.MasterConnectSetting.TblName)
	query := fmt.Sprintf("CREATE TABLE IF NOT EXISTS `%s` ( `ts` varchar(26) NOT NULL, `server_id` int(10) unsigned NOT NULL, `file` varchar(255) DEFAULT NULL, `position` bigint(20) unsigned DEFAULT NULL, `relay_master_log_file` varchar(255) DEFAULT NULL, `exec_master_log_pos` bigint(20) unsigned DEFAULT NULL, PRIMARY KEY (`server_id`) ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;", global.HeartbeatSetting.MasterConnectSetting.TblName)

	_, err := con.RunExec(query)
	if err != nil {
		return errors.Wrapf(err, "query sql error %s", query)
	}
	return nil
}

// 打印主从监控数据库初始化 sql
func EchoDBInitCmd() {
	fmt.Printf("CREATE DATABASE `%s` DEFAULT CHARACTER SET utf8mb4;\n",
		global.HeartbeatSetting.MasterConnectSetting.DbName)

	fmt.Printf("CREATE USER '%s' identified by '%s';\n",
		global.HeartbeatSetting.MasterConnectSetting.UserName,
		global.HeartbeatSetting.MasterConnectSetting.Password)

	fmt.Printf("GRANT SELECT,UPDATE,INSERT,CREATE ON %s.* To '%s'@'%%';\n",
		global.HeartbeatSetting.MasterConnectSetting.DbName,
		global.HeartbeatSetting.MasterConnectSetting.UserName)

	fmt.Printf("GRANT SUPER,REPLICATION CLIENT ON *.* To '%s'@'%%';\n",
		global.HeartbeatSetting.MasterConnectSetting.UserName)
}
