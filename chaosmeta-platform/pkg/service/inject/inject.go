package inject

import (
	models "chaosmeta-platform/pkg/models/common"
	"chaosmeta-platform/pkg/models/inject/basic"
	"context"
	"fmt"
)

type InjectService struct{}

type ScopeType string

const (
	PodScopeType        ScopeType = "pod"
	NodeScopeType       ScopeType = "node"
	KubernetesScopeType ScopeType = "kubernetes"

	ExecInject = "inject"
)

var (
	PodScope        = basic.Scope{Name: string(PodScopeType), NameCn: "容器组", Description: "fault injection can be performed on any Pod in the Kubernetes cluster", DescriptionCn: "可以对 Kubernetes 集群中的任意 Pod 进行故障注入"}
	NodeScope       = basic.Scope{Name: string(NodeScopeType), NameCn: "节点", Description: "fault injection can be performed on any Node in the Kubernetes cluster", DescriptionCn: "可以对 Kubernetes 集群中的任意 Node 进行故障注入"}
	KubernetesScope = basic.Scope{Name: string(KubernetesScopeType), NameCn: "节点", Description: "faults can be injected into Kubernetes resource instances such as pod, deployment, and node to achieve the exception of the Kubernetes cluster itself or the exception of the operator application", DescriptionCn: "可以对 pod、deployment、node 等Kubernetes资源实例注入故障，达到 Kubernetes 集群自身的异常或者 operator 应用的异常"}
)

func Init() error {
	scope, target, fault, args := basic.Scope{}, basic.Target{}, basic.Fault{}, basic.Args{}
	if _, err := models.GetORM().Raw(fmt.Sprintf("TRUNCATE TABLE %s", scope.TableName())).Exec(); err != nil {
		return err
	}

	if _, err := models.GetORM().Raw(fmt.Sprintf("TRUNCATE TABLE %s", target.TableName())).Exec(); err != nil {
		return err
	}

	if _, err := models.GetORM().Raw(fmt.Sprintf("TRUNCATE TABLE %s", fault.TableName())).Exec(); err != nil {
		return err
	}

	if _, err := models.GetORM().Raw(fmt.Sprintf("TRUNCATE TABLE %s", args.TableName())).Exec(); err != nil {
		return err
	}
	ctx := context.Background()

	// 开始初始化数据
	scopes := []basic.Scope{PodScope, NodeScope}
	for _, scope := range scopes {
		scopeId, err := basic.InsertScope(ctx, &scope)
		if err != nil {
			return err
		}
		scope.ID = int(scopeId)
		if err := InitTarget(ctx, scope); err != nil {
			return err
		}

	}

	if _, err := basic.InsertScope(ctx, &KubernetesScope); err != nil {
		return err
	}
	return InitK8STarget(ctx, KubernetesScope)
}

func InitTarget(ctx context.Context, scope basic.Scope) error {
	var (
		CpuTarget       = basic.Target{Name: "cpu", NameCn: "cpu", Description: "fault injection capabilities related to cpu faults", DescriptionCn: "cpu故障相关的故障注入能力"}
		MemTarget       = basic.Target{Name: "mem", NameCn: "mem", Description: "fault injection capabilities related to memory faults", DescriptionCn: "内存故障相关的故障注入能力"}
		DiskTarget      = basic.Target{Name: "disk", NameCn: "disk", Description: "fault injection capabilities related to disk failures", DescriptionCn: "磁盘故障相关的故障注入能力"}
		DiskioTarget    = basic.Target{Name: "diskio", NameCn: "diskio", Description: "fault injection capabilities related to disk IO faults", DescriptionCn: "磁盘IO故障相关的故障注入能力"}
		NetworkTarget   = basic.Target{Name: "network", NameCn: "network", Description: "fault injection capabilities related to disk failures", DescriptionCn: "磁盘故障相关的故障注入能力"}
		ProcessTarget   = basic.Target{Name: "process", NameCn: "process", Description: "process-dependent fault injection capability", DescriptionCn: "进程相关的故障注入能力"}
		FileTarget      = basic.Target{Name: "file", NameCn: "file", Description: "file-related fault injection capabilities", DescriptionCn: "文件相关的故障注入能力"}
		KernelTarget    = basic.Target{Name: "kernel", NameCn: "kernel", Description: "kernel-related fault injection capabilities", DescriptionCn: "内核相关的故障注入能力"}
		JvmTarget       = basic.Target{Name: "jvm", NameCn: "jvm", Description: "jvm related fault injection capabilities", DescriptionCn: "围绕jvm相关的故障注入能力"}
		ContainerTarget = basic.Target{Name: "container", NameCn: "container", Description: "container runtime-related fault injection capabilities", DescriptionCn: "容器运行时相关的故障注入能力"}
	)

	CpuTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &CpuTarget); err != nil {
		return err
	}
	if err := InitCpuFault(ctx, CpuTarget); err != nil {
		return err
	}
	MemTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &MemTarget); err != nil {
		return err
	}
	if err := InitMemFault(ctx, CpuTarget); err != nil {
		return err
	}
	DiskTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &DiskTarget); err != nil {
		return err
	}
	if err := InitDiskFault(ctx, DiskTarget); err != nil {
		return err
	}
	DiskioTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &DiskioTarget); err != nil {
		return err
	}
	if err := InitDiskioFault(ctx, DiskioTarget); err != nil {
		return err
	}
	NetworkTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &NetworkTarget); err != nil {
		return err
	}
	if err := InitNetworkFault(ctx, NetworkTarget); err != nil {
		return err
	}
	ProcessTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &ProcessTarget); err != nil {
		return err
	}
	if err := InitProcessFault(ctx, ProcessTarget); err != nil {
		return err
	}
	FileTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &FileTarget); err != nil {
		return err
	}
	if err := InitFileFault(ctx, FileTarget); err != nil {
		return err
	}
	KernelTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &KernelTarget); err != nil {
		return err
	}
	if err := InitKernelFault(ctx, KernelTarget); err != nil {
		return err
	}
	JvmTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &JvmTarget); err != nil {
		return err
	}
	if err := InitJvmFault(ctx, JvmTarget); err != nil {
		return err
	}
	ContainerTarget.ScopeId = scope.ID
	if err := basic.InsertTarget(ctx, &ContainerTarget); err != nil {
		return err
	}
	return InitContainerFault(ctx, ContainerTarget)
}

