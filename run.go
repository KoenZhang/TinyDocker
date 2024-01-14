package main

import (
	"TinyDocker/cgroups"
	"TinyDocker/cgroups/subsystems"
	"TinyDocker/container"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func Run(tty bool, comArray []string, res *subsystems.ResourceConfig) {
	parent, writePipe := container.NewParentProcess(tty)
	if parent == nil {
		log.Errorf("New parent process error")
		return
	}
	/*
		Start 方法是真正开始前面创建好的 command 调用，它首先会clone出来一个namespace隔离的进程，然后在子进程里，调用/proc/self/exec, 也就是调用自己
		发送 init 参数，调用我们写的 init 方法，去初始化容器一些资源
	*/
	if err := parent.Start(); err != nil {
		log.Error(err)
	}

	// 使用 mydocer-cgroup 作为 cgroup name
	// 创建 cgroup manager，并通过调用 set 和 apply 设置资源限制并限制在容器上生效
	cgroupManager := cgroups.NewCgroupManager("mydocker-cgroup")
	defer cgroupManager.Destory()

	// 设置资源限制
	cgroupManager.Set(res)
	// 将容器进程加入到各个 subsystem 挂载的对应 cgroup 中
	cgroupManager.Apply(parent.Process.Pid)
	// 对容器设置完限制后，初始化容器
	sendInitCommand(comArray, writePipe)

	parent.Wait()
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	log.Infof("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}
