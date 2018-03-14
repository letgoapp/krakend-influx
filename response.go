package influxdb

import (
	"regexp"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	responseCounterPattern = `krakend\.router\.response\.([\/_a-zA-Z]+)\.status\.([\d]{3})\.count`
	responseCounterRegexp  = regexp.MustCompile(responseCounterPattern)
)

func (cw clientWrapper) responseCounterValues(now time.Time, counters map[string]int64) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !responseCounterRegexp.MatchString(k) {
			cw.logger.Debug("counter", k, "doesn't match")
			continue
		}
		params := responseCounterRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"name":   params[0],
			"status": params[1],
		}
		fields := map[string]interface{}{
			"count": int(count),
		}

		countersPoint, err := client.NewPoint("responses", tags, fields, now)
		if err != nil {
			cw.logger.Error("creating counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	return res
}
