package util

import (
	"fmt"
	"log"
	"os"
)

// 处理基本的cli
func DealBaseCli() bool {
	// 设置框架logo
	setLogo();
	// 验证cli参数是否正确
	if !checkCliParams() {
		os.Exit(0)
	}
	return true
}

// 验证cli参数
func checkCliParams() bool {
	if len(os.Args) != 2 {
		log.Printf("%v", "错误的cli参数(Missing cli parameters)")
		return false
	}
	switch os.Args[1] {
	case "start" :
		startCli()
		fmt.Printf("开源地址:https://github.com/MicroSocketTeam/deployMicroSocket/tree/master\n\r" +
			"运行正常...\n\r")
		break
	default :
		log.Printf("%v", "错误的参数(Missing parameters)")
		return false
	}
	return true
}

// cli参数验证成功后需要执行的操作
func startCli()  {

}

// cli参数验证失败后需要执行的操作
func stopCli()  {

}

// 框架Logo
func setLogo()  {
	fmt.Printf("%v",
		"==================================================\n\r" +
			"                   microSocket\n\r" +
			"==================================================\n\r")
}

