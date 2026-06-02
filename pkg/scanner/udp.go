package pkg/scanner

import(
	"net"
)

type UDP struct{
}

func (udp *UDP) Scan(target net.IP, port int) {
}
