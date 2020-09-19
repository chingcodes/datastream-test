package main

import (
	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"

	"github.com/chingcodes/datastream-test/tm-server/tm"
)

var (
	hz, size int

	useJson bool

	wg = sync.WaitGroup{}

	grpcCmd = &cobra.Command{
		Use: "grpc",
		Run: func(cmd *cobra.Command, args []string) {
			wg.Add(1)
			defer wg.Done()
			tm.RunGrpcServer(hz, size)
		},
	}

	nats_addr, nats_subject string

	natsCmd = &cobra.Command{
		Use: "nats",
		Run: func(cmd *cobra.Command, args []string) {
			wg.Add(1)
			defer wg.Done()
			tm.PushToNats(nats_addr, nats_subject, useJson, hz, size)
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
	natsCmd.Flags().BoolVar(&useJson, "json", false, "Use Json encoding (not used for gprc)")

	natsCmd.Flags().StringVar(&nats_addr, "addr", "localhost:4222", "Nats Server Address")
	natsCmd.Flags().StringVar(&nats_subject, "subject", "tm.timetest.1", "Nats Subject to push to")

	cmd.AddCommand(grpcCmd, natsCmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
