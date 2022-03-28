package heartbeatconf

type HeartbeatSettings struct {
	MasterConnectSetting ConnectSettingS   `yaml:"masterConnectSetting"`
	SlaveConnectSetting  []ConnectSettingS `yaml:"slaveConnectSettings"`
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
