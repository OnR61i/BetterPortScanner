package scanner

import(
	"net"
	"time"
)

type Job struct{
	Target net.IP
	Port int
	Strategy Strategy
	UsedInterface *net.Interface
	Timeout time.Duration
	State string
}

func NewJob(target net.IP, port int, strategy Strategy, intrfc *net.Interface, timeout time.Duration) Job{
	return Job{
		Target:		target,
		Port:		port,
		Strategy:	strategy,
		UsedInterface:	intrfc,
		Timeout: 	timeout,
		State:		"undefined",
	}
}

