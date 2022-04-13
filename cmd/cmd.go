package cmd

import (
	"go-heartbeat/global"
	"go-heartbeat/internal/cronjobs/masterupdate"
	"go-heartbeat/internal/cronjobs/slavecheck"
	"go-heartbeat/internal/cronjobs/slaveselect"
	"go-heartbeat/internal/serverinit"
	"go-heartbeat/pkg/cronjob"
	"go-heartbeat/pkg/rolling"

	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"

	"github.com/spf13/cobra"
)

var heartbeatVars struct {
	printDBInitCmd bool
	check          bool
	run            bool
}

var HeartbeatCmd = &cobra.Command{
	Use:   "go-heartbeat",
	Short: "master slave heartbeat",
	Long:  "MySQL 主从延迟监控",
	Run: func(cmd *cobra.Command, args []string) {

		// 基于配置文件打印被监控数据库要执行的语句
		if heartbeatVars.printDBInitCmd {
			serverinit.EchoDBInitCmd()
		}

		// 校验配置文件配置的 master slave 连通性
		if heartbeatVars.check {
			log.Infoln("check master and slave connection status")
		}

		// 开启主从监控
		if heartbeatVars.run {
			jobs := cronjob.New(cron.WithSeconds(), cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger)))

			// 与 Master 建立连接
			con, err := masterupdate.MasterNewConnect()
			if err != nil {
				log.Fatal(err)
			}

			// 记录开启时间，程序运行前 10s 不发送报警
			serverinit.GlobalVarInit()

			// 开启 Master 定时任务
			err = masterupdate.MasterCronRun(jobs, con)
			if err != nil {
				log.Fatal(err)
			}

			// 创建 Slave 相关定时任务
			for _, slaveSetting := range global.HeartbeatSetting.SlaveConnectSetting {

				con, err := slaveselect.SlaveNewConnect(slaveSetting)
				if err != nil {
					log.Fatal(err)
				}
				var maxSaveSecond int64
				for _, s := range slaveSetting.MonitorRole {
					if maxSaveSecond < int64(s.During) {
						maxSaveSecond = int64(s.During)
					}
				}

				log.Infof("%s max save seconds: %d", slaveSetting.Name, maxSaveSecond)
				rollingTime := rolling.NewTiming(maxSaveSecond)
				jobs.AddJob("* * * * * *", slaveselect.SlaveConnectionS{
					Con: con, SlaveSetting: slaveSetting, RollingTiming: rollingTime})

				s := slavecheck.SlaveCheckInit(slaveSetting, rollingTime)
				jobs.AddJob("* * * * * *", slavecheck.SlaveCheck{
					Name: s.Name, MonitorRole: s.MonitorRole, RollingTiming: s.RollingTiming})
			}
			jobs.Run()
		}
	},
}

func Execute() error {
	return HeartbeatCmd.Execute()
}

func init() {
	HeartbeatCmd.Flags().BoolVarP(&heartbeatVars.printDBInitCmd, "print", "p", false, "打印要执行的操作")
	HeartbeatCmd.Flags().BoolVarP(&heartbeatVars.check, "check", "c", false, "校验数据库连接可用性")
	HeartbeatCmd.Flags().BoolVarP(&heartbeatVars.run, "run", "r", false, "启动")
}
