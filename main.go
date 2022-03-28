package main

import (
	"go-heartbeat/cmd"
	"go-heartbeat/global"
	"go-heartbeat/internal/heartbeatconf"
	"log"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}
}

func main() {
	log.Println(global.HeartbeatSetting)
	err := cmd.Execute()
	if err != nil {
		log.Fatalf("cmd.Execute err: %v", err)
	}
}

func setupSetting() error {
	// confDirpath := "/Users/xiangyan/workSave/esConfigs"
	// confDirpath := "/home/xxx/workSave/configs"
	confDirpath := "./configs"
	setting, err := heartbeatconf.NewSetting(confDirpath)
	if err != nil {
		return err
	}
	err = setting.ReadHeartbeatSetting("GoHeartbeatTools", &global.HeartbeatSetting)
	if err != nil {
		return err
	}

	return nil
}
