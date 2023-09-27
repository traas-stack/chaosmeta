/*
 * Copyright 2022-2023 Chaos Meta Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
	ExecFlow   = "flow"
)

var (
	PodScope        = basic.Scope{Name: string(PodScopeType), NameCn: "Pod", Description: "Fault injection can be performed on any Pod in the Kubernetes cluster", DescriptionCn: "可以对Kubernetes集群中的任意Pod进行故障注入"}
	NodeScope       = basic.Scope{Name: string(NodeScopeType), NameCn: "Node", Description: "Fault injection can be performed on any Node in the Kubernetes cluster", DescriptionCn: "可以对Kubernetes集群中的任意Node进行故障注入"}
	KubernetesScope = basic.Scope{Name: string(KubernetesScopeType), NameCn: "Kubernetes", Description: "Faults can be injected into Kubernetes resource instances such as pod, deployment, and node to achieve the exception of the Kubernetes cluster itself or the exception of the operator application", DescriptionCn: "可以对Pod、Deployment、Node 等Kubernetes资源实例注入故障，达到Kubernetes集群自身的异常或者operator应用的异常"}
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
		CpuTarget       = basic.Target{Name: "cpu", NameCn: "cpu", Description: "Fault injection capabilities related to cpu faults", DescriptionCn: "cpu故障相关的故障注入能力"}
		MemTarget       = basic.Target{Name: "mem", NameCn: "mem", Description: "Fault injection capabilities related to memory faults", DescriptionCn: "内存故障相关的故障注入能力"}
		DiskTarget      = basic.Target{Name: "disk", NameCn: "disk", Description: "Fault injection capabilities related to disk failures", DescriptionCn: "磁盘故障相关的故障注入能力"}
		DiskioTarget    = basic.Target{Name: "diskIO", NameCn: "diskIO", Description: "Fault injection capabilities related to disk IO faults", DescriptionCn: "磁盘IO故障相关的故障注入能力"}
		NetworkTarget   = basic.Target{Name: "network", NameCn: "network", Description: "Fault injection capabilities related to disk failures", DescriptionCn: "磁盘故障相关的故障注入能力"}
		ProcessTarget   = basic.Target{Name: "process", NameCn: "process", Description: "Process-dependent fault injection capability", DescriptionCn: "进程相关的故障注入能力"}
		FileTarget      = basic.Target{Name: "file", NameCn: "file", Description: "File-related fault injection capabilities", DescriptionCn: "文件相关的故障注入能力"}
		KernelTarget    = basic.Target{Name: "kernel", NameCn: "kernel", Description: "Kernel-related fault injection capabilities", DescriptionCn: "内核相关的故障注入能力"}
		JvmTarget       = basic.Target{Name: "jvm", NameCn: "jvm", Description: "Jvm related fault injection capabilities", DescriptionCn: "围绕jvm相关的故障注入能力"}
		ContainerTarget = basic.Target{Name: "container", NameCn: "container", Description: "Container runtime-related fault injection capabilities", DescriptionCn: "容器运行时相关的故障注入能力"}
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
	if err := InitMemFault(ctx, MemTarget); err != nil {
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
		CpuFaultBurn = basic.Fault{TargetId: cpuTarget.ID, Name: "burn", NameCn: "cpu使用率", Description: "The CPU usage rate soars,provide at least one of the count and list parameters,when both are provided, count will prevail and list will be ignored.", DescriptionCn: "cpu使用率飙高,count和list参数至少提供一个,都提供的时候，以count为准,忽略list"}
		CpuFaultLoad = basic.Fault{TargetId: cpuTarget.ID, Name: "load", NameCn: "cpu负载", Description: "High cpu usage", DescriptionCn: "cpu平均负载飙高"}
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
		CpuArgsPercent = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "使用率", Unit: "%", UnitCn: "%", Description: "Target cpu usage", DescriptionCn: "目标cpu使用率", ValueType: "int", Required: true, ValueRule: "1-100"}
		CpuArgsCount   = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "核数", Unit: "", UnitCn: "", DefaultValue: "0", Description: "Number of faulty CPU cores, 0 means all cores", DescriptionCn: "故障cpu核数,0表示全部核", ValueType: "int", ValueRule: ">=0"}
		CpuArgsList    = basic.Args{InjectId: cpuFault.ID, ExecType: ExecInject, Key: "list", KeyCn: "列表", Unit: "", UnitCn: "", Description: "Faulty cpu list, comma separated core number list, can be confirmed from /proc/cpuinfo", DescriptionCn: "故障cpu列表,逗号分隔的核编号列表,可以从/proc/cpuinfo确认", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&CpuArgsPercent, &CpuArgsCount, &CpuArgsList})
}

func InitLoadTargetArgsLoad(ctx context.Context, loadFault basic.Fault) error {
	var LoadArgsCount = basic.Args{InjectId: loadFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "负载数", Unit: "", UnitCn: "", Description: "Number of loads to add", DescriptionCn: "需要增加的负载数, 如果为0表示cpu核数的4倍", Required: true, ValueType: "int", ValueRule: ">=0"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&LoadArgsCount})
}

func InitMemFault(ctx context.Context, memTarget basic.Target) error {
	var (
		MemFaultFill = basic.Fault{TargetId: memTarget.ID, Name: "fill", NameCn: "内存填充", Description: "Memory usage spikes, when percentage and bytes arguments are both provided, same as percentage, bytes ignored", DescriptionCn: "内存使用率飙高,percent和bytes参数都提供的时候,以percent为准,忽略bytes"}
		MemFaultOom  = basic.Fault{TargetId: memTarget.ID, Name: "oom", NameCn: "内存oom", Description: "The system memory oom will cause the machine to hang up", DescriptionCn: "系统内存oom,会使机器宕机挂掉"}
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
		MemArgsPercent = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "内存使用率", Unit: "%", UnitCn: "%", Description: "Target mem usage", DescriptionCn: "目标内存使用率", ValueType: "int", Required: true, ValueRule: "1-100"}
		MemArgsBytes   = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "bytes", KeyCn: "填充量", Unit: "KB,MB,GB,TB", UnitCn: "KB,MB,GB,TB", Description: "Memory fill", DescriptionCn: "内存填充量", ValueType: "string"}
		MemArgsMode    = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "填充模式", Unit: "mode", UnitCn: "模式", Description: "Memory filling mode, ram is the way to apply for process memory, cache is the way to use tmpfs", DescriptionCn: "内存填充模式,ram是使用进程内存申请的方式,cache是使用tmpfs的方式", ValueType: "string", Required: true, ValueRule: "ram,cache"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&MemArgsPercent, &MemArgsBytes, &MemArgsMode})
}

func InitMemTargetArgsOom(ctx context.Context, memFault basic.Fault) error {
	var MemArgsMode = basic.Args{InjectId: memFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "内存填充模式", Unit: "mode", UnitCn: "模式", DefaultValue: "cache", Description: "Memory filling mode, ram is the way to apply for process memory, cache is the way to use tmpfs", DescriptionCn: "内存填充模式,ram是使用进程内存申请的方式,cache是使用tmpfs的方式", ValueType: "string", Required: true, ValueRule: "ram,cache"}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&MemArgsMode})
}

func InitDiskFault(ctx context.Context, diskTarget basic.Target) error {
	var DiskFaultFill = basic.Fault{TargetId: diskTarget.ID, Name: "fill", NameCn: "磁盘填充", Description: "The disk usage is so high, when both the percent and bytes parameters are provided, the percent will prevail and bytes will be ignored", DescriptionCn: "磁盘使用率飙高,percent和bytes参数都提供的时候,以percent为准,忽略bytes"}
	if err := basic.InsertFault(ctx, &DiskFaultFill); err != nil {
		return err
	}
	return InitDiskTargetArgsFill(ctx, DiskFaultFill)
}

func InitDiskTargetArgsFill(ctx context.Context, diskFault basic.Fault) error {
	var (
		DiskArgsPercent = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "磁盘使用率", Unit: "%", UnitCn: "%", Description: "Target disk usage", DescriptionCn: "目标磁盘使用率", ValueType: "int", Required: true, ValueRule: "1-100"}
		DiskArgsBytes   = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "bytes", KeyCn: "填充量", Unit: "KB,MB,GB,TB", UnitCn: "KB,MB,GB,TB", Description: "Memory fill", DescriptionCn: "磁盘填充量", ValueType: "string"}
		DiskArgsDir     = basic.Args{InjectId: diskFault.ID, ExecType: ExecInject, Key: "dir", KeyCn: "目录", Unit: "", UnitCn: "", DefaultValue: "/tmp", Description: "Target population directory", DescriptionCn: "目标填充目录", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskArgsPercent, &DiskArgsBytes, &DiskArgsDir})
}

func InitDiskioFault(ctx context.Context, diskioTarget basic.Target) error {
	var (
		DiskioFaultBurn  = basic.Fault{TargetId: diskioTarget.ID, Name: "burn", NameCn: "磁盘IO高负载", Description: "Disk IO soars", DescriptionCn: "磁盘IO飙高"}
		DiskioFaultHang  = basic.Fault{TargetId: diskioTarget.ID, Name: "hang", NameCn: "磁盘IO hang", Description: "The target process generates a disk IO hang; provide at least one of the pid-list and key parameters. When both are provided, the pid-list will prevail and the key will be ignored; the principle of this injection capability is to limit the process to only 1byte of IO per second.Therefore, it has little impact on processes with too small IO", DescriptionCn: "目标进程产生磁盘IO hang;pid-list和key参数至少提供一个,都提供的时候,以pid-list为准,忽略key;此注入能力的原理是限制进程每秒只能进行1byte大小的IO,所以对IO过小的进程影响不大"}
		DiskioFaultLimit = basic.Fault{TargetId: diskioTarget.ID, Name: "limit", NameCn: "磁盘IO limit", Description: "The target process generates a disk IO limit; when both pid-list and key parameters are provided, the pid-list shall prevail and the key shall be ignored; at least one of the four limit parameters must be provided, and multiple limits are in an \"AND\" relationship", DescriptionCn: "目标进程产生磁盘IO limit;pid-list和key参数都提供的时候,以pid-list为准,忽略key;四个限制参数至少提供一个,多个限制是“与”的关系"}
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
		DiskioArgsDir   = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dir", KeyCn: "目录", Unit: "", UnitCn: "", DefaultValue: "/tmp", Description: "Target directory for high IO operations", DescriptionCn: "进行高IO操作的目标目录", ValueType: "string", Required: true}
		DiskioArgsMode  = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "IO模式", Unit: "", UnitCn: "", DefaultValue: "read", Description: "IO mode", DescriptionCn: "IO模式", ValueType: "string", Required: true, ValueRule: "read,write"}
		DiskioArgsBlock = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "block", KeyCn: "IO块", Unit: "KB,MB", UnitCn: "KB,MB", DefaultValue: "10MB", Description: "The block size of a single IO, ranging from 1K-1024M", DescriptionCn: "单次IO的块大小,范围为1K-1024M", ValueType: "string", Required: true, ValueRule: "1-1024"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDir, &DiskioArgsMode, &DiskioArgsBlock})
}

func InitDiskioTargetArgsHang(ctx context.Context, diskioFault basic.Fault) error {
	var (
		DiskioArgsDevList = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dev-list", KeyCn: "设备列表", Unit: "", UnitCn: "", Description: "Target disk device list, use the command lsblk -a | grep disk to obtain the primary and secondary device numbers of the target device, such as 8:0", DescriptionCn: "目标磁盘设备列表,使用命令lsblk -a | grep disk获取目标设备的主备设备号,比如8:0", ValueType: "stringlist"}
		DiskioArgsPidList = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "pid-list", KeyCn: "进程pid列表", Unit: "", UnitCn: "", Description: "Affected process pid list, for example: 7850, 7690", DescriptionCn: "受影响的进程pid列表,比如:7850,7690", ValueType: "stringlist"}
		DiskioArgsKey     = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Unit: "", UnitCn: "", Description: "keywords used to filter affected processes will be filtered using ps -ef | grep [key]", DescriptionCn: "用来筛选受影响进程的关键词,会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		DiskioArgsMode    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "IO模式", Unit: "mode", UnitCn: "模式", DefaultValue: "all", Description: "Affected IO operation", DescriptionCn: "受影响的IO操作", ValueType: "string", ValueRule: "all,read,write"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDevList, &DiskioArgsPidList, &DiskioArgsKey, &DiskioArgsMode})
}

func InitDiskioTargetArgsLimit(ctx context.Context, diskioFault basic.Fault) error {
	var (
		DiskioArgsDevList    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "dev-list", KeyCn: "设备列表", Unit: "", UnitCn: "", Description: "Target disk device list, use the command lsblk -a | grep disk to obtain the primary and secondary device numbers of the target device, such as 8:0", DescriptionCn: "目标磁盘设备列表,使用命令lsblk -a | grep disk获取目标设备的主备设备号,比如8:0", ValueType: "stringlist"}
		DiskioArgsPidList    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "pid-list", KeyCn: "进程pid列表", Unit: "", UnitCn: "", Description: "Affected process pid list, such as 7850, 7690", DescriptionCn: "受影响的进程pid列表,比如7850,7690", ValueType: "stringlist"}
		DiskioArgsKey        = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Unit: "", UnitCn: "", Description: "Keywords used to filter affected processes will be filtered using ps -ef | grep [key]", DescriptionCn: "用来筛选受影响进程的关键词,会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		DiskioArgsReadBytes  = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "read-bytes", KeyCn: "读字节数", Unit: "B,KB,MB,GB,TB", UnitCn: "B,KB,MB,GB,TB", Description: "Number of bytes that can be read per second", DescriptionCn: "每秒能读的字节数", ValueType: "string"}
		DiskioArgsWriteBytes = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "write-bytes", KeyCn: "写字节数", Unit: "B,KB,MB,GB,TB", UnitCn: "B,KB,MB,GB,TB", Description: "Number of bytes that can be written per second", DescriptionCn: "每秒能写的字节数", ValueType: "string"}
		DiskioArgsReadIO     = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "read-io", KeyCn: "读IO次数", Unit: "", UnitCn: "", Description: "Number of IO operations that can be read per second", DescriptionCn: "每秒能读的IO次数", ValueType: "int", ValueRule: ">0"}
		DiskioArgsWriteIO    = basic.Args{InjectId: diskioFault.ID, ExecType: ExecInject, Key: "write-io", KeyCn: "写IO次数", Unit: "", UnitCn: "", Description: "Number of IO operations that can be written per second", DescriptionCn: "每秒能写的IO次数", ValueType: "int", ValueRule: ">0"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&DiskioArgsDevList, &DiskioArgsPidList, &DiskioArgsKey, &DiskioArgsReadBytes, &DiskioArgsWriteBytes, &DiskioArgsReadIO, &DiskioArgsWriteIO})
}

func InitNetworkFault(ctx context.Context, networkTarget basic.Target) error {
	var (
		NetworkFaultOccupy  = basic.Fault{TargetId: networkTarget.ID, Name: "occupy", NameCn: "端口占用", Description: "Kill the service on the target port and occupy this port", DescriptionCn: "把目标端口的服务kill掉并把这个端口占用掉"}
		NetworkFaultLimit   = basic.Fault{TargetId: networkTarget.ID, Name: "limit", NameCn: "带宽限制", Description: "To inject bandwidth restrictions on network packets in the outflow direction from the faulty machine, the interface must be provided. The packet filtering parameters can be optionally provided. The relationship between \"and\"", DescriptionCn: "对从故障机器流出方向的网络数据包注入带宽限制,interface必须提供,数据包筛选参数可以选择性提供,“与”的关系"}
		NetworkFaultDelay   = basic.Fault{TargetId: networkTarget.ID, Name: "delay", NameCn: "网络延迟", Description: "Network packet injection delay in the outflow direction from the faulty machine; interface must be provided, packet filtering parameters can be optionally provided, the relationship between \"and\"", DescriptionCn: "从故障机器流出方向的网络数据包注入延迟;interface必须提供,数据包筛选参数可以选择性提供,“与”的关系"}
		NetworkFaultLoss    = basic.Fault{TargetId: networkTarget.ID, Name: "loss", NameCn: "网络丢包", Description: "Network data packets in the outflow direction from the faulty machine are injected into packet loss; interface must be provided, and packet filtering parameters can be optionally provided, the relationship between \"and\"", DescriptionCn: "从故障机器流出方向的网络数据包注入丢包;interface 必须提供,数据包筛选参数可以选择性提供,“与”的关系"}
		NetworkFaultCorrupt = basic.Fault{TargetId: networkTarget.ID, Name: "corrupt", NameCn: "网络包损坏", Description: "Inject packet damage into the network data packets flowing out from the faulty machine; the interface must be provided, and the packet filtering parameters can be optionally provided, the relationship between \"and\"", DescriptionCn: "对从故障机器流出方向的网络数据包注入包损坏;interface必须提供,数据包筛选参数可以选择性提供,“与”的关系"}

		NetworkFaultDuplicate = basic.Fault{TargetId: networkTarget.ID, Name: "duplicate", NameCn: "网络包重复", Description: "Repeat for network packet injection packets in the outflow direction from the faulty machine; interface must be provided, packet filtering parameters can be optionally provided, \"relationship with\"", DescriptionCn: "对从故障机器流出方向的网络数据包注入包重复;interface必须提供，数据包筛选参数可以选择性提供,“与”的关系"}
		NetworkFaultReorder   = basic.Fault{TargetId: networkTarget.ID, Name: "reorder", NameCn: "网络包乱序", Description: "Inject packet reordering into network data packets flowing out from the faulty machine; the interface must be provided, and the packet filtering parameters can be optionally provided, the relationship between \"and\"", DescriptionCn: "对从故障机器流出方向的网络数据包注入包乱序;interface必须提供，数据包筛选参数可以选择性提供,“与”的关系"}
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
	if err := InitNetworkTargetArgsLoss(ctx, NetworkFaultLoss); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &NetworkFaultCorrupt); err != nil {
		return err
	}
	if err := InitNetworkTargetArgsCorrupt(ctx, NetworkFaultCorrupt); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &NetworkFaultDuplicate); err != nil {
		return err
	}
	if err := InitNetworkTargetArgsDuplicate(ctx, NetworkFaultDuplicate); err != nil {
		return err
	}
	if err := basic.InsertFault(ctx, &NetworkFaultReorder); err != nil {
		return err
	}
	return InitNetworkTargetArgsReorder(ctx, NetworkFaultReorder)
}

func getNetworkCommonFilterParameters(fault basic.Fault) []*basic.Args {
	var (
		NetworkArgsInterface = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "interface", KeyCn: "网卡", Description: "The network card included in the faulty machine, such as eth0", DescriptionCn: "故障机器包含的网卡,比如eth0", Required: true, ValueType: "string"}
		NetworkArgsDstIP     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "dst-ip", KeyCn: "目标ip列表", Description: "Target IP list, such as: 1.2.3.4, 2.3.4.5, 192.168.1.1/24", DescriptionCn: "目标ip列表,比如1.2.3.4,2.3.4.5,192.168.1.1/24", ValueType: "stringlist"}
		NetworkArgsSrcIP     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "src-ip", KeyCn: "数据包筛选参数：源ip列表", Description: "Source IP list, for example: 1.2.3.4,2.3.4.5,192.168.1.1/24", DescriptionCn: "源ip列表，比如1.2.3.4,2.3.4.5,192.168.1.1/24", ValueType: "stringlist"}
		NetworkArgsDstPort   = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "dst-port", KeyCn: "数据包筛选参数：目标端口列表", Description: "Target port list, for example: 8080-8090,8095,9099", DescriptionCn: "目标端口列表，比如8080-8090,8095,9099", ValueType: "stringlist"}
		NetworkArgsSrcPort   = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "src-port", KeyCn: "数据包筛选参数：源端口列表", Description: "Source port list, such as 8080-8090,8095,9099", DescriptionCn: "源端口列表，比如8080-8090,8095,9099", ValueType: "stringlist"}
		NetworkArgsMode      = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "数据包筛选模式", Unit: "", UnitCn: "", DefaultValue: "normal", Description: "normal: inject fault to selected targets, exclude: do not inject fault to selected targets", DescriptionCn: "normal:对选中的目标注入故障,exclude:对选中的目标不注入故障", ValueType: "string", ValueRule: "normal,exclude"}
		NetworkArgsForce     = basic.Args{InjectId: fault.ID, ExecType: ExecInject, Key: "force", KeyCn: "是否强制覆盖", DefaultValue: "false", Description: "Whether to force overwrite", DescriptionCn: "是否强制覆盖", ValueType: "bool", ValueRule: "true,false"}
	)
	return []*basic.Args{&NetworkArgsInterface, &NetworkArgsDstIP, &NetworkArgsSrcIP, &NetworkArgsDstPort, &NetworkArgsSrcPort, &NetworkArgsMode, &NetworkArgsForce}
}

func InitNetworkTargetArgsOccupy(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsPort       = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "port", KeyCn: "端口", Description: "target port", DescriptionCn: "目标端口", ValueType: "int"}
		NetworkArgsProtocol   = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "protocol", KeyCn: "协议", DefaultValue: "tcp", Description: "target protocol", DescriptionCn: "目标协议", ValueType: "string", ValueRule: "tcp,udp,tcp6,udp6"}
		NetworkArgsRecoverCmd = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "recover-cmd", KeyCn: "恢复命令", Description: "resume command, it will be executed last when resuming operation", DescriptionCn: "恢复命令，恢复操作时会最后执行", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&NetworkArgsPort, &NetworkArgsProtocol, &NetworkArgsRecoverCmd})
}

func InitNetworkTargetArgsLimit(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsRate = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "rate", KeyCn: "每秒的网络带宽限制", Unit: "bit,kbit,mbit,gbit,tbit", UnitCn: "bit,kbit,mbit,gbit,tbit", Description: "Network bandwidth limit per second", DescriptionCn: "每秒的网络带宽限制", ValueType: "int", Required: true}
	)
	argList := []*basic.Args{&NetworkArgsRate}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsDelay(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsLatency = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "latency", KeyCn: "延迟时间", Unit: "us,ms,s", UnitCn: "us,ms,s", Description: "Delay time", DescriptionCn: "延迟时间", ValueType: "int", Required: true}
		NetworkArgsJitter  = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "jitter", KeyCn: "抖动值", Unit: "us,ms,s", UnitCn: "us,ms,s", Description: "Jitter value, the fluctuation range of each delay", DescriptionCn: "抖动值,每次延迟的波动范围", DefaultValue: "0", ValueType: "int", Required: true}
	)

	argList := []*basic.Args{&NetworkArgsLatency, &NetworkArgsJitter}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsLoss(ctx context.Context, networkFault basic.Fault) error {
	var NetworkArgsPercent = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "丢包率", Description: "Packet loss rate", DescriptionCn: "丢包率", ValueType: "int", Required: true, ValueRule: "1-100"}

	argList := []*basic.Args{&NetworkArgsPercent}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsCorrupt(ctx context.Context, networkFault basic.Fault) error {
	var NetworkArgsPercent = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "包损坏率", Description: "Packet loss rate", DescriptionCn: "包损坏率", ValueType: "int", Required: true, ValueRule: "1-100"}

	argList := []*basic.Args{&NetworkArgsPercent}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsDuplicate(ctx context.Context, networkFault basic.Fault) error {
	var NetworkArgsPercent = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "percent", KeyCn: "包重复率", Description: "packet loss rate", DescriptionCn: "丢包率", ValueType: "int", Required: true, ValueRule: "1-100"}

	argList := []*basic.Args{&NetworkArgsPercent}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitNetworkTargetArgsReorder(ctx context.Context, networkFault basic.Fault) error {
	var (
		NetworkArgsLatency = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "latency", KeyCn: "包延迟", Unit: "us,ms,s", UnitCn: "us,ms,s", DefaultValue: "1s", Description: "delay", DescriptionCn: "延迟时间", ValueType: "int", Required: true}
		NetworkArgsGap     = basic.Args{InjectId: networkFault.ID, ExecType: ExecInject, Key: "gap", KeyCn: "KeyCn", DefaultValue: "3", Description: "Select the interval. For example, a gap of 3 means that packets with serial numbers 1, 3, 6, 9, etc. will not be delayed, and the remaining packets will be delayed", DescriptionCn: "选中间隔,比如gap为3表示序号为1、3、6、9等的包不延迟,其余的包会延迟", ValueType: "int", ValueRule: ">0"}
	)
	argList := []*basic.Args{&NetworkArgsLatency, &NetworkArgsGap}
	argList = append(argList, getNetworkCommonFilterParameters(networkFault)...)
	return basic.InsertArgsMulti(ctx, argList)
}

func InitProcessFault(ctx context.Context, processTarget basic.Target) error {
	var (
		processFaultKill = basic.Fault{TargetId: processTarget.ID, Name: "kill", NameCn: "杀进程", Description: "To kill the target process, provide at least one of the pid and key parameters. When both are provided, the pid will prevail and the key will be ignored.", DescriptionCn: "把目标进程杀掉,pid和key参数至少提供一个,都提供的时候,以pid为准,忽略key"}
		processFaultStop = basic.Fault{TargetId: processTarget.ID, Name: "stop", NameCn: "停止进程", Description: "Stop the target process. Provide at least one of the pid and key parameters. When both are provided, the pid will prevail and the key will be ignored.", DescriptionCn: "停止目标进程,pid和key参数至少提供一个,都提供的时候,以pid为准,忽略key"}
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
		ProcessArgsKey        = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Description: "keywords used to filter affected processes; Will use ps -ef | grep [key] to filter", DescriptionCn: "用来筛选受影响进程的关键词;会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		ProcessArgsPid        = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程pid", Description: "the pid of the living process", DescriptionCn: "存活进程的pid", ValueType: "int"}
		ProcessArgsSignal     = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "signal", KeyCn: "信号", DefaultValue: "9", Description: "the signal sent to the process;consistent with the signal value supported by the kill command", DescriptionCn: "对进程发送的信号;和kill命令支持的信号数值一致", ValueType: "int"}
		ProcessArgsRecoverCmd = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "recover-cmd", KeyCn: "恢复命令", Description: "resume command, it will be executed last when resuming operation", DescriptionCn: "恢复命令，恢复操作时会最后执行", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&ProcessArgsKey, &ProcessArgsPid, &ProcessArgsSignal, &ProcessArgsRecoverCmd})
}

func InitProcessTargetArgsStop(ctx context.Context, processFault basic.Fault) error {
	var (
		ProcessArgsKey = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "关键词", Description: "Keywords used to filter affected processes, will use ps -ef | grep [key] to filter", DescriptionCn: "用来筛选受影响进程的关键词;会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		ProcessArgsPid = basic.Args{InjectId: processFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程", Description: "The pid of the living process", DescriptionCn: "存活进程的pid", ValueType: "int"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&ProcessArgsKey, &ProcessArgsPid})
}

func InitFileFault(ctx context.Context, fileTarget basic.Target) error {
	var (
		fileFaultChmod  = basic.Fault{TargetId: fileTarget.ID, Name: "chmod", NameCn: "篡改权限", Description: "File access permissions have been modified", DescriptionCn: "文件访问权限被修改"}
		fileFaultDelete = basic.Fault{TargetId: fileTarget.ID, Name: "del", NameCn: "删除文件", Description: "Delete target file", DescriptionCn: "删除目标文件"}
		fileFaultAppend = basic.Fault{TargetId: fileTarget.ID, Name: "append", NameCn: "追加文件", Description: "Append content to the target file, often used for exception log injection", DescriptionCn: "对目标文件追加内容，常用于异常日志注入"}
		fileFaultAdd    = basic.Fault{TargetId: fileTarget.ID, Name: "add", NameCn: "增加文件", Description: "Add file", DescriptionCn: "增加文件"}
		fileFaultMv     = basic.Fault{TargetId: fileTarget.ID, Name: "mv", NameCn: "移动文件", Description: "Move file", DescriptionCn: "移动文件"}
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
		FileArgsPath       = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "路径", Description: "Target file path", DescriptionCn: "目标文件路径", ValueType: "string", Required: true}
		FileArgsPermission = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "permission", KeyCn: "权限", Description: "Target permissions; 3 integers in [0, 7], according to the Unix permission description specification", DescriptionCn: "目标权限,3个在[0, 7]的整数,按照unix的权限描述规范", ValueType: "string", Required: true}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsPermission})
}

func InitFileTargetArgsDelete(ctx context.Context, fileFault basic.Fault) error {
	var FileArgsPath = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "路径", Description: "Target file path", DescriptionCn: "目标文件路径", ValueType: "string", Required: true}
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath})
}

func InitFileTargetArgsAppend(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsPath    = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "目标文件", Description: "Target file path", DescriptionCn: "目标文件路径", ValueType: "string", Required: true}
		FileArgsContent = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "content", KeyCn: "追加内容", Description: "Append content", DescriptionCn: "追加内容", ValueType: "string"}
		FileArgsRaw     = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "raw", KeyCn: "是否追加纯字符串", DefaultValue: "false", Description: "Whether to append pure string. By default, it will add some additional identifiers to delete the appended content when recovering; if true, it will append pure string, and the appended content will not be deleted when recovering", DescriptionCn: "是否追加纯字符串,默认false会添加一些额外标识,用于恢复时删掉追加的内容;true追加纯字符串，则恢复时不删掉追加的内容", ValueType: "bool", ValueRule: "true,false"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsContent, &FileArgsRaw})
}

func InitFileTargetArgsAdd(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsPath       = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "path", KeyCn: "目标文件", Description: "Target file path", DescriptionCn: "目标文件路径", ValueType: "string", Required: true}
		FileArgsContent    = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "content", KeyCn: "文件内容", Description: "File content", DescriptionCn: "文件内容", ValueType: "string"}
		FileArgsPermission = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "permission", KeyCn: "创建的文件权限", Unit: "", UnitCn: "", DefaultValue: "", Description: "Created file permission; integers within [0, 7], according to the unix permission description specification", DescriptionCn: "创建的文件权限; 3个在[0, 7]的整数，按照unix的权限描述规范", ValueType: "string"}
		FileArgsForce      = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "force", KeyCn: "是否覆盖已存在的文件", DefaultValue: "false", Description: "Whether to overwrite existing files", DescriptionCn: "是否覆盖已存在的文件", ValueType: "bool", ValueRule: "true,false"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsPath, &FileArgsContent, &FileArgsPermission, &FileArgsForce})
}

func InitFileTargetArgsMove(ctx context.Context, fileFault basic.Fault) error {
	var (
		FileArgsSrc = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "src", KeyCn: "源文件路径", Description: "source file path", DescriptionCn: "源文件路径", ValueType: "string", Required: true}
		FileArgsDst = basic.Args{InjectId: fileFault.ID, ExecType: ExecInject, Key: "dst", KeyCn: "移动后文件路径", Description: "moved file path", DescriptionCn: "移动后文件路径", ValueType: "string", Required: true}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&FileArgsSrc, &FileArgsDst})
}

func InitKernelFault(ctx context.Context, kernelTarget basic.Target) error {
	var (
		kernelFaultFdfull = basic.Fault{TargetId: kernelTarget.ID, Name: "fdfull", NameCn: "系统fd耗尽", Description: "System fd exhausted, cannot use new fd (open file, create socket, create process), only affects non-root processes", DescriptionCn: "系统fd耗尽,无法使用新fd(打开文件、新建socket、新建进程),只对非root进程有影响;fill模式下,如果系统最大fd数过大可能会先导致oom"}
		kernelFaultNproc  = basic.Fault{TargetId: kernelTarget.ID, Name: "nproc", NameCn: "系统nproc满", Description: "Target user process count is full, target user cannot create new processes", DescriptionCn: "目标用户进程数打满,目标用户无法创建新进程"}
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
		KernelArgsCount = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "消耗量", DefaultValue: "0", Description: "fd consumption: an integer greater than or equal to 0, 0 means exhausted", DescriptionCn: "fd消耗量:大于等于0的整数,0表示耗尽", ValueType: "int", ValueRule: ">=0"}
		KernelArgsMode  = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "mode", KeyCn: "执行模式", DefaultValue: "conf", Description: "Execution mode, conf: by modifying the maximum number of FDs in the system, fill: by real consumption", DescriptionCn: "执行模式,conf:通过修改系统最大fd数,fill:通过真实消耗", ValueType: "string", ValueRule: "conf,fill"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&KernelArgsCount, &KernelArgsMode})
}

func InitKernelTargetArgsNproc(ctx context.Context, KernelFault basic.Fault) error {
	var (
		KernelArgsCount = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "count", KeyCn: "消耗进程数", DefaultValue: "0", Description: "The number of consumed processes, an integer greater than or equal to 0, 0 means exhausted", DescriptionCn: "消耗进程数,大于等于0的整数,0表示耗尽", ValueType: "int", ValueRule: ">=0"}
		KernelArgsUser  = basic.Args{InjectId: KernelFault.ID, ExecType: ExecInject, Key: "user", KeyCn: "影响用户", Description: "Affected users, root users are not supported", DescriptionCn: "影响用户,不支持root用户", ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&KernelArgsCount, &KernelArgsUser})
}

func InitJvmFault(ctx context.Context, jvmTarget basic.Target) error {
	var (
		jvmFaultMethodDelay      = basic.Fault{TargetId: jvmTarget.ID, Name: "methoddelay", NameCn: "Java运行时方法调用延迟", Description: "inject method call delay into running Java process", DescriptionCn: "对运行中的Java进程注入方法调用延迟"}
		jvmFaultMethodReturn     = basic.Fault{TargetId: jvmTarget.ID, Name: "methodreturn", NameCn: "Java运行时篡改返回值", Description: "Mock the return value of the specified method in the running Java process", DescriptionCn: "对运行中的Java进程Mock指定方法的返回值"}
		javaFaultMethodException = basic.Fault{TargetId: jvmTarget.ID, Name: "methodexception", NameCn: "Java运行时方法抛出异常", Description: "make the specified method in the running Java process throw an exception when called", DescriptionCn: "使运行中的Java进程的指定方法被调用时抛出异常"}
	)
	if err := basic.InsertFault(ctx, &jvmFaultMethodDelay); err != nil {
		return err
	}
	if err := InitJvmTargetArgsMethod(ctx, jvmFaultMethodDelay, "目标方法以及延迟值", "Comma separated list, element format: class@method@delay milliseconds", "逗号分隔的列表,元素格式:类@方法@延迟毫秒"); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &jvmFaultMethodReturn); err != nil {
		return err
	}
	if err := InitJvmTargetArgsMethod(ctx, jvmFaultMethodReturn, "目标方法以及返回值", "Comma separated list, element format: class@method@return value, return integer: Client@say@10, return string: Client@say@\"test\", return variable: Client@say@var", "逗号分隔的列表,元素格式：类@方法@返回值,返回整数:Client@say@10,返回字符串:Client@say@\"test\",返回变量:Client@say@var"); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &javaFaultMethodException); err != nil {
		return err
	}
	return InitJvmTargetArgsMethod(ctx, javaFaultMethodException, "目标方法以及异常信息", "Comma separated list, element format: class@method@exception description information, example: Client@say@test\"", "逗号分隔的列表,元素格式:类@方法@异常描述信息,样例:Client@say@test")
}

func InitJvmTargetArgsMethod(ctx context.Context, javaFault basic.Fault, argsMethodKeyCn string, argsMethodDescription string, argsMethodDescriptionCn string) error {
	var (
		argsKey    = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "key", KeyCn: "进程关键词", Description: "will use ps -ef | grep [key] to filter", DescriptionCn: "用来筛选受影响进程的关键词;会使用ps -ef | grep [key]来筛选", ValueType: "string"}
		argsPid    = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "pid", KeyCn: "进程pid", Description: "pid of the running process", DescriptionCn: "存活进程的pid", ValueType: "string"}
		argsMethod = basic.Args{InjectId: javaFault.ID, ExecType: ExecInject, Key: "method", KeyCn: argsMethodKeyCn, Description: argsMethodDescription, DescriptionCn: argsMethodDescriptionCn, ValueType: "string"}
	)
	return basic.InsertArgsMulti(ctx, []*basic.Args{&argsKey, &argsPid, &argsMethod})
}

func InitContainerFault(ctx context.Context, containerTarget basic.Target) error {
	var (
		containerFaultKill    = basic.Fault{TargetId: containerTarget.ID, Name: "kill", NameCn: "杀容器", Description: "Kill the target container", DescriptionCn: "把目标容器kill掉"}
		containerFaultPause   = basic.Fault{TargetId: containerTarget.ID, Name: "pause", NameCn: "暂停容器", Description: "Pause the target container", DescriptionCn: "把目标容器暂停运行"}
		containerFaultRm      = basic.Fault{TargetId: containerTarget.ID, Name: "rm", NameCn: "删除容器", Description: "Forcefully remove the target container", DescriptionCn: "强制删除容器"}
		containerFaultRestart = basic.Fault{TargetId: containerTarget.ID, Name: "restart", NameCn: "重启容器", Description: "Restart the target container", DescriptionCn: "重启容器"}
	)
	if err := basic.InsertFault(ctx, &containerFaultKill); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultKill, false); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultPause); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultPause, false); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultRm); err != nil {
		return err
	}
	if err := InitContainerArgs(ctx, containerFaultRm, false); err != nil {
		return err
	}

	if err := basic.InsertFault(ctx, &containerFaultRestart); err != nil {
		return err
	}
	return InitContainerArgs(ctx, containerFaultRestart, true)
}

func InitContainerArgs(ctx context.Context, containerFault basic.Fault, isHaveWaitTimeArg bool) error {
	var (
		ContainerArgsId       = basic.Args{InjectId: containerFault.ID, ExecType: ExecInject, Key: "container-id", KeyCn: "目标容器ID", Description: "Target container, do not specify the default attack local", DescriptionCn: "目标容器,不指定默认攻击本地", ValueType: "int"}
		ContainerArgsRuntime  = basic.Args{InjectId: containerFault.ID, ExecType: ExecInject, Key: "container-runtime", KeyCn: "目标容器runtime", Description: "Docker or pouch, if container-id is specified and runtime is not specified, it defaults to docker", DescriptionCn: "可选docker、pouch,如果指定了container-id,不指定runtime则默认为docker", ValueType: "string", ValueRule: "docker,pouch"}
		ContainerArgsWaitTime = basic.Args{InjectId: containerFault.ID, ExecType: ExecInject, Key: "wait-time(s)", KeyCn: "容器重启可容忍耗时秒数", Description: "The number of seconds a container can tolerate restarting", DescriptionCn: "容器重启可容忍耗时秒数", ValueType: "int", ValueRule: ">0", DefaultValue: "10"}
	)

	args := []*basic.Args{&ContainerArgsId, &ContainerArgsRuntime}
	if isHaveWaitTimeArg {
		args = append(args, &ContainerArgsWaitTime)
	}

	return basic.InsertArgsMulti(ctx, args)
}
