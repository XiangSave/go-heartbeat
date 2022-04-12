package heartbeatconf

import "sort"

type MonitorRoleS []MonitorRoleSettingS

type MonitorRoleSettingS struct {
	During       int `yaml:"during"`
	LaterSeconds int `yaml:"laterSeconds"`
}

func (m MonitorRoleS) Len() int           { return len(m) }
func (m MonitorRoleS) Less(i, j int) bool { return m[i].During < m[j].During }
func (m MonitorRoleS) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m MonitorRoleS) Sort()              { sort.Sort(m) }
