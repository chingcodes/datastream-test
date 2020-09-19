package tm

import (
	"context"
	"fmt"
	"sync"
	"time"

	pb "github.com/chingcodes/datastream-test/pb"
)

var (
	dpChans     = map[chan (*pb.DataPoint)]struct{}{}
	dpChansLock sync.Mutex

	dpWorkerLock   sync.Mutex
	dpWorkerActive bool
)

func addChannel(c chan (*pb.DataPoint)) {
	dpChansLock.Lock()
	defer dpChansLock.Unlock()
	dpChans[c] = struct{}{}
}

func rmChannel(c chan (*pb.DataPoint)) {
	dpChansLock.Lock()
	defer dpChansLock.Unlock()
	delete(dpChans, c)
	close(c)
}

func startWorker(hz, size int) {
	dpWorkerLock.Lock()
	defer dpWorkerLock.Unlock()

	if !dpWorkerActive {
		go func() {
			t := time.NewTicker(time.Second / time.Duration(hz))
			defer t.Stop()

			dp := &pb.DataPoint{
				Dummy: make([]byte, size),
			}

			for {
				dp.TimeNs = uint64(time.Now().UnixNano())
				dpChansLock.Lock()
				for c, _ := range dpChans {
					select {
					case c <- dp:
					default:
						//skip if blocked
						fmt.Println("Datapoint dropped")
					}
				}
				dpChansLock.Unlock()
				<-t.C // wait for next tick
			}
		}()
		dpWorkerActive = true
	}
}

func NewGen(ctx context.Context, hz int, size int) chan (*pb.DataPoint) {
	startWorker(hz, size) // Can only start 1st
	c := make(chan (*pb.DataPoint))

	addChannel(c)

	go func() {
		<-ctx.Done()
		rmChannel(c)
	}()

	return c
}
