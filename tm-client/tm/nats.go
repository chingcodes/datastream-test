package tm

import (
	"context"
	"fmt"

	"github.com/chingcodes/datastream-test/pb"
	"github.com/nats-io/nats.go"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func SubscribeNats(nats_addr, nats_subject string, useJson bool) (chan (*pb.DataPoint), error) {
	nc, err := nats.Connect(nats_addr)
	if err != nil {
		return nil, err
	}

	sub, err := nc.SubscribeSync(nats_subject)
	if err != nil {
		nc.Close()
		return nil, err
	}

	var unmarshal func([]byte, proto.Message) error
	if useJson {
		unmarshal = protojson.Unmarshal
	} else {
		unmarshal = proto.Unmarshal
	}

	c := make(chan (*pb.DataPoint))

	go func() {
		defer nc.Close()

		dp := &pb.DataPoint{}

		for {
			msg, err := sub.NextMsgWithContext(context.Background())
			if err != nil {
				fmt.Println(err)
				return
			}

			err = unmarshal(msg.Data, dp)
			if err != nil {
				fmt.Println(err)
			}
			c <- dp
		}

		fmt.Println("sub ended")
	}()

	return c, nil
}
