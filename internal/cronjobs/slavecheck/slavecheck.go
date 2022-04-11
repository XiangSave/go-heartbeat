package slavecheck

import (
	"go-heartbeat/internal/heartbeatconf"
	"go-heartbeat/pkg/rolling"
)

type SlaveCheck struct {
	SlaveSetting  heartbeatconf.SlaveConnectSettingS
	RollingTiming *rolling.Timing
}
