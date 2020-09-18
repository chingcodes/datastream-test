package tm

import (
	"context"
	"time"

	pb "github.com/chingcodes/datastream-test/pb"
)

func NewGen(ctx context.Context, hz int, size int) chan (*pb.DataPoint) {
	c := make(chan (*pb.DataPoint))

	go func() {
		defer close(c)

		t := time.NewTicker(time.Second / time.Duration(hz))
		defer t.Stop()

		dp := &pb.DataPoint{
			Dummy: make([]byte, size),
		}

		done := ctx.Done()

		for {
			select {
			case <-t.C:
				dp.TimeNs = uint64(time.Now().UnixNano())
				select {
				case c <- dp:
					//do nothing
				case <-done:
					return
				}
			case <-done:
				return
			}
		}
	}()

	return c
}
