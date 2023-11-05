package shared

type SystemResult struct {
	TotalMemory uint64
	FreeMemory  uint64
	CPUUseage   float64
	Hostname    string
	Version     string //version of vm-manager-client
	Networks    []Network
}

type Network struct {
	IP   string
	MAC  string
	Name string
}
