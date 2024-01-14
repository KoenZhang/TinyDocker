package cgroups

import (
	"TinyDocker/cgroups/subsystems"
	"github.com/sirupsen/logrus"
)

/*
CgroupManager 将资源限制的配置，以及将进程移动到cgroup中的操作交给各个subsystem去处理
将不同的 subsystem 中的 cgroup 管理起来，并与容器建立关系

CgroupManager 在配置容器资源限制时，首先会初始化Subsystem的实例，
然后遍历调用Subsystem实例的Set方法，创建和配置不同Subsystem挂载的hierarchy中的cgroup，
最后再通过调用Subsystem实例将容器的进程分别加入到那些cgroup中，实现容器的资源限制
*/
type CgroupManager struct {
	Path     string                     // cgroup 在 hierarchy 中的路径，相当于创建的 cgroup 目录相对于各 root cgroup 目录的路径
	Resource *subsystems.ResourceConfig // 资源配置
}

func NewCgroupManager(path string) *CgroupManager {
	return &CgroupManager{
		Path: path,
	}
}

// Apply 将进程PID加入到每个cgroup中
func (c *CgroupManager) Apply(pid int) error {
	for _, subSysIns := range subsystems.SubsytemIns {
		subSysIns.Apply(c.Path, pid)
	}
	return nil
}

// Set 设置各个 subsystem 挂载中的 cgroup 资源限制
func (c *CgroupManager) Set(res *subsystems.ResourceConfig) error {
	for _, subSysIns := range subsystems.SubsytemIns {
		subSysIns.Set(c.Path, res)
	}
	return nil
}

// Destory 释放各个 subsystem 挂载中的 cgroup
func (c *CgroupManager) Destory() error {
	for _, subSysIns := range subsystems.SubsytemIns {
		if err := subSysIns.Remove(c.Path); err != nil {
			logrus.Warnf("remove cgroup fail %v", err)
		}
	}
	return nil
}
