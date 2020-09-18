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

	wg = sync.WaitGroup{}

	grpcCmd = &cobra.Command{
		Use: "grpc",
		Run: func(cmd *cobra.Command, args []string) {
			wg.Add(1)
			defer wg.Done()
			tm.RunGrpcServer(hz, size)
		},
	}

	cmd = &cobra.Command{
		Use: "tm-server",
	}
)

func init() {
	cmd.Flags().IntVarP(&hz, "hz", "", 1, "Hertz to run Generator at")
	cmd.Flags().IntVarP(&size, "size", "", 0, "Size in bytes of dummy payload")
	cmd.AddCommand(grpcCmd)
}

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
