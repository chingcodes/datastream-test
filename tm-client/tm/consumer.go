package tm

import (
	"fmt"

	"github.com/chingcodes/datastream-test/pb"
	"google.golang.org/protobuf/encoding/prototext"
)

func ConsumeDataPoint(c chan (*pb.DataPoint)) {
	for dp := range c {
		fmt.Print(prototext.Format(dp))
	}
}
