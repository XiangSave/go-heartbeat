package cmd

import (
	"go-heartbeat/internal/cronjobs/masterupdate"
	"go-heartbeat/pkg/cronjob"

	log "github.com/sirupsen/logrus"

	"github.com/robfig/cron/v3"
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
			log.Info("grants users ......")
		}

		// 校验配置文件配置的 master slave 连通性
		if heartbeatVars.check {
			log.Infoln("check master and slave connection status")
		}

		// 开启主从监控
		if heartbeatVars.run {
			jobs := cronjob.New(cron.WithSeconds(), cron.WithChain(
				cron.SkipIfStillRunning(cron.DefaultLogger)))

			_, err := jobs.AddFunc("* * * * * *", masterupdate.MasterUpdate)
			if err != nil {
				log.Fatal(err)
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