func InitCpuFault(ctx context.Context, cpuTarget basic.Target) error {
	var (
		CpuFaultBurn = basic.Fault{TargetId: cpuTarget.ID, Name: "burn", NameCn: "cpu燃烧", Description: "high cpu usage", DescriptionCn: "cpu使用率飙高"}
		CpuFaultLoad = basic.Fault{TargetId: cpuTarget.ID, Name: "burn", NameCn: "cpu燃烧", Description: "high cpu usage", DescriptionCn: "cpu使用率飙高"}
	)
	if err := basic.InsertFault(ctx, &CpuFaultBurn); err != nil {
		return err
	}
	if err := InitCpuTargetArgsBurn(ctx, CpuFaultBurn); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &CpuFaultLoad); err != nil {
		return err
	}
	return InitLoadTargetArgsLoad(ctx, CpuFaultLoad)
}

func InitCpuTargetArgsBurn(ctx context.Context, cpuFault basic.Fault) error {
	var (
		CpuArgsPercent = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "使用率", Unit: "numbers from 1-100", UnitCn: "1-100的数字", Description: "target cpu usage", DescriptionCn: "目标cpu使用率", ValueType: "int"}
		CpuArgsCount   = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "核数", Unit: "an integer greater than or equal to 0, 0 means all cores", UnitCn: "大于等于0的整数，0表示全部核", DefaultValue: "0", Description: "number of faulty CPU cores", DescriptionCn: "故障cpu核数", ValueType: "int"}
		CpuArgsList    = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "list", KeyCn: "列表", Unit: "comma separated list of core numbers, can be confirmed from /proc/cpuinfo", UnitCn: "逗号分隔的核编号列表，可以从/proc/cpuinfo确认", Description: "cpu fault list", DescriptionCn: "故障cpu列表", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&CpuArgsPercent, &CpuArgsCount, &CpuArgsList})
}

func InitLoadTargetArgsLoad(ctx context.Context, loadFault basic.Fault) error {
	var LoadArgsCount = basic.Args{InjectId: loadFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "负载数", Unit: "an integer greater than or equal to 0, 0 means 4 times the number of CPU cores", UnitCn: "大于等于0的整数，如果为0表示“cpu核数的4倍”", Description: "number of loads to add", DescriptionCn: "需要增加的负载数", ValueType: "int"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&LoadArgsCount})
}

func InitMemFault(ctx context.Context, memTarget basic.Target) error {
	var (
		MemFaultFill = basic.Fault{TargetId: memTarget.ID, Name: "fill", NameCn: "内存填充", Description: "memory usage soars", DescriptionCn: "内存使用率飙高"}
		MemFaultOom  = basic.Fault{TargetId: memTarget.ID, Name: "oom", NameCn: "内存oom", Description: "the system memory oom will cause the machine to hang up", DescriptionCn: "系统内存oom，会使机器宕机挂掉"}
	)

	if err := basic.InsertFault(ctx, &MemFaultFill); err != nil {
		return err
	}
	if err := InitMemTargetArgsFill(ctx, MemFaultFill); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &MemFaultOom); err != nil {
		return err
	}

	return InitMemTargetArgsOom(ctx, MemFaultOom)
}

func InitMemTargetArgsFill(ctx context.Context, memFault basic.Fault) error {
	var (
		MemArgsPercent = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "使用率", Unit: "numbers from 1-100", UnitCn: "1-100的数字", Description: "target mem usage", DescriptionCn: "目标内存使用率", ValueType: "int"}
		MemArgsBytes   = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "bytes", KeyCn: "填充量", Unit: "units: KB (default), MB, GB, TB", UnitCn: "单位：KB（默认）、MB、GB、TB", DefaultValue: "KB", Description: "memory fill", DescriptionCn: "内存填充量", ValueType: "string"}
		MemArgsMode    = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "填充模式", Unit: "optional: ram, cache.ram is the way to use process memory application, cache is the way to use tmpfs", UnitCn: "可选：ram、cache。ram是使用进程内存申请的方式，cache是使用tmpfs的方式", Description: "memory fill mode", DescriptionCn: "内存填充模式", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&MemArgsPercent, &MemArgsBytes, &MemArgsMode})
}

func InitMemTargetArgsOom(ctx context.Context, memFault basic.Fault) error {
	var MemArgsMode = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "填充模式", Unit: "optional: ram, cache.ram is the way to use process memory application, cache is the way to use tmpfs", UnitCn: "可选：ram、cache。ram是使用进程内存申请的方式，cache是使用tmpfs的方式", DefaultValue: "cache", Description: "memory fill mode", DescriptionCn: "内存填充模式", ValueType: "string"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&MemArgsMode})
}

