package subsystems

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
)

// memory subsystem 的实现，实现了 Subsystem 接口，应该实现相应的4个方法：
type MemorySubSystem struct {
}

// 设置 cgroupPath 对应的 cgroup 的内存资源限制
func (s *MemorySubSystem) Set(cgroupPath string, res *ResourceConfig) error {
	//  获取当前 subsystem 在虚拟文件系统中的路径
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if res.MemoryLimit != "" {
			// 设置这个 cgroup 的内存限制，即将限制写入到 cgroup 对应目录的 memory.limit_in_bytes 文件中
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(res.MemoryLimit), 0644); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}

// 将一个进程添加到 cgroupPath 中
func (s *MemorySubSystem) Apply(cgroupPath string, pid int) error {
	//  获取当前 subsystem 在虚拟文件系统中的路径
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		// 把进程的PID写到 cgroup 的虚拟文件系统对应目录下的 "tasks" 文件中
		if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "tasks"), []byte(strconv.Itoa(pid)), 0644); err != nil {
			return fmt.Errorf("set cgroup proc fail %v", err)
		}
		return nil
	} else {
		return fmt.Errorf("get cgroup %s error: %v", cgroupPath, err)
	}
}

// 移除某个 cgroup
func (s *MemorySubSystem) Remove(cgroupPath string) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, false); err == nil {
		// 删除 cgroup 便是删除对应的 cgroupPath 的目录
		return os.RemoveAll(subsysCgroupPath)
	} else {
		return err
	}
}

// 返回 subsystem 的名称
func (s *MemorySubSystem) Name() string {
	return "memory"
}
