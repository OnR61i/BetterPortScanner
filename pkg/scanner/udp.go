package scanner

import(
	"net"
	"time"
)

type UDP struct{
}

func (udp *UDP) Scan(targetIp net.IP, srcIp net.IP, targetPort int, intrfc *net.Interface, timeOut time.Duration) (bool, error){
	return true, nil
}