func InitDiskFault(ctx context.Context, diskTarget basic.Target) error {
	var DiskFaultFill = basic.Fault{TargetId: diskTarget.ID, Name: "fill", NameCn: "磁盘填充", Description: "disk usage soars", DescriptionCn: "磁盘使用率飙高"}
	if err := basic.InsertFault(ctx, &DiskFaultFill); err != nil {
		return err
	}
	return InitDiskTargetArgsFill(ctx, DiskFaultFill)
}

func InitDiskTargetArgsFill(ctx context.Context, diskFault basic.Fault) error {
	var (
		DiskArgsPercent = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "使用率", Unit: "numbers from 1-100", UnitCn: "1-100的数字", Description: "target disk usage", DescriptionCn: "目标磁盘使用率", ValueType: "int"}
		DiskArgsBytes   = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "bytes", KeyCn: "填充量", Unit: "units: KB (default), MB, GB, TB", UnitCn: "单位：KB（默认）、MB、GB、TB", DefaultValue: "KB", Description: "disk fill", DescriptionCn: "磁盘填充量", ValueType: "string"}
		DiskArgsDir     = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "dir", KeyCn: "目录", Unit: "target directory", UnitCn: "目标目录", DefaultValue: "/tmp", Description: "target population directory", DescriptionCn: "目标填充目录", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskArgsPercent, &DiskArgsBytes, &DiskArgsDir})
}

func InitDiskioFault(ctx context.Context, diskioTarget basic.Target) error {
	var (
		DiskioFaultBurn  = basic.Fault{TargetId: diskioTarget.ID, Name: "burn", NameCn: "磁盘IO高负载", Description: "disk IO soars", DescriptionCn: "磁盘IO飙高"}
		DiskioFaultHang  = basic.Fault{TargetId: diskioTarget.ID, Name: "hang", NameCn: "磁盘IO hang", Description: "the target process generates a disk IO hang", DescriptionCn: "目标进程产生磁盘IO hang"}
		DiskioFaultLimit = basic.Fault{TargetId: diskioTarget.ID, Name: "limit", NameCn: "磁盘IO limit", Description: "target process generates disk IO limit", DescriptionCn: "目标进程产生磁盘IO limit"}
	)
	if err := basic.InsertFault(ctx, &DiskioFaultBurn); err != nil {
		return err
	}
	if err := InitDiskioTargetArgsBurn(ctx, DiskioFaultBurn); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &DiskioFaultHang); err != nil {
		return err
	}
	if err := InitDiskioTargetArgsHang(ctx, DiskioFaultHang); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &DiskioFaultLimit); err != nil {
		return err
	}
	return InitDiskioTargetArgsLimit(ctx, DiskioFaultLimit)
}

func InitDiskioTargetArgsBurn(ctx context.Context, diskioFault basic.Fault) error {
	var (
		DiskioArgsDir   = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dir", KeyCn: "目录", Unit: "target directory", UnitCn: "目标目录", DefaultValue: "/tmp", Description: "target directory for high IO operations", DescriptionCn: "进行高IO操作的目标目录", ValueType: "string"}
		DiskioArgsMode  = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "IO模式", Unit: "read、write", UnitCn: "读、写", DefaultValue: "read", Description: "IO mode", DescriptionCn: "IO模式", ValueType: "string"}
		DiskioArgsBlock = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "block", KeyCn: "IO块", Unit: "the range is 1K-1024M, and the optional units are: KB (default), MB", UnitCn: "范围为1K-1024M，可选单位为：KB（默认）、MB", DefaultValue: "10MB", Description: "the block size of a single IO", DescriptionCn: "单次IO的块大小", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDir, &DiskioArgsMode, &DiskioArgsBlock})
}

func InitDiskioTargetArgsHang(ctx context.Context, diskioFault basic.Fault) error {
	var (
		DiskioArgsDevList = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dev-list", KeyCn: "设备列表", Unit: "use the command lsblk -a | grep disk to get the major and minor device numbers of the target device. For example: 8:0", UnitCn: "使用命令lsblk -a | grep disk获取目标设备的主备设备号。比如：8:0", Description: "target disk device list", DescriptionCn: "目标磁盘设备列表", ValueType: "string"}
		DiskioArgsPidList = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "pid-list", KeyCn: "进程pid列表", Unit: "For example: 7850,7690", UnitCn: "比如：7850,7690", Description: "affected process pid list", DescriptionCn: "受影响的进程pid列表", ValueType: "string"}
		DiskioArgsKey     = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Unit: "will use ps -ef | grep [key] to filter", UnitCn: "会使用ps -ef | grep [key]来筛选", Description: "keyword used to filter affected processes", DescriptionCn: "用来筛选受影响进程的关键词", ValueType: "string"}
		DiskioArgsMode    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "IO模式", Unit: "all、read、write", UnitCn: "受影响的IO操作", DefaultValue: "all", Description: "affected IO operation", DescriptionCn: "受影响的IO操作", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDevList, &DiskioArgsPidList, &DiskioArgsKey, &DiskioArgsMode})
}

