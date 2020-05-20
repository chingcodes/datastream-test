package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"net"
	"time"

	pb "github.com/chingcodes/datastream-test/pb"
	"google.golang.org/grpc"
)

func main() {
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()

	pb.RegisterDataStreamServiceServer(grpcServer, &DataStreamGen{})

	fmt.Printf("Starting Server on port %d\n", port)
	grpcServer.Serve(lis)

}

type DataStreamGen struct {
}

func (s *DataStreamGen) Subscribe(req *pb.SubscribeReq, srv pb.DataStreamService_SubscribeServer) (err error) {
	fmt.Printf("Got Subscribe request: %v\n", req)

	defer func() {
		fmt.Printf("Subscribe finished: %v\n", err)
	}()

	var gen func() float64

	getTime := func() float64 { return time.Now().Sub(time.Unix(0, 0)).Seconds() }

	var last float64
	start := time.Now()

	switch req.Query {
	case "/counter/1":
		gen = func() float64 {
			last += 1
			return last
		}
	case "/constant/42":
		gen = func() float64 { return 42 }
	case "/math/sin":
		gen = func() float64 {
			last += 0.10
			return math.Sin(last)
		}
	case "/time/now":
		gen = func() float64 {
			return getTime()
		}
	case "/time/duration":
		gen = func() float64 {
			return time.Now().Sub(start).Seconds()
		}
	default:
		return errors.New("Unknown Query. Supported ones: '/counter/1', '/constant/42', '/time/now', '/time/duration', and '/math/sin'")
	}

	t := time.NewTicker(time.Millisecond * 250)

	srv.Send(&pb.DataPoint{
		Name:  req.Query,
		Time:  getTime(),
		Value: gen(),
	})

	for {
		select {
		case <-t.C:
			err := srv.Send(&pb.DataPoint{
				Time:  getTime(),
				Value: gen(),
			})
			if err != nil {
				return err
			}
		case <-srv.Context().Done():
			return nil
		}
	}
}
