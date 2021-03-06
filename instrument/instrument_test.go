package instrument

import (
	"math/rand"
	"testing"
	"time"

	"github.com/golang/glog"
)

func Test(t *testing.T) {
	StartMetrics()
	for {
		v := (rand.NormFloat64() * 200) + 10
		SuccessDuration <- &Measurements{
			Name:  "Normal",
			Value: v,
		}

		FailedCounter <- &Measurements{
			Name:  "FailedMethod",
			Value: v,
		}

		InProcess <- &Measurements{
			Name:  "InProcessMethod",
			Value: 1,
		}

		TransferRate <- &Measurements{
			Name:  "GetFile",
			Value: v,
		}

		time.Sleep(time.Duration(100) * time.Millisecond)

		q, err := GetTransferRateQuantile("GetFile", 0.99)
		if err != nil {
			glog.Infof("error %v", err)
		}
		glog.Infof("0.99: %v", q)
	}
}