func InitDiskioTargetArgsLimit(ctx context.Context, diskioFault basic.Fault) error {
	var (
		DiskioArgsDir        = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dir", KeyCn: "目录", Unit: "target directory", UnitCn: "目标目录", DefaultValue: "/tmp", Description: "target directory for high IO operations", DescriptionCn: "进行高IO操作的目标目录", ValueType: "string"}
		DiskioArgsMode       = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "IO模式", Unit: "read、write", UnitCn: "读、写", DefaultValue: "read", Description: "IO mode", DescriptionCn: "IO模式", ValueType: "string"}
		BlockArgsMode        = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "block", KeyCn: "IO块", Unit: "the range is 1K-1024M, and the optional units are: KB (default), MB", UnitCn: "范围为1K-1024M，可选单位为：KB（默认）、MB", DefaultValue: "10MB", Description: "the block size of a single IO", DescriptionCn: "单次IO的块大小", ValueType: "string"}
		DiskioArgsReadBytes  = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "read-bytes", KeyCn: "读字节数", Unit: "optional unit: B (default), KB, MB, GB, TB, for example: 1KB", UnitCn: "可选单位：B（默认）、KB、MB、GB、TB，比如：1KB", Description: "number of bytes that can be read per second", DescriptionCn: "每秒能读的字节数", ValueType: "string"}
		DiskioArgsWriteBytes = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "write-bytes", KeyCn: "写字节数", Unit: "optional unit: B (default), KB, MB, GB, TB, for example: 1KB", UnitCn: "可选单位：B（默认）、KB、MB、GB、TB，比如：1KB", Description: "number of bytes that can be written per second", DescriptionCn: "每秒能写的字节数", ValueType: "string"}
		DiskioArgsReadIO     = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "read-io", KeyCn: "读IO次数", Unit: "positive integer", UnitCn: "大于0的整数", Description: "number of IO operations that can be read per second", DescriptionCn: "每秒能读的IO次数", ValueType: "int"}
		DiskioArgsWriteIO    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "write-io", KeyCn: "写IO次数", Unit: "positive integer", UnitCn: "大于0的整数", Description: "number of IO operations that can be written per second", DescriptionCn: "每秒能写的IO次数", ValueType: "int"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDir, &DiskioArgsMode, &BlockArgsMode, &DiskioArgsReadBytes, &DiskioArgsWriteBytes, &DiskioArgsReadIO, &DiskioArgsWriteIO})
}

func InitNetworkFault(ctx context.Context, networkTarget basic.Target) error {
	var (
		NetworkFaultOccupy  = basic.Fault{TargetId: networkTarget.ID, Name: "occupy", NameCn: "端口占用", Description: "kill the service of the target port and occupy this port", DescriptionCn: "把目标端口的服务kill掉并把这个端口占用掉"}
		NetworkFaultLimit   = basic.Fault{TargetId: networkTarget.ID, Name: "limit", NameCn: "带宽限制", Description: "bandwidth limit for network packet injection in the outbound direction from the failed machine", DescriptionCn: "对从故障机器流出方向的网络数据包注入带宽限制"}
		NetworkFaultDelay   = basic.Fault{TargetId: networkTarget.ID, Name: "delay", NameCn: "网络延迟", Description: "network packet injection delay in outbound direction from the failed machine", DescriptionCn: "从故障机器流出方向的网络数据包注入延迟"}
		NetworkFaultLoss    = basic.Fault{TargetId: networkTarget.ID, Name: "loss", NameCn: "网络丢包", Description: "network packet injection in outbound direction from faulty machine drops packets", DescriptionCn: "从故障机器流出方向的网络数据包注入丢包"}
		NetworkFaultCorrupt = basic.Fault{TargetId: networkTarget.ID, Name: "corrupt", NameCn: "网络包损坏", Description: "inject packet corruption on network packets in the outbound direction from the failed machine", DescriptionCn: "对从故障机器流出方向的网络数据包注入包损坏"}

		NetworkFaultDuplicate = basic.Fault{TargetId: networkTarget.ID, Name: "duplicate", NameCn: "网络包重复", Description: "repeat for network packet injection packets in the outbound direction from the failed machine", DescriptionCn: "对从故障机器流出方向的网络数据包注入包重复"}
		NetworkFaultReorder   = basic.Fault{TargetId: networkTarget.ID, Name: "reorder", NameCn: "网络包乱序", Description: "inject out-of-sequence packets into network packets in the outflow direction from the faulty machine", DescriptionCn: "对从故障机器流出方向的网络数据包注入包乱序"}
	)

	if err := basic.InsertFault(ctx, &NetworkFaultOccupy); err != nil {
		return err
	}
	if err := InitNetworkTargetArgsOccupy(ctx, NetworkFaultOccupy); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &NetworkFaultLimit); err != nil {
		return err
	}
	if err := InitNetworkTargetArgsLimit(ctx, NetworkFaultLimit); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &NetworkFaultDelay); err != nil {
		return err
	}
	if err := InitNetworkTargetArgsDelay(ctx, NetworkFaultDelay); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &NetworkFaultLoss); err != nil {
		return err
	}
	if err := InitNetworkTargetArgs(ctx, NetworkFaultLoss); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &NetworkFaultCorrupt); err != nil {
		return err
	}
	if err := InitNetworkTargetArgs(ctx, NetworkFaultCorrupt); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &NetworkFaultDuplicate); err != nil {
		return err
	}
	if err := InitNetworkTargetArgs(ctx, NetworkFaultDuplicate); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &NetworkFaultReorder); err != nil {
		return err
	}
	return InitNetworkTargetReorder(ctx, NetworkFaultReorder)
}

