package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	pb "github.com/chingcodes/datastream-test/pb"
	"google.golang.org/grpc"
)

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: ./client SERVERADDR:PORT QUERY\n")
		os.Exit(1)
	}

	conn, err := grpc.Dial(args[0], grpc.WithInsecure())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	client := pb.NewDataStreamServiceClient(conn)

	srv, err := client.Subscribe(context.Background(), &pb.SubscribeReq{Query: args[1]})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		dp, err := srv.Recv()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		fmt.Println(dp)
	}
}
