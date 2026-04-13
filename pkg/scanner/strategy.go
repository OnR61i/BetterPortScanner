package pkg/scanner

import(
	"net"
)

type Strategy interface{
}

func scan(target net.IP, port int) error{
}
