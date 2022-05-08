package slavemonitor

import (
	"go-heartbeat/internal/heartbeatconf"
	"sync"
)

type MonitorMsgs struct {
	// 上次告警时间
	LastMonitorTimestamp int16
	// 是否正在出发告警
	Monitoring bool
	// 连续告警次数
	MonitorCount int
	// 规则及延迟信息
	RoleMsgs []MonitorRoleMsg
	// 锁
	Mutex *sync.RWMutex
}

type MonitorRoleMsg struct {
	heartbeatconf.MonitorRoleSettingS
	// 此规则当前是否正在触发告警
	Monitoring bool
	// 当前真实延迟
	NewLaterSecond int64
}
