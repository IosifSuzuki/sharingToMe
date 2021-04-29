package sysInfoManager

import (
	"IosifSuzuki/sharingToMe/internal/models"
	"fmt"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"runtime"
	"strings"
)

func GetSysInfo() (*models.SysInfo, error) {
	cpuInfos, err := cpu.Info()
	if err != nil {
		return nil, err
	}
	percentage, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}
	var sysInfo *models.SysInfo = new(models.SysInfo)
	sysInfo.GolangVersion = runtime.Version()
	var infoAboutCPUs = make([]models.InfoAboutCPU, 0, 1)
	for _, cpuInfo := range cpuInfos {
		var infoAboutCPU models.InfoAboutCPU
		infoAboutCPU.Vendor = cpuInfo.VendorID
		infoAboutCPU.Family = cpuInfo.Family
		infoAboutCPU.NumberOfCores = cpuInfo.Cores
		infoAboutCPU.ModelName = cpuInfo.ModelName
		infoAboutCPU.Speed = cpuInfo.Mhz
		infoAboutCPU.CasheSize = cpuInfo.CacheSize
		infoAboutCPU.UtilizationPerCore = percentage
		infoAboutCPUs = append(infoAboutCPUs, infoAboutCPU)
	}
	sysInfo.CPUInfos = infoAboutCPUs

	hostInfo, err := host.Info()
	if err != nil {
		return sysInfo, err
	}
	var minutes = hostInfo.Uptime / 60 % 60
	var hours = hostInfo.Uptime / 3600 % 24
	var days = hostInfo.Uptime / 3600 / 24
	if days == 0 {
		sysInfo.UpTime = fmt.Sprintf("%02d:%02d hours", hours, minutes)
	} else {
		var daysText string
		if days == 1 {
			daysText = "day"
		} else {
			daysText = "days"
		}
		sysInfo.UpTime = fmt.Sprintf("%d %s %02d:%02d hours", days, daysText, hours, minutes)
	}
	sysInfo.NumberOfRunningProcesses = hostInfo.Procs
	sysInfo.OS = strings.Title(hostInfo.OS)

	diskInfo, err := disk.Usage("/")
	if err != nil {
		return sysInfo, err
	}
	sysInfo.StorageInfo.TotalDiskSpace = diskInfo.Total
	sysInfo.StorageInfo.FreeDiskSpace = diskInfo.Free
	sysInfo.StorageInfo.UsedDiskSpace = diskInfo.Used
	sysInfo.StorageInfo.PercentageSpaceUsage = diskInfo.UsedPercent

	memoryInfo, err := mem.VirtualMemory()
	if err != nil {
		return sysInfo, err
	}
	sysInfo.RAMInfo.TotalMemory = memoryInfo.Total
	sysInfo.RAMInfo.FreeMemory = memoryInfo.Free
	sysInfo.RAMInfo.PercentageUsedMemory = memoryInfo.UsedPercent
	return sysInfo, nil
}
