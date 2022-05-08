package slavecheck

import (
	"go-heartbeat/global"
	"go-heartbeat/internal/heartbeatconf"
	"go-heartbeat/pkg/rolling"
	"sync"

	"time"

	"github.com/pkg/errors"

	log "github.com/sirupsen/logrus"
)

type SlaveCheck struct {
	Name          string
	MonitorRole   []heartbeatconf.MonitorRoleSettingS
	RollingTiming *rolling.Timing
}

func SlaveCheckInit(c heartbeatconf.SlaveConnectSettingS, r *rolling.Timing) *SlaveCheck {
	c.MonitorRole.Sort()

	return &SlaveCheck{
		Name:          c.Name,
		MonitorRole:   c.MonitorRole,
		RollingTiming: r,
	}
}

func (s SlaveCheck) getRollingAverage(during int) (int64, error) {
	return s.RollingTiming.RangeAverage(during)
}

func (s SlaveCheck) setGlobalMonitorMsgs(name string, during int, later int64) error {
	global.SlaveMonitorMsgs[name].Mutex.Lock()
	defer global.SlaveMonitorMsgs[name].Mutex.Unlock()
	for _, gmr := range global.SlaveMonitorMsgs[name].RoleMsgs {
		if gmr.During == during {
			gmr.NewLaterSecond = later
		}
	}

	return nil
}

func (s SlaveCheck) Run() {
	wg := sync.WaitGroup{}
	wg.Add(len(s.MonitorRole))

	for i := 0; i < len(s.MonitorRole); i++ {
		go func(i int) {
			role := s.MonitorRole[i]
			later, err := s.getRollingAverage(role.During)

			if err != nil && !errors.Is(err, rolling.ErrNoMatchLens) ||
				(err != nil && i == len(s.MonitorRole) && time.Now().Sub(global.StartTime) > 15*1000000000) {
				log.Info(len(s.RollingTiming.Buckets))
				log.Errorf("%+v", err)
			}

			if later >= int64(role.LaterSeconds) {
				log.Warnf("%s 状态异常报警：时间区间：%d，告警阈值：%d，当前延迟:%d",
					s.Name, role.During, role.LaterSeconds, later)
			} else {
				log.Infof("%s 状态正常报警：时间区间：%d，告警阈值：%d，当前延迟:%d",
					s.Name, role.During, role.LaterSeconds, later)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
}
