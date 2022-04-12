package masterupdate

import (
	"go-heartbeat/global"
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/internal/serverinit"
	"go-heartbeat/pkg/cronjob"
	"go-heartbeat/pkg/mysql"

	log "github.com/sirupsen/logrus"
)

type MasterConnectionS struct {
	Con *mysql.DBModel
}

func MasterCronRun(jobs *cronjob.CronjobServer, con *mysql.DBModel) error {
	var err error

	// 查询 master server id 并赋给全局变量
	global.MasterServerId, err = query.GetServerId(con)
	if err != nil {
		return err
	}

	// 监控表不存在则创建监控表
	err = serverinit.MasterCreateTable(con)
	if err != nil {
		return err
	}

	// 手动触发一次 master update，第一次 slave 查询时 master 未启动导致的时间不准确
	masterupdate(con)
	_, err = jobs.AddJob("* * * * * *", MasterConnectionS{Con: con})
	if err != nil {
		return err
	}

	return nil
}

// 创建无传参函数，供 cron 调用
func (c MasterConnectionS) Run() {
	err := masterupdate(c.Con)
	if err != nil {
		log.Errorf("%+v", err)
	}
}
