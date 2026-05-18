package main

import(
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"syscall"
	"net"
	"os"
	"log"
	"strings"
//	"strconv"
	"regexp"
	"errors"
//	"reflect"
)

type TCPStrategy struct{
}

func scan(target net.IP, port int) error{
	// --> Getting mac addresses...
	intrfc, err := net.InterfaceByName("eno1")
	if err != nil {
		log.Fatal(err)
	}

	srcMac := intrfc.HardwareAddr 

	dstMac, err := getMacAddr(target)
	if(err != nil){
		log.Fatal(err)
	}

 	// Creating networklayer...
	eth := layers.Ethernet{
		SrcMAC: srcMac,
		DstMAC: dstMac,
		EthernetType: layers.EthernetTypeIPv4,
	}

	// Creating iplayer...
	srcIp, err := getLocalIp(intrfc)
	if(err != nil){
		log.Fatal(err)
	}

	if(err != nil){
		log.Fatal(err)
	}

	dstIp := target

	ipv4 := layers.IPv4{
		SrcIP: srcIp,
		DstIP: dstIp,
	}

	// Creating tcplayer...
	tcp := layers.TCP{
		SrcPort: layers.TCPPort(666),
		DstPort: layers.TCPPort(port),
		SYN:	 true,
	}
	
	tcp.SetNetworkLayerForChecksum(&ipv4)

	// Serializing into bytestream...
	buf := gopacket.NewSerializeBuffer()
	ops := gopacket.SerializeOptions{
		ComputeChecksums: true,
		FixLengths: 	true,
	}
	err = gopacket.SerializeLayers(
		buf,
		ops,
		&eth,
		&ipv4,
		&tcp,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Prepared for request")
	packetData := buf.Bytes()

	//Sending bytestream...
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	log.Print("Set up socket")
	if err != nil {
		log.Fatal(err)
	}

	err = syscall.BindToDevice(fd, "eno1")
    if err != nil {
        syscall.Close(fd)
        panic(err)
    }

	go sendByteStream(srcIp, port, dstMac, packetData, *intrfc, fd)

	// Waiting for response...
	// Pls make shure to deactivate your Firewall with: sudo iptables -A OUTPUT -p tcp --tcp-flags RST RST -j DROP
	packet, err := getResponse(fd, dstIp, port)
	if err != nil{
		log.Fatal(err)
	}
	log.Print(packet)
	return nil
}

func sendByteStream(ip net.IP, port int, dstMac net.HardwareAddr, buf []byte, ifi net.Interface, fd int) error {
	addr := syscall.SockaddrLinklayer{
		Protocol: 	htons(syscall.ETH_P_ALL),
		Ifindex: 	ifi.Index,
		Halen: 		6,
	}

	copy(addr.Addr[:6], dstMac)
	for{
	err := syscall.Sendto(fd, buf, 0, &addr)
	log.Print("Send stream")
	if err != nil {
		log.Fatal(err)
	}
	log.Print("Check")
}
	return nil
}

func getResponse(fd int, dstIp net.IP, port int) (gopacket.Packet, error) {

		buf := make([]byte, 64*1024)

	
	bufOob := make([]byte, 64*1024)
 	log.Print("Waiting for response")
	for{
			
		buf = make([]byte, 64*1024)
		recvBts, _, _, _, _ := syscall.Recvmsg(fd, buf, bufOob, 0)
		log.Println(buf[:recvBts])
	}
				
}

func getMacAddr(host net.IP) (net.HardwareAddr, error) {

	// Extracting content from arp file...
	file, err := os.Open("/proc/net/arp")
	if err != nil {
		log.Fatal(err)
	}

	data := make([]byte, 1000)
	_, err = file.Read(data)
	if err != nil{
		log.Fatal(err)
	}

	content := string(data)
	
	// Scanning for mac addresses...
	lines := strings.Split(content, "\n")
	
	hostStr := host.String()
	for _, line := range lines{

		if(strings.Contains(line, hostStr)){		

			regex := regexp.MustCompile("([0-9A-Fa-f]{2}[-:]){5}([0-9A-Fa-f]{2})")
			mac := regex.FindString(line)

			return net.ParseMAC(mac)
		}
	}
	return nil, errors.New("Wasn't able to find macaddr")
}

func getLocalIp(i *net.Interface) (net.IP, error) {

	addrs, err := i.Addrs()
	if err != nil {
		log.Fatal(err)
	}

	for _, addr := range addrs {
		
		log.Print(addr)
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP, nil
			}
		}
	}
	return nil, errors.New("Wasn't able to find local ip") 
}

func htons(in uint16) uint16 {
	return (in<<8)&0xff00 | in>>8
}
 
/*
func getSrcIp(interfaces []net.Interface) net.IP {
	
	for _, i := range interfaces {

		addrs, err := i.Addr()
		if err != nil {
			log.Fatal(err)
		}
		
		for _, addr := range addrs {

			var ip net.IP
			switch v := addr.(type) {
				case: *net.IPNet:
				ip = v.IP
				case: *net.IPAddr:
				ip = v.IP
			}
*/
