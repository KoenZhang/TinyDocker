package container

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

/*
init 函数是在容器内部执行的，也就是说，代码执行到这里，容器所在的进程已经创建出来了，这是容器执行的第一个进程：
使用mount先挂载 proc 文件系统，以便后面通过 ps 等系统命令去查看当前进程资源情况
*/
func RunContainerInitProcess() error {
	cmdArray := readUserCommand()
	if cmdArray == nil || len(cmdArray) == 0 {
		return fmt.Errorf("Run container get user command error, cmdArray is nil")
	}

	/*
		MS_NOEXEC: 在本文件系统中不允许运行其他程序
		MS_NOSUID： 在本系统运行过程中，不允许 set-user-ID 或者 set-group-ID
		MS_NODEV: 自从 Linux2.4 以来，所有mount的系统都会默认设定的参数
	*/
	//mount proc
	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	// 调用 exec.LookPath，可以在系统的 PATH 里面寻找命令的绝对路径
	path, err := exec.LookPath(cmdArray[0])
	if err != nil {
		log.Errorf("Exec loop path error %v", err)
		return err
	}

	log.Infof("Find path %s", path)

	/*
		首先，使用Docker创建起来一个容器之后，会发现容器内的第一个程序，也就是PID为1的那个进程，是指定的前台进程。那么，根据3.1.1小节所讲的过程发现，
		容器创建之后，执行的第一个进程并不是用户的进程，而是init初始化的进程。这时候，如果通过ps命令查看就会发现，容器内第一个进程变成了自己的init，
		这和预想的是不一样的。你可能会想，大不了把第一个进程给kill掉。但这里又有一个令人头疼的问题，PID 为1的进程是不能被kill掉的，
		如果该进程被kill掉，我们的容器也就退出了。那么，有什么办法呢？这里的execve系统调用就可以大显神威了。
		syscall.Exec这个方法，其实最终调用了Kernel的int execve（const char*filename，char*const argv[]，char*const envp[]）；
		这个系统函数。它的作用是执行当前filename对应的程序。它会覆盖当前进程的镜像、数据和堆栈等信息，包括PID，这些都会被将要运行的进程覆盖掉。
		也就是说，调用这个方法，将用户指定的进程运行起来，把最初的init进程给替换掉，这样当进入到容器内部的时候，就会发现容器内的第一个程序就是我们指定的进程了
	*/
	if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
		log.Errorf(err.Error())
	}

	return nil
}

func readUserCommand() []string {
	// 通过匿名通道建立子进程，一个进程默认有3个文件描述符，分别是：标准输入(0)、标准输出(1)、标准错误(2)、所以新建立的这个匿名通道(文件)就是4(3)
	// uintptr(3) 就是 index 为 3 的文件描述符，也就是传递进来的管道的一端
	pipe := os.NewFile(uintptr(3), "pipe")
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		log.Errorf("init read pipe error %v", err)
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
