package pkg/scanner

import(
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	
	"net"
	"errors"
	"syscall"
	"strings"
	"reflect"
	"regexp"
	"time"
//	"log"
	"os"
)

const ARP_FILE_PATH = "/proc/net/arp"

var MAC_REGEX = regexp.MustCompile("([0-9A-Fa-f]{2}[-:]){5}([0-9A-Fa-f]{2})")

var TIMEOUT_ARP_REQUEST = 3 * time.Second

func ResolveTargetHwAddr(ifindex int, fd int, targetIp net.IP, srcIp net.IP) (net.HardwareAddr, error, error) {
	// Try to detect targets MAC-Address with given IP-Address...
	if mac, err1 := scanARPFile(targetIp); mac == nil {
		if mac, err2 := broadcastARPRequest(ifindex, fd, targetIp, srcIp); mac == nil {
			return nil, err1, err2
		} else {
			return mac, err1, nil
		}
	} else { 
		return mac, nil, nil
	}
}

func scanARPFile(target net.IP) (net.HardwareAddr, error) {
	// Scanning ARP-Cache for targets MAC-Address...
	// --> Extracting content from ARP-Cache...
	file, err := os.Open(ARP_FILE_PATH)
	if err != nil {
		return nil, errors.New("Unable to read ARP-Cache")
	}
	
	// --> Recieving data...
	data := make([]byte, 10000)
	_, err = file.Read(data)
	if err != nil {
		return nil, errors.New("Unable to recieve ARP-Cache-Data")
	}

	content := string(data)

	// --> Scanning for MAC-Address...
	lines := strings.Split(content, "\n")

	targetStr := target.String()
	for _, line := range lines {
		if (strings.Contains(line, targetStr)){
			mac := MAC_REGEX.FindString(line)
			hwAddr, err := net.ParseMAC(mac)
			if err != nil {
				return nil, errors.New("Unable to parse found HwAddr string to net.HardwareAddr")
			} else {
				return hwAddr, nil
			}
		}
	}
	return nil, nil
}

// Attention: This function only works for forgein mashines. If you broadcast for the sender mashine there will be a difference in answer...
// To fix it adapt listenForARPResponse to this behavior...
func broadcastARPRequest(ifindex int, fd int, targetIp net.IP, srcIp net.IP) (net.HardwareAddr, error) {
	// Resolving targets MAC-Address with ARP-Protocol...
	// --> Creating ARP-Request...
	intrfc, err := net.InterfaceByIndex(ifindex)
	if err != nil {
		return nil, errors.New("Unable to resolve network interface by index")
	}

	eth := &layers.Ethernet{
		SrcMAC: 	intrfc.HardwareAddr,
		DstMAC: 	net.HardwareAddr([]byte{
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff,
				}),
		EthernetType:	layers.EthernetType(0x0806),
	}

	arp := &layers.ARP{
		AddrType: 	layers.LinkType(0x01),
		Protocol: 	layers.EthernetType(0x0800), 	//EthernetTypeARP EthernetType = 0x0806
		HwAddressSize: 	0x06,
		ProtAddressSize: 0x04,
		Operation: 	0x01,
		SourceHwAddress:	intrfc.HardwareAddr,
		SourceProtAddress:	srcIp,	
		DstHwAddress:	[]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00,},
		DstProtAddress:	targetIp,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:	false,
		ComputeChecksums:	false,
	}

	err = gopacket.SerializeLayers(
		buf,
		opts,
		eth,
		arp,
	)

	data := buf.Bytes()

	// --> Start listening for ARP-Responses...
	responseChan := make(chan net.HardwareAddr)
	stopChan := make(chan bool)

	go listenForARPResponse(ifindex, intrfc.HardwareAddr, srcIp, targetIp, responseChan, stopChan)

	// --> Prepare ARP-Request...
	addr := &syscall.SockaddrLinklayer{
		Ifindex: 	intrfc.Index,
	}

	// --> Save time for Packet Capture Live opening...
	time.Sleep(2 * time.Second)

	// --> Sending ARP-Request...
	err = syscall.Sendto(fd, data, 0, addr)
	if err != nil {
		return nil, errors.New("Unable to send ARP-Request")
	}

	// --> Waiting for ARP-Response...
	go func(){
		time.Sleep(TIMEOUT_ARP_REQUEST)
		close(stopChan)
		close(responseChan)
	}()

	for resp := range responseChan {
		// --> Checking ARP-Response
		dstMac := resp
		if(dstMac != nil){
			return dstMac, nil
		} else {
			return nil, nil
		}
	}

	return nil, nil
}

func listenForARPResponse(ifindex int, srcMac []byte, srcIp []byte, dstIp []byte, responseChan chan net.HardwareAddr, stopChan chan bool) {
	// Packet Capture for ARP-Response...
	intrfc, err := net.InterfaceByIndex(ifindex)
	if err != nil {
		responseChan <- net.HardwareAddr([]byte{
			0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		})
		return
	}
	device := intrfc.Name

	// --> Open Live for Packet Capture...
	handle, err := pcap.OpenLive(device, 1600, true, pcap.BlockForever)

	if err != nil {
		responseChan <- net.HardwareAddr([]byte{
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				})
		return 
	}

	defer handle.Close()

	// --> Create Packet Source...
	packetSrc := gopacket.NewPacketSource(handle, handle.LinkType())
			
	for {
		select{
		case <- stopChan:
			return
		case packet := <- packetSrc.Packets():
			// --> Checking if Packet is Response...
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer != nil {
				arp, _ := arpLayer.(*layers.ARP)
				if (arp.Operation == uint16(0x02)) {
					if (reflect.DeepEqual(arp.DstHwAddress, srcMac)) {
						if (reflect.DeepEqual(arp.DstProtAddress, srcIp)) {
							if (reflect.DeepEqual(arp.SourceProtAddress, dstIp)) {
								responseChan <- net.HardwareAddr(arp.SourceHwAddress)
								break
							}
						}
					}
				}
			}

		}
	}
}
