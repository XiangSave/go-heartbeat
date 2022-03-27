package serverinit

import "fmt"

func EchoDBInitCmd() {
	// 创建用户
	fmt.Println("CREATE USER '%s' IDENTIFIED BY '%s'")

	//用户配置权限

	// 初始化库表

}
