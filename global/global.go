package global

import (
	"go-heartbeat/internal/cronjobs/slavemonitor"
	"go-heartbeat/internal/heartbeatconf"
	"time"
)

var HeartbeatSetting heartbeatconf.HeartbeatSettingS
var MasterServerId int
var StartTime time.Time
var SlaveMonitorMsgs = make(slavemonitor.AllMonitorMsgs)
