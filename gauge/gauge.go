package gauge

import (
	"time"

	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client/v2"
)

func Points(hostname string, now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := make([]*client.Point, 2)

	in := map[string]interface{}{
		"gauge": int(counters["krakend.router.connected-gauge"]),
	}
	incoming, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "in"}, in, now)
	if err != nil {
		logger.Error("creating incoming connection counters point:", err.Error())
		return res
	}
	res[0] = incoming

	out := map[string]interface{}{
		"gauge": int(counters["krakend.router.disconnected-gauge"]),
	}
	outgoing, err := client.NewPoint("router", map[string]string{"host": hostname, "direction": "out"}, out, now)
	if err != nil {
		logger.Error("creating outgoing connection counters point:", err.Error())
		return res
	}
	res[1] = outgoing

	return res
}
