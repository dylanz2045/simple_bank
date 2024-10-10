package utils

import (
	"fmt"
	"os"
	"regexp"

	"go.uber.org/zap"
)

func Initlogger() error {
	config := zap.NewDevelopmentConfig()
	config.DisableStacktrace = true
	//检测是否有目录
	r, _ := regexp.Compile(`(.*\/)(.*)`)
	matches := r.FindStringSubmatch("./logs/") //此处用配置文件获取，最后一定要有'/'
	fmt.Println(matches)
	// 检查目录是否存在
	if _, err := os.Stat(matches[1]); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(matches[1], os.ModePerm)
		if err != nil {
			fmt.Printf("无法创建目录：%v\n", err)
			return nil
		}
		fmt.Println("目录已创建：")
	} else {
		// 目录已存在
		fmt.Println("目录已存在：")
	}
	//生成地
	config.OutputPaths = []string{"stdout", "./logs/logs.txt"}
	level := zap.DebugLevel
	//设置日志级别
	config.Level = zap.NewAtomicLevelAt(level)
	config.Encoding = "json"
	Logger, err := config.Build()
	if err != nil {
		return err
	}
	zap.ReplaceGlobals(Logger)
	Logger.Info("-------> Zaplog  init successful ")
	return nil
}
