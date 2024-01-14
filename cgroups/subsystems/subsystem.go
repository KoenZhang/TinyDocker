package subsystems

// 用于传递资源限制配置的结构体,包括内存限制、CPU时间片权重、CPU核心数
type ResourceConfig struct {
	MemoryLimit string // 内存限制
	CpuShare    string // CPU时间片权重
	CpuSet      string // CPU核心数
}

// Subsystem接口，每个 Subsystem 可以实现以下4个接口
// 这里将 cgroup 抽象为了 path,原因是 cgroup 在 hierarchy 的路径，便是虚拟文件系统中的虚拟路径
type Subsystem interface {
	Name() string                               // 返回 subsystem 的名称
	Set(path string, res *ResourceConfig) error // 设置某个 cgroup 在这个 Subsystem 中的资源限制
	Apply(path string, pid int) error           // 将进程添加到某个 cgroup 中
	Remove(path string) error                   // 移除某个 cgroup
}

// 通过不同的 subsystem 初始化实例创建资源限制处理链数组
var (
	SubsytemIns = []Subsystem{
		&CpusetSubSystem{},
		&MemorySubSystem{},
		&CpuSubSystem{},
	}
)
