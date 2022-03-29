package main

import (
	"go-heartbeat/cmd"
	"go-heartbeat/global"
	"go-heartbeat/internal/heartbeatconf"
	"os"

	log "github.com/sirupsen/logrus"
)

func init() {
	err := setupSetting()
	if err != nil {
		log.Fatalf("init.setupSetting err: %v", err)
	}

	err = setupLogger(global.HeartbeatSetting.LogPath, log.DebugLevel)
	if err != nil {
		log.Fatalf("init.setupLogger err: %v", err)
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

func setupLogger(logPath string, logLevel log.Level) error {
	log.SetFormatter(&log.JSONFormatter{})

	file, err := os.OpenFile(logPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	log.SetOutput(file)
	log.SetLevel(logLevel)

	return nil
}
