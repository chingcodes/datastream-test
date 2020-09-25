package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/chingcodes/datastream-test/tm-server/tm"
)

var (
	hz, size int

	useJson bool

	grpcCmd = &cobra.Command{
		Use:   "grpc",
		Short: "Starts grpc server",
		Run: func(cmd *cobra.Command, args []string) {
			err := tm.RunGrpcServer(hz, size)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	nats_addr, nats_subject string

	natsCmd = &cobra.Command{
		Use:   "nats",
		Short: "Start NATS.io service",
		Run: func(cmd *cobra.Command, args []string) {
			err := tm.PushToNats(nats_addr, nats_subject, useJson, hz, size)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	udp_addr string

	udpCmd = &cobra.Command{
		Use:   "udp",
		Short: "Start broadcasting udp packets",
		Run: func(cmd *cobra.Command, args []string) {
			err := tm.SendUdp(udp_addr, useJson, hz, size)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		},
	}

	cmd = &cobra.Command{
		Use: "tm-server",
	}
)

func init() {
	grpcCmd.Flags().IntVar(&hz, "hz", 1, "Hertz to run Generator at")
	grpcCmd.Flags().IntVar(&size, "size", 0, "Size in bytes of dummy payload")

	natsCmd.Flags().IntVar(&hz, "hz", 1, "Hertz to run Generator at")
	natsCmd.Flags().IntVar(&size, "size", 0, "Size in bytes of dummy payload")
	natsCmd.Flags().BoolVar(&useJson, "json", false, "Use Json encoding")

	natsCmd.Flags().StringVar(&nats_addr, "addr", "localhost:4222", "Nats Server Address")
	natsCmd.Flags().StringVar(&nats_subject, "subject", "tm.timetest.1", "Nats Subject to push to")

	udpCmd.Flags().IntVar(&hz, "hz", 1, "Hertz to run Generator at")
	udpCmd.Flags().IntVar(&size, "size", 0, "Size in bytes of dummy payload")
	udpCmd.Flags().BoolVar(&useJson, "json", false, "Use Json encoding")

	udpCmd.Flags().StringVar(&udp_addr, "addr", "224.0.0.42:4242", "UDP address to send to")

	cmd.AddCommand(grpcCmd, natsCmd, udpCmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
