package slaveselect

import (
	"go-heartbeat/internal/cronjobs/query"
	"go-heartbeat/internal/heartbeatconf"
	"go-heartbeat/pkg/mysql"

	log "github.com/sirupsen/logrus"

	"github.com/pkg/errors"
)

type SlaveConnectionS struct {
	SlaveSetting heartbeatconf.SlaveConnectSettingS
	Con          *mysql.DBModel
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
	log.Infof("%s: %d", s.SlaveSetting.Name, timestamp)
}
