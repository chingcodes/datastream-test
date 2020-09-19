package tm

import (
	"context"
	"net"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func SendUdp(addr string, useJson bool, hz, size int) error {
	conn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return err
	}
	defer conn.Close()

	dst, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	var marshal func(proto.Message) ([]byte, error)
	if useJson {
		marshal = protojson.Marshal
	} else {
		marshal = proto.Marshal
	}

	for dp := range NewGen(context.Background(), hz, size) {
		buf, err := marshal(dp)
		if err != nil {
			return err
		}
		_, err = conn.WriteTo(buf, dst)
		if err != nil {
			return err
		}
	}
	return nil
}
