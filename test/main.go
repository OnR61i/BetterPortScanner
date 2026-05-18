package main

import ( 
	"net"
	"log"
)

func main() {
	target:= net.ParseIP("192.168.178.37")

	portRange := 80	

	err := scan(target, portRange)

	if err != nil {
		log.Fatal(err)
	}
}
		
