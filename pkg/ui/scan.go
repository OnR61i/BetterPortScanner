package pkg/ui

import(
)

type Scan struct{
	TargetRange []string
	PortRange string
	Strategy string
	NetInterface string
	Timeout string
}

func NewScan(target string, portRange string, strategy string, netInterface string, timeout string) {
	return Scan{
		Target: 	target,
		PortRange: 	portRange,
		Strategy:	strategy,
		NetInterface:	netInterface,
		Timeout: 	timeout,	
	}
}