func getNetworkCommonFilterParameters(fault basic.Fault) []*basic.Args {
	var (
		NetworkArgsInterface = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "interface", KeyCn: "网卡", UnitCn: "故障机器包含的网卡，比如：eth0", Unit: "the network card included in the faulty machine, such as: eth0", Description: "network interface", DescriptionCn: "网卡", ValueType: "string"}
		NetworkArgsDstIP     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "dst-ip", KeyCn: "数据包筛选参数：目标ip列表", UnitCn: "比如：1.2.3.4,2.3.4.5,192.168.1.1/24", Unit: "for example: 1.2.3.4, 2.3.4.5, 192.168.1.1/24", Description: "packet filtering parameter: destination IP list", DescriptionCn: "数据包筛选参数：目标ip列表", ValueType: "string"}
		NetworkArgsSrcIP     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "src-ip", KeyCn: "数据包筛选参数：源ip列表", UnitCn: "比如：1.2.3.4,2.3.4.5,192.168.1.1/24", Unit: "for example: 1.2.3.4,2.3.4.5,192.168.1.1/24", Description: "packet filtering parameter: source IP list", DescriptionCn: "数据包筛选参数：源ip列表", ValueType: "string"}
		NetworkArgsDstPort   = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "dst-port", KeyCn: "数据包筛选参数：目标端口列表", UnitCn: "比如：8080-8090,8095,9099", Unit: "for example: 8080-8090,8095,9099", Description: "packet filtering parameter: destination port list", DescriptionCn: "数据包筛选参数：目标端口列表", ValueType: "string"}
		NetworkArgsSrcPort   = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "src-port", KeyCn: "数据包筛选参数：源端口列表", UnitCn: "比如：8080-8090,8095,9099", Unit: "for example: 8080-8090,8095,9099", Description: "packet filtering parameter: source port list", DescriptionCn: "数据包筛选参数：源端口列表", ValueType: "string"}
		NetworkArgsMode      = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "数据", Unit: "normal: inject fault to selected targets, exclude: do not inject fault to selected targets", UnitCn: "normal：对选中的目标注入故障 exclude：对选中的目标不注入故障", DefaultValue: "normal", Description: "packet filtering mode", DescriptionCn: "数据包筛选模式", ValueType: "string"}
		NetworkArgsForce     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "force", KeyCn: "是否强制覆盖", Unit: "true or false", UnitCn: "是、否", DefaultValue: "false", Description: "whether to force overwrite", DescriptionCn: "是否强制覆盖", ValueType: "bool"}
	)
	return []*basic.Args{&NetworkArgsInterface, &NetworkArgsDstIP, &NetworkArgsSrcIP, &NetworkArgsDstPort, &NetworkArgsSrcPort, &NetworkArgsMode, &NetworkArgsForce}
}

func InitNetworkTargetArgsOccupy(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsPort       = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "port", KeyCn: "端口", Unit: "本机端口号", UnitCn: "local port number", Description: "target port", DescriptionCn: "目标端口", ValueType: "int"}
		NetworkArgsProtocol   = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "protocol", KeyCn: "协议", Unit: "optional protocols: tcp, udp, tcp6, udp6", UnitCn: "可选协议：tcp、udp、tcp6、udp6", DefaultValue: "tcp", Description: "target protocol", DescriptionCn: "目标协议", ValueType: "string"}
		NetworkArgsRecoverCmd = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "recover-cmd", KeyCn: "恢复命令", Unit: "user provided", UnitCn: "用户提供", Description: "resume command, it will be executed last when resuming operation", DescriptionCn: "恢复命令，恢复操作时会最后执行", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&NetworkArgsPort, &NetworkArgsProtocol, &NetworkArgsRecoverCmd})
}

