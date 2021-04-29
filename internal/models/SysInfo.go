package models

type InfoAboutCPU struct {
	Vendor string
	Family string
	NumberOfCores int32
	ModelName string
	Speed float64
	UtilizationPerCore []float64
	CasheSize int32
	UpTime string
}

type StorageInfo struct {
	TotalDiskSpace uint64
	UsedDiskSpace uint64
	FreeDiskSpace uint64
	PercentageSpaceUsage float64
	FileSystem string
}

type RAMInfo struct {
	TotalMemory uint64
	FreeMemory uint64
	PercentageUsedMemory float64
}

type SysInfo struct {
	GolangVersion string
	CPUInfos []InfoAboutCPU
	UpTime string
	NumberOfRunningProcesses uint64
	OS string
	StorageInfo StorageInfo
	RAMInfo RAMInfo
}
