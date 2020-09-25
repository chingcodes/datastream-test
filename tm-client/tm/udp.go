package tm

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/chingcodes/datastream-test/pb"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func ListenUdp(udp_addr string, useJson bool) error {
	addr := strings.Split(udp_addr, ":")

	ip := net.ParseIP(addr[0])
	port, _ := strconv.Atoi(addr[1])

	udpAddr := &net.UDPAddr{IP: ip, Port: port}

	var conn *net.UDPConn
	var err error

	if ip.IsMulticast() {
		fmt.Println("Using mulitcast")
		conn, err = net.ListenMulticastUDP("udp", nil, udpAddr)
		if err != nil {
			return err
		}
	} else {
		conn, err = net.ListenUDP("udp", udpAddr)
		if err != nil {
			return err
		}
	}
	defer conn.Close()

	var unmarshal func([]byte, proto.Message) error
	if useJson {
		unmarshal = protojson.Unmarshal
	} else {
		unmarshal = proto.Unmarshal
	}

	dp := &pb.DataPoint{}

	buf := make([]byte, 10000) // 10KB
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			return err
		}
		err = unmarshal(buf[:n], dp)
		if err != nil {
			return err
		}

		fmt.Print(prototext.Format(dp))
	}
}
