package main

import (
	"fmt"
	"net"
	"strconv"
	"time"

	"github.com/cesbo/go-mpegts"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/transform"
)

func process(udpAddr string, udpPort string) {

	port, _ := strconv.Atoi(udpPort)
	sdt := map[string]string{}
	conn, _ := openSocket4(nil, net.ParseIP(udpAddr), port)
	defer conn.Close()

	conn.SetReadDeadline(time.Now().Add(200 * time.Millisecond))

	var slicer mpegts.Slicer
	buffer := make([]byte, 32*1024)
	cnt := 0
	for {

		n, _, _ := conn.ReadFrom(buffer)
		for packet := slicer.Begin(buffer[:n]); packet != nil; packet = slicer.Next() {
			pid := packet.PID()
			if pid == 0x11 {
				conn.Close()
				length_prov := int(packet[24] - 1)
				sdt["provider"] = string(packet[26 : 26+length_prov])

				sdt["servicename"] = string(packet[26+length_prov+1 : int(packet[7]+4)])
				sdt["pnr"] = strconv.Itoa(int(packet[16])<<8 | int(packet[17]))

				if packet[27+length_prov] == 0x01 {
					decoder := charmap.ISO8859_5.NewDecoder()
					as, _, _ := transform.Bytes(decoder, []byte(sdt["servicename"]))
					sdt["servicename"] = string(as)
				}
				fmt.Println(udpAddr + ":" + udpPort + "\t" + sdt["pnr"] + "\t" + sdt["provider"] + "\t" + sdt["servicename"])
				return
			}
		}
		cnt++
		if cnt > 1000 {
			fmt.Println(udpAddr + ":" + udpPort + "\t?")
			return
		}
		// if err := slicer.Err(); err != nil {
		// 	fmt.Println(err)
		// }
	}

}