func InitNetworkTargetArgsLimit(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsPort       = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "port", KeyCn: "端口", UnitCn: "本机端口号", Unit: "local port number", Description: "target port", DescriptionCn: "目标端口", ValueType: "int"}
		NetworkArgsProtocol   = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "protocol", KeyCn: "协议", UnitCn: "optional protocols: tcp, udp, tcp6, udp6", Unit: "可选协议：tcp、udp、tcp6、udp6", DefaultValue: "tcp", Description: "target protocol", DescriptionCn: "目标协议", ValueType: "string"}
		NetworkArgsRecoverCmd = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "recover-cmd", KeyCn: "恢复命令", UnitCn: "user provided", Unit: "用户提供", Description: "resume command, it will be executed last when resuming operation", DescriptionCn: "恢复命令，恢复操作时会最后执行", ValueType: "string"}
		NetworkArgsRate       = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "rate", KeyCn: "带宽限制", Unit: "support unit: bit (default), kbit, mbit, gbit, tbit, for example: 10bit", UnitCn: "支持单位：bit（默认）、kbit、mbit、gbit、tbit，比如：10bit", Description: "network bandwidth limit per second", DescriptionCn: "每秒的网络带宽限制", ValueType: "string"}
	)
	argList := []*basic.Args{&NetworkArgsPort, &NetworkArgsProtocol, &NetworkArgsRecoverCmd, &NetworkArgsRate}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsDelay(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsLatency = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "latency", KeyCn: "延迟时间", Unit: "optional units: s, ms, us (default), for example: 3s", UnitCn: "可选单位：s、ms、us（默认），比如：3s", Description: "delay", DescriptionCn: "延迟时间", ValueType: "string"}
		NetworkArgsJitter  = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "jitter", KeyCn: "抖动值", Unit: "optional units: s, ms, us (default), for example: 3s", UnitCn: "可选单位：s、ms、us（默认），比如：3s", Description: "jitter value, the fluctuation range of each delay", DescriptionCn: "抖动值，每次延迟的波动范围", DefaultValue: "0", ValueType: "string"}
	)

	argList := []*basic.Args{&NetworkArgsLatency, &NetworkArgsJitter}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgs(ctx context.Context, networkFault basic.Fault) error {
	var NetworkArgsPercent = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "丢包率", Unit: "an integer conforming to (0, 100]", UnitCn: "符合(0, 100]的整数", Description: "packet loss rate", DescriptionCn: "丢包率", ValueType: "int"}

	argList := []*basic.Args{&NetworkArgsPercent}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetReorder(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsLatency = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "latency", KeyCn: "延迟时间", Unit: "optional units: s, ms, us (default), for example: 3s", UnitCn: "可选单位：s、ms、us（默认），比如：3s", DefaultValue: "1s", Description: "delay", DescriptionCn: "延迟时间", ValueType: "string"}
		NetworkArgsGap     = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "gap", KeyCn: "KeyCn", Unit: "an integer greater than 0", UnitCn: "大于0的整数", DefaultValue: "3", Description: "select the interval, for example, if the gap is 3, it means that the packets with serial numbers 1, 3, 6, 9, etc. will not be delayed, and the rest of the packets will be delayed", DescriptionCn: "选中间隔，比如gap为3表示序号为1、3、6、9等的包不延迟，其余的包会延迟", ValueType: "string"}
	)
	argList := []*basic.Args{&NetworkArgsLatency, &NetworkArgsGap}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitProcessFault(ctx context.Context, processTarget basic.Target) error {
	var (
		processFaultKill = basic.Fault{TargetId: processTarget.ID, Name: "kill", NameCn: "杀进程", Description: "kill the target process", DescriptionCn: "把目标进程杀掉"}

		processFaultStop = basic.Fault{TargetId: processTarget.ID, Name: "stop", NameCn: "停止进程", Description: "stop target process", DescriptionCn: "停止目标进程"}
	)
	if err := basic.InsertFault(ctx, &processFaultKill); err != nil {
		return err
	}
	if err := InitProcessTargetArgsKill(ctx, processFaultKill); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &processFaultStop); err != nil {
		return err
	}
	return InitProcessTargetArgsStop(ctx, processFaultStop)
}

func InitProcessTargetArgsKill(ctx context.Context, processFault basic.Fault) error {
	var (
		ProcessArgsKey        = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Unit: "Will use ps -ef | grep [key] to filter", UnitCn: "会使用ps -ef | grep [key]来筛选", Description: "keywords used to filter affected processes", DescriptionCn: "用来筛选受影响进程的关键词", ValueType: "string"}
		ProcessArgsPid        = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程", Unit: "the pid of the living process", UnitCn: "存活进程的pid", Description: "process pid", DescriptionCn: "进程pid", ValueType: "string"}
		ProcessArgsSignal     = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "signal", KeyCn: "信号", Unit: "consistent with the signal value supported by the kill command", UnitCn: "和kill命令支持的信号数值一致", DefaultValue: "9", Description: "the signal sent to the process", DescriptionCn: "对进程发送的信号", ValueType: "string"}
		ProcessArgsRecoverCmd = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "recover-cmd", KeyCn: "恢复命令", Unit: "user provided", UnitCn: "用户提供", Description: "resume command, it will be executed last when resuming operation", DescriptionCn: "恢复命令，恢复操作时会最后执行", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&ProcessArgsKey, &ProcessArgsPid, &ProcessArgsSignal, &ProcessArgsRecoverCmd})
}

func InitProcessTargetArgsStop(ctx context.Context, processFault basic.Fault) error {
	var (
		ProcessArgsKey = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Unit: "Will use ps -ef | grep [key] to filter", UnitCn: "会使用ps -ef | grep [key]来筛选", Description: "keywords used to filter affected processes", DescriptionCn: "用来筛选受影响进程的关键词", ValueType: "string"}
		ProcessArgsPid = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程", Unit: "the pid of the living process", UnitCn: "存活进程的pid", Description: "process pid", DescriptionCn: "进程pid", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&ProcessArgsKey, &ProcessArgsPid})
}

func InitFileFault(ctx context.Context, fileTarget basic.Target) error {
	var (
		fileFaultChmod  = basic.Fault{TargetId: fileTarget.ID, Name: "chmod", NameCn: "篡改权限", Description: "file access permissions have been modified", DescriptionCn: "文件访问权限被修改"}
		fileFaultDelete = basic.Fault{TargetId: fileTarget.ID, Name: "del", NameCn: "删除文件", Description: "delete target file", DescriptionCn: "删除目标文件"}
		fileFaultAppend = basic.Fault{TargetId: fileTarget.ID, Name: "append", NameCn: "追加文件", Description: "append content to the target file, often used for exception log injection", DescriptionCn: "对目标文件追加内容，常用于异常日志注入"}
		fileFaultAdd    = basic.Fault{TargetId: fileTarget.ID, Name: "add", NameCn: "增加文件", Description: "add file", DescriptionCn: "增加文件"}
		fileFaultMv     = basic.Fault{TargetId: fileTarget.ID, Name: "mv", NameCn: "移动文件", Description: "move file", DescriptionCn: "移动文件"}
	)
	if err := basic.InsertFault(ctx, &fileFaultChmod); err != nil {
		return err
	}
	if err := InitFileTargetArgsChmod(ctx, fileFaultChmod); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &fileFaultDelete); err != nil {
		return err
	}
	if err := InitFileTargetArgsDelete(ctx, fileFaultDelete); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &fileFaultAppend); err != nil {
		return err
	}
	if err := InitFileTargetArgsAppend(ctx, fileFaultAppend); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &fileFaultAdd); err != nil {
		return err
	}

	if err := InitFileTargetArgsAdd(ctx, fileFaultAdd); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &fileFaultMv); err != nil {
		return err
	}

	return InitFileTargetArgsMove(ctx, fileFaultMv)
}

