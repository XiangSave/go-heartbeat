package slaveselect

import (
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/internal/heartbeatconf"
	"go-heartbeat/pkg/mysql"
	"go-heartbeat/pkg/rolling"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type SlaveConnectionS struct {
	SlaveSetting  heartbeatconf.SlaveConnectSettingS
	Con           *mysql.DBModel
	RollingTiming *rolling.Timing
}

func SlaveNewConnect(c heartbeatconf.SlaveConnectSettingS) (*mysql.DBModel, error) {
	dbInfo := &mysql.DBInfo{
		DBType:             c.DBType,
		Host:               c.Host,
		UserName:           c.UserName,
		Password:           c.Password,
		Charset:            c.Charset,
		DatabaseName:       c.DbName,
		MaxIdleConnections: c.MaxIdleConnections,
		MaxOpenConnections: c.MaxOpenConnections,
	}
	con := mysql.NewDBModel(dbInfo)
	err := con.Connect()
	if err != nil {
		return nil, errors.Wrapf(err, "connect slave db: %s failed", c.Name)
	}

	return con, nil
}

func (s SlaveConnectionS) Run() {

	timestamp, err := query.GetTimestamp(s.Con, s.SlaveSetting.TblName)
	if err != nil {
		log.Error(err)
	}
	now := time.Now()
	log.Info(now.Sub(time.Unix(0, timestamp)))

	s.RollingTiming.Add(now.Sub(time.Unix(0, timestamp)))
	log.Println(s.SlaveSetting.Name, s.RollingTiming.Average())
}
