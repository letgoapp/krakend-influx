package gauge

import (
	"regexp"
	"time"

	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client"
)

var (
	connectionCounterPattern = `krakend\.router\.(dis)?connected-gauge`
	connectionCounterRegexp  = regexp.MustCompile(connectionCounterPattern)
)

func Points(now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !connectionCounterRegexp.MatchString(k) {
			logger.Debug("counter-gauge", k, "doesn't match")
			continue
		}
		fields := map[string]interface{}{
			"count": int(count),
		}

		countersPoint, err := client.NewPoint(k[15:], map[string]string{}, fields, now)
		if err != nil {
			logger.Error("creating counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	return res
}