func InitFileTargetArgsChmod(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsPath       = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "路径", Unit: "file path", UnitCn: "文件路径", Description: "target file", DescriptionCn: "目标文件", ValueType: "string"}
		FileArgsPermission = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "permission", KeyCn: "权限", Unit: "3 integers in [0, 7], according to the Unix permission description specification", UnitCn: "3个在[0, 7]的整数，按照unix的权限描述规范", Description: "target permissions", DescriptionCn: "目标权限", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsPermission})
}

func InitFileTargetArgsDelete(ctx context.Context, fileFault basic.Fault) error {
	var FileArgsPath = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "路径", Unit: "file path", UnitCn: "文件路径", Description: "target file", DescriptionCn: "目标文件", ValueType: "string"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath})
}

func InitFileTargetArgsAppend(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsPath    = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "目标文件", UnitCn: "文件路径", Unit: "file path", Description: "target file", DescriptionCn: "目标文件", ValueType: "string"}
		FileArgsContent = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "content", KeyCn: "追加内容", Unit: "any string", UnitCn: "任意字符串", Description: "append content", DescriptionCn: "追加内容", ValueType: "string"}
		FileArgsRaw     = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "raw", KeyCn: "是否追加纯字符串", Unit: "true or false", UnitCn: "是、否", DefaultValue: "false", Description: "whether to append pure string. By default, it will add some additional identifiers to delete the appended content when recovering; if true, it will append pure string, and the appended content will not be deleted when recovering", DescriptionCn: "是否追加纯字符串。默认false会添加一些额外标识，用于恢复时删掉追加的内容；true追加纯字符串，则恢复时不删掉追加的内容", ValueType: "bool"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsContent, &FileArgsRaw})
}

func InitFileTargetArgsAdd(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsPath       = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "目标文件", Unit: "file path", UnitCn: "文件路径", Description: "target file", DescriptionCn: "目标文件", ValueType: "string"}
		FileArgsContent    = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "content", KeyCn: "文件内容", Unit: "any string", UnitCn: "任意字符串", Description: "file content", DescriptionCn: "文件内容", ValueType: "string"}
		FileArgsPermission = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "permission", KeyCn: "创建的文件权限", Unit: "3 integers within [0, 7], according to the unix permission description specification", UnitCn: "3个在[0, 7]的整数，按照unix的权限描述规范", DefaultValue: "system default permission", Description: "created file permission", DescriptionCn: "创建的文件权限", ValueType: "string"}
		FileArgsForce      = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "force", KeyCn: "是否覆盖已存在的文件", Unit: "true or false", UnitCn: "true、false", DefaultValue: "false", Description: "whether to overwrite existing files", DescriptionCn: "是否覆盖已存在的文件", ValueType: "bool"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsContent, &FileArgsPermission, &FileArgsForce})
}

func InitFileTargetArgsMove(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsSrc = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "src", KeyCn: "源文件路径", UnitCn: "文件路径", Unit: "file path", Description: "source file path", DescriptionCn: "源文件路径", ValueType: "string"}
		FileArgsDst = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "dst", KeyCn: "移动后文件路径", UnitCn: "文件路径", Unit: "file path", Description: "moved file path", DescriptionCn: "移动后文件路径", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsSrc, &FileArgsDst})
}

func InitKernelFault(ctx context.Context, kernelTarget basic.Target) error {
	var (
		kernelFaultFdfull = basic.Fault{TargetId: kernelTarget.ID, Name: "fdfull", NameCn: "系统fd耗尽", Description: "system fd exhausted, cannot use new fd (open file, create socket, create process), only affects non-root processes", DescriptionCn: "系统fd耗尽，无法使用新fd（打开文件、新建socket、新建进程），只对非root进程有影响"}
		kernelFaultNproc  = basic.Fault{TargetId: kernelTarget.ID, Name: "nproc", NameCn: "系统nproc满", Description: "target user process count is full, target user cannot create new processes", DescriptionCn: "目标用户进程数打满，目标用户无法创建新进程"}
	)
	if err := basic.InsertFault(ctx, &kernelFaultFdfull); err != nil {
		return err
	}
	if err := InitKernelTargetArgsFdfull(ctx, kernelFaultFdfull); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &kernelFaultNproc); err != nil {
		return err
	}
	return InitKernelTargetArgsNproc(ctx, kernelFaultNproc)
}

