package pkg/scanner

import(
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"

	"net"
	"syscall"
	"reflect"
	"errors"
	"time"
)

const SRC_PORT = 51190

func Scan(ifindex int, targetIp net.IP, srcIp net.IP, targetPort int, timeOut time.Duration) (bool, error) {
	// Scanning target's IP-Address given port for state...
	// --> Creating Network-Socket, Type: RAW...
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(syscall.ETH_P_ALL))
	if err != nil {
		return false, errors.New("Critical Error! Unable to create Network-Socket")
	}

	// --> Preparing essential values...
	intrfc, err := net.InterfaceByIndex(ifindex)
	srcMac := intrfc.HardwareAddr
	dstMac, _, _ := ResolveTargetHwAddr(ifindex, fd, targetIp, srcIp)

	if dstMac == nil{
		return false, nil
	}
	// --> Creating TCP-Package...
	ethernetLayer := &layers.Ethernet{
		SrcMAC:		srcMac,
		DstMAC: 	dstMac,
		EthernetType:	layers.EthernetTypeIPv4,	// --> uint16(0x0800)
	}

	ipLayer := &layers.IPv4{
		SrcIP: 		srcIp,
		DstIP:		targetIp,
		Version: 	0x04,
		IHL:		0x05,
		Protocol:	layers.IPProtocolTCP,		// --> uint8(0x06)
		TTL: 		0x08,
	}

	tcpLayer := &layers.TCP{
		SrcPort: 	layers.TCPPort(SRC_PORT),
		DstPort: 	layers.TCPPort(targetPort),
		DataOffset:	5, 
		SYN:		true,
		Window:		0xffff,
	}

	// --> Use IP-Layer to compute Checksum for TCP-Layer...
	tcpLayer.SetNetworkLayerForChecksum(ipLayer)

	// --> Converting data to []byte...
	buffer := gopacket.NewSerializeBuffer()
	options := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths: true,
	}

	err = gopacket.SerializeLayers(
		buffer, 
		options,
		ethernetLayer,
		ipLayer,
		tcpLayer,
	)
	if err != nil {
		return false, errors.New("Critical Error! Unable to create Network-Socket")
	}

	packet := buffer.Bytes()

	// --> Setting up Response-Listener...
	responseChan := make(chan int) 		// --> Reponses: -1 (Error), 1 (Responded), 2 (!Responded)
	stopChan := make(chan bool)
	
	go listenForResponse(ifindex, fd, srcMac, dstMac, srcIp, targetIp, SRC_PORT, layers.TCPPort(targetPort), responseChan, stopChan)

	// --> Sending packet...
	addr := &syscall.SockaddrLinklayer{
		Ifindex:	ifindex, 
	}

	time.Sleep(time.Second)
	err = syscall.Sendto(fd, packet, 0, addr)
	if err != nil {
		return false, errors.New("Critical Error! Unable to send TCP-Packet")
	}
	
	// --> Waiting for response...
	go func(){
		time.Sleep(timeOut)
		close(responseChan)
		close(stopChan)
	}()

	for resp := range responseChan {

		if(resp == -1) {
			return false, errors.New("Critical Error! Unable to capture packages")
		}

		if(resp == 2) {
			return false, nil
		}
		if(resp == 1) {
			return true, nil
		}
	}

	return false, nil
}

func listenForResponse(ifindex int, fd int, srcMac net.HardwareAddr, dstMac net.HardwareAddr, srcIp net.IP, dstIp net.IP, srcPort layers.TCPPort, dstPort layers.TCPPort, responseChan chan int, stopChan chan bool) {
	// Packet Capture for target response...
	intrfc, err := net.InterfaceByIndex(ifindex)
	if err != nil {
		responseChan <- -1
		return
	}

	device := intrfc.Name

	// --> Open Live for Packet Capture...
	handle, err := pcap.OpenLive(device, 0x0640, true, pcap.BlockForever)
	if err != nil {
		responseChan <- -1 		// --> this boolean signals something went wrong.
		return 
	}
	defer handle.Close()

	// --> Creating Packet Source...
	packetSrc := gopacket.NewPacketSource(handle, handle.LinkType())
	
	for {
		select{
		case <- stopChan:
			responseChan <- 2
			return
		case packet := <- packetSrc.Packets():
			// --> Checking whether packet is response...
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			var ip *layers.IPv4
			if ipLayer != nil {
				ip, _ = ipLayer.(*layers.IPv4)
			} else {
				continue
			}

			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			var tcp *layers.TCP
			if tcpLayer != nil {
				tcp, _ = tcpLayer.(*layers.TCP)
			} else {
				continue
			}

			if(reflect.DeepEqual(ip.DstIP, srcIp)) {
				if(reflect.DeepEqual(ip.SrcIP, dstIp)) {
					if(reflect.DeepEqual(tcp.DstPort, srcPort)) {
						if(reflect.DeepEqual(tcp.SrcPort, dstPort)) {
							responseChan <- 1
						}
					}
				}
			}
		}
	}
}

