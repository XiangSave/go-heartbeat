package heartbeatconf

type HeartbeatSettingS struct {
	LogPath              string                 `yaml:"logPath"`
	MasterConnectSetting MasterConnectSettingS  `yaml:"masterConnectSetting"`
	SlaveConnectSetting  []SlaveConnectSettingS `yaml:"slaveConnectSetting"`
}

type SlaveConnectSettingS struct {
	Name               string                 `yaml:"name"`
	DBType             string                 `yaml:"dbType"`
	Host               string                 `yaml:"host"`
	UserName           string                 `yaml:"userName"`
	Password           string                 `yaml:"password"`
	Charset            string                 `yaml:"charset"`
	DbName             string                 `yaml:"dbName"`
	TblName            string                 `yaml:"tblName"`
	MaxIdleConnections int                    `yaml:"maxIdleConnections"`
	MaxOpenConnections int                    `yaml:"maxOpenConnections"`
	MonitorRoler       []MonitorRolerSettingS `yaml:"monitorRoler"`
}

type MonitorRolerSettingS struct {
	During       int `yaml:"during"`
	LaterSeconds int `yaml:"laterSeconds"`
}

type MasterConnectSettingS struct {
	Name               string `yaml:"name"`
	DBType             string `yaml:"dbType"`
	Host               string `yaml:"host"`
	UserName           string `yaml:"userName"`
	Password           string `yaml:"password"`
	Charset            string `yaml:"charset"`
	DbName             string `yaml:"dbName"`
	TblName            string `yaml:"tblName"`
	MaxIdleConnections int    `yaml:"maxIdleConnections"`
	MaxOpenConnections int    `yaml:"maxOpenConnections"`
}

func (s *Setting) ReadHeartbeatSetting(k string, v interface{}) error {
	err := s.vp.UnmarshalKey(k, v)
	if err != nil {
		return err
	}
	return nil
}
