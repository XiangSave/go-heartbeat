package heartbeatconf

type HeartbeatSettingS struct {
	LogPath              string            `yaml:"logPath"`
	MasterConnectSetting ConnectSettingS   `yaml:"masterConnectSetting"`
	SlaveConnectSetting  []ConnectSettingS `yaml:"slaveConnectSetting"`
}

type ConnectSettingS struct {
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
