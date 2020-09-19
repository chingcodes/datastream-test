package tm

import (
	"context"
	"fmt"

	"github.com/chingcodes/datastream-test/pb"
	"google.golang.org/grpc"
)

func GrpcClient(addr string) (chan (*pb.DataPoint), error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	client := pb.NewDataStreamServiceClient(conn)

	sub, err := client.Subscribe(context.Background(), &pb.SubscribeReq{})
	if err != nil {
		return nil, err
	}

	c := make(chan (*pb.DataPoint))

	go func() {
		defer close(c)
		defer conn.Close()

		for {
			dp, err := sub.Recv()
			if err != nil {
				fmt.Println(err)
				return
			}
			c <- dp
		}
	}()

	return c, nil
}