func InitKernelTargetArgsFdfull(ctx context.Context, KernelFault basic.Fault) error {
	var (
		KernelArgsCount = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "消耗量", Unit: "integer greater than or equal to 0, 0 means exhausted", UnitCn: "大于等于0的整数，0表示耗尽", DefaultValue: "0", Description: "fd consumption", DescriptionCn: "fd消耗量", ValueType: "int"}
		KernelArgsMode  = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "模式", Unit: "conf: by modifying the maximum fd number of the system, fill: by real consumption", UnitCn: "conf：通过修改系统最大fd数, fill：通过真实消耗", DefaultValue: "conf", Description: "execution mode", DescriptionCn: "执行模式", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&KernelArgsCount, &KernelArgsMode})
}

func InitKernelTargetArgsNproc(ctx context.Context, KernelFault basic.Fault) error {
	var (
		KernelArgsCount = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "消耗量", Unit: "integer greater than or equal to 0, 0 means exhausted", UnitCn: "大于等于0的整数，0表示耗尽", DefaultValue: "0", Description: "fd consumption", DescriptionCn: "fd消耗量", ValueType: "int"}
		KernelArgsUser  = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "user", KeyCn: "影响用户", Unit: "system user, does not support root user", UnitCn: "系统用户，不支持root用户", DefaultValue: "", Description: "affected user", DescriptionCn: "影响用户", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&KernelArgsCount, &KernelArgsUser})
}

func InitJvmFault(ctx context.Context, jvmTarget basic.Target) error {
	var (
		jvmFaultMethodDelay      = basic.Fault{TargetId: jvmTarget.ID, Name: "methoddelay", NameCn: "Java 运行时方法调用延迟", Description: "inject method call delay into running Java process", DescriptionCn: "对运行中的Java进程注入方法调用延迟"}
		jvmFaultMethodReturn     = basic.Fault{TargetId: jvmTarget.ID, Name: "methodreturn", NameCn: "Java 运行时方法返回值", Description: "Mock the return value of the specified method in the running Java process", DescriptionCn: "对运行中的Java进程Mock指定方法的返回值"}
		javaFaultMethodException = basic.Fault{TargetId: jvmTarget.ID, Name: "methodexception", NameCn: "Java 运行时方法抛出异常", Description: "make the specified method in the running Java process throw an exception when called", DescriptionCn: "使运行中的Java进程的指定方法被调用时抛出异常"}
	)
	if err := basic.InsertFault(ctx, &jvmFaultMethodDelay); err != nil {
		return err
	}
	if err := InitJvmTargetArgsMethod(ctx, jvmFaultMethodDelay); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &jvmFaultMethodReturn); err != nil {
		return err
	}
	if err := InitJvmTargetArgsMethod(ctx, jvmFaultMethodReturn); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &javaFaultMethodException); err != nil {
		return err
	}
	return InitJvmTargetArgsMethod(ctx, javaFaultMethodException)
}

func InitJvmTargetArgsMethod(ctx context.Context, javaFault basic.Fault) error {
	var (
		argsKey    = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "用来筛选受影响进程的关键词", Unit: "will use ps -ef | grep [key] to filter", UnitCn: "会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		argsPid    = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程pid", Unit: "pid of the running process", UnitCn: "存活进程的pid", ValueType: "string"}
		argsMethod = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "method", KeyCn: "目标方法以及延迟值，逗号分隔的列表，元素格式：类@方法@延迟毫秒", Unit: "target method and delay value, comma-separated list, element format: class@method@delay millisecond", UnitCn: "目标方法以及延迟值，逗号分隔的列表，元素格式：类@方法@延迟毫秒", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsKey, &argsPid, &argsMethod})
}

func InitContainerFault(ctx context.Context, containerTarget basic.Target) error {
	var (
		containerFaultKill    = basic.Fault{TargetId: containerTarget.ID, Name: "kill", NameCn: "杀容器", Description: "kill the target container", DescriptionCn: "把目标容器kill掉"}
		containerFaultPause   = basic.Fault{TargetId: containerTarget.ID, Name: "pause", NameCn: "暂停容器", Description: "pause the target container", DescriptionCn: "把目标容器暂停运行"}
		containerFaultRm      = basic.Fault{TargetId: containerTarget.ID, Name: "rm", NameCn: "删除容器", Description: "forcefully remove the target container", DescriptionCn: "强制删除容器"}
		containerFaultRestart = basic.Fault{TargetId: containerTarget.ID, Name: "restart", NameCn: "重启容器", Description: "restart the target container", DescriptionCn: "重启容器"}
	)
	if err := basic.InsertFault(ctx, &containerFaultKill); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultKill); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultPause); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultPause); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultRm); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultRm); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultRestart); err != nil {
		return err
	}
	return InitContainerArgs(ctx, containerFaultRestart)
}

func InitContainerArgs(ctx context.Context, containerFault basic.Fault) error {
	var (
		ContainerArgsId      = basic.Args{InjectId: containerFault.ID, ExecType: ExecInject, Key: "container-id", KeyCn: "容器ID", Unit: "id of the native container", UnitCn: "本机容器的ID", Description: "target container, do not specify the default attack local", DescriptionCn: "目标容器, 不指定默认攻击本地", ValueType: "string"}
		ContainerArgsRuntime = basic.Args{InjectId: containerFault.ID, ExecType: ExecInject, Key: "container-runtime", KeyCn: "目标容器运行时", Unit: "可选docker、pouch，如果指定了container-id，不指定runtime则默认为docker", UnitCn: "docker or pouch, if container-id is specified and runtime is not specified, it defaults to docker", ValueType: "string", DefaultValue: ""}
	)

	return basic.InsertArgsMulti(ctx, []*basic.Args{&ContainerArgsId, &ContainerArgsRuntime})
}
