package container

import (
	"os"
	"os/exec"
	"syscall"
)

func NewParentProcess(tty bool, command string) *exec.Cmd {
	args := []string{"init", command}
	/*
		指定被 fork 出来的新进程内的初始命令, 默认使用 /proc/self/exec,
		/proc/self/ 指当前运行进程自己的环境，exec 就是自己调用自己，使用这种方式对创建出来的进程进行初始化
		args 是参数，其中 init 是传递给本进程的第一个参数，也就是调用 initCommand 初始化进程的一些环境和资源
	*/
	cmd := exec.Command("/proc/self/exe", args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		/* fork 一个新的进程出来，使用了 namespace 隔离新创建的进程和外部环境
		CLONE_NEWUTS: 创建一个 UTS Namespace, 每个 TS namespace 允许有自己的 hostname
		CLONE_NEWIPC: IPC Namespace用来隔离System V IPC和POSIX message queues
		syscall.CLONE_NEWPID: 隔离进程ID
		syscall.CLONE_NEWNS: Mount Namespace用来隔离各个进程看到的挂载点视图
		syscall.CLONE_NEWUSER: User Namespace 主要是隔离用户的用户组ID
		syscall.CLONE_NEWNET: Network Namespace 是用来隔离网络设备、IP地址端口等网络栈的Namespace
		*/
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWNET | syscall.CLONE_NEWIPC,
	}

	// 如果用户制定了 -ti 参数，就需要把当前进程的输入输出导入到标准输入输出上
	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd
}
