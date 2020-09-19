package tm

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func PushToNats(addr string, subject string, useJson bool, hz, size int) error {
	fmt.Println("Connecting to NATS server")
	nc, err := nats.Connect(addr)
	if err != nil {
		return err
	}
	defer nc.Close()

	gen := NewGen(context.Background(), hz, size)

	var marshal func(proto.Message) ([]byte, error)
	if useJson {
		marshal = protojson.Marshal
	} else {
		marshal = proto.Marshal
	}

	for dp := range gen {
		buf, err := marshal(dp)
		if err != nil {
			return err
		}
		err = nc.Publish(subject, buf)
		if err != nil {
			return err
		}
	}
	return nil
}
