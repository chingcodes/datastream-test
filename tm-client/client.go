package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chingcodes/datastream-test/tm-client/tm"
)

var (
	useJson bool

	grpc_addr string

	grpcCmd = &cobra.Command{
		Use:   "grpc",
		Short: "Connects to grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := tm.GrpcClient(grpc_addr)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			tm.ConsumeDataPoint(c)
		},
	}

	nats_addr, nats_subject string

	natsCmd = &cobra.Command{
		Use:   "nats",
		Short: "Subscribes to NATS.io service",
		Run: func(cmd *cobra.Command, args []string) {
			c, err := tm.SubscribeNats(nats_addr, nats_subject, useJson)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			tm.ConsumeDataPoint(c)
		},
	}

	udp_addr string

	udpCmd = &cobra.Command{
		Use:   "udp",
		Short: "Listens for udp packets",
		Run: func(cmd *cobra.Command, args []string) {
			//err := tm.SendUdp(udp_addr, useJson, hz, size)
			//if err != nil {
			//	fmt.Println(err)
			//	os.Exit(1)
			//}
		},
	}

	cmd = &cobra.Command{
		Use: "tm-client",
	}
)

func init() {
	grpcCmd.Flags().StringVar(&grpc_addr, "addr", "localhost:8081", "Grpc Server Address")

	natsCmd.Flags().BoolVar(&useJson, "json", false, "Use Json encoding")

	natsCmd.Flags().StringVar(&nats_addr, "addr", "localhost:4222", "Nats Server Address")
	natsCmd.Flags().StringVar(&nats_subject, "subject", "tm.timetest.1", "Nats Subject to subscribe to")

	udpCmd.Flags().BoolVar(&useJson, "json", false, "Use Json encoding")

	udpCmd.Flags().StringVar(&udp_addr, "addr", "244.0.0.42:2042", "UDP address listen on")

	cmd.AddCommand(grpcCmd, natsCmd, udpCmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
