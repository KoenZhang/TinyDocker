package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

//  挂载了 memory subsystem 的 hierarchy 的根目录位置
const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/memory"

func main() {
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()

		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}

		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		os.Exit(-1)
	}

	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(1)
	} else {
		// 得到 fork 出来进程映射在外部命名空间的pid
		fmt.Printf("%v", cmd.Process.Pid)

		// 在系统默认创建挂在了 memory subsystem 的 Hierarchy 上创建 cgroup
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit"), 0755)
		// 将容器进程加入到这个 cgroup
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "tasks"), []byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		// 限制 cgroup 进程的使用
		ioutil.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmemorylimit", "memory.limit_in_bytes"), []byte("100m"), 0644)

	}

	cmd.Process.Wait()

	//// 指定被 fork 出来的新进程内的初始命令，默认使用 /bin/sh
	//cmd := exec.Command("sh")
	//cmd.SysProcAttr = &syscall.SysProcAttr{
	//	/**
	//	* CLONE_NEWUTS: 创建一个 UTS Namespace, 每个 TS namespace 允许有自己的 hostname
	//	* CLONE_NEWIPC: IPC Namespace用来隔离System V IPC和POSIX message queues
	//	* syscall.CLONE_NEWPID: 隔离进程ID
	//	* syscall.CLONE_NEWNS: Mount Namespace用来隔离各个进程看到的挂载点视图
	//	* syscall.CLONE_NEWUSER: User Namespace 主要是隔离用户的用户组ID
	//	* syscall.CLONE_NEWNET: Network Namespace 是用来隔离网络设备、IP地址端口等网络栈的Namespace
	//	 */
	//	Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	//	UidMappings: []syscall.SysProcIDMap{
	//		{
	//			ContainerID: 0,
	//			HostID:      0,
	//			Size:        1,
	//		},
	//	},
	//	GidMappings: []syscall.SysProcIDMap{
	//		{
	//			ContainerID: 0,
	//			HostID:      0,
	//			Size:        1,
	//		},
	//	},
	//}
	//cmd.Stdin = os.Stdin
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//
	//if err := cmd.Run(); err != nil {
	//	log.Fatal(err)
	//}
	//
	//os.Exit(-1)
}
