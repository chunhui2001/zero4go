package x

// #cgo CFLAGS: -Wno-gnu-folding-constant
import "C"

import (
	"log"
	"runtime"

	"github.com/klauspost/cpuid/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"

	"github.com/chunhui2001/zero4go/pkg/utils"
)

func Info() {

	log.Printf("--------------------- CPU INFO ---------------------")
	log.Printf("Hostname: %s", utils.Hostname())
	log.Printf("OutboundIP: %s", utils.OutboundIP())

	// ---------------- 操作系统 ----------------
	hostInfo, _ := host.Info()
	log.Println("操作系统:", hostInfo.OS)

	// ---------------- 架构 ----------------
	log.Println("处理器架构:", runtime.GOARCH)

	// ---------------- CPU 信息 ----------------
	logicalCPU, _ := cpu.Counts(true)
	physicalCPU, _ := cpu.Counts(false)

	cpuInfos, _ := cpu.Info()

	log.Println("逻辑CPU数量:", logicalCPU)
	log.Println("物理CPU数量:", physicalCPU)

	if len(cpuInfos) > 0 {
		log.Println("总核数:", cpuInfos[0].Cores)
		log.Println("CPU厂商:", cpuInfos[0].VendorID)
		log.Println("CPU型号:", cpuInfos[0].ModelName)
	}

	// ---------------- 内存 ----------------
	vmem, _ := mem.VirtualMemory()
	swap, _ := mem.SwapMemory()

	log.Printf("总内存: %.2f GB\n", float64(vmem.Total)/1024/1024/1024)
	log.Printf("Swap(交换分区)总量: %.2f GB\n", float64(swap.Total)/1024/1024/1024)

	//---------------- CPU 指令集 ----------------
	log.Println("支持SSE?:", cpuid.CPU.Supports(cpuid.SSE))
	log.Println("支持AVX?:", cpuid.CPU.Supports(cpuid.AVX))

	// ---------------- 超线程 ----------------
	// 简单判断：逻辑核 > 物理核
	log.Println("支持HyperThreading?:", logicalCPU > physicalCPU)

	log.Printf("----------------------------------------------------")
}
