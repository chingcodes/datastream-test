package tm

import (
	"fmt"
	"net"

	pb "github.com/chingcodes/datastream-test/pb"
	"google.golang.org/grpc"
)

type DataStreamGen struct {
	pb.UnimplementedDataStreamServiceServer
	hz   int
	size int
}

func RunGrpcServer(hz int, size int) error {
	port := 8081
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()

	dsg := &DataStreamGen{
		hz:   hz,
		size: size,
	}

	pb.RegisterDataStreamServiceServer(grpcServer, dsg)

	fmt.Printf("Starting GRPC Server on port %d\n", port)
	return grpcServer.Serve(lis)
}

func (s *DataStreamGen) Subscribe(req *pb.SubscribeReq, srv pb.DataStreamService_SubscribeServer) (err error) {
	fmt.Printf("Got GRPC Subscribe request: %v\n", req)

	defer func() {
		fmt.Printf("Subscribe GRPC finished: %v\n", err)
	}()

	gen := NewGen(srv.Context(), s.hz, s.size)

	for dp := range gen {
		err := srv.Send(dp)
		if err != nil {
			return err
		}
	}
	return nil
}
