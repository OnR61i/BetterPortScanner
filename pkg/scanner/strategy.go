package pkg/scanner

import(
	"net"
	"time"
)

type Strategy interface{
	Scan(targetIp net.IP, srcIp net.IP, targetPort int, intrfc *net.Interface, timeOut time.Duration) (bool, error);
}


