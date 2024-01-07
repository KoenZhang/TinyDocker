package main

import (
	"TinyDocker/container"
	log "github.com/sirupsen/logrus"
	"os"
)

func Run(tty bool, command string) {
	parent := container.NewParentProcess(tty, command)
	/*
		Start 方法是真正开始前面创建好的 command 调用，它首先会clone出来一个namespace隔离的进程，然后在子进程里，调用/proc/self/exec, 也就是调用自己
		发送 init 参数，调用我们写的 init 方法，去初始化容器一些资源
	*/
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	parent.Wait()
	os.Exit(-1)
}
