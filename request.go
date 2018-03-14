package influxdb

import (
	"regexp"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

var (
	requestCounterPattern = `krakend\.proxy\.requests\.layer\.([a-zA-Z]+)\.name\.([\/_a-zA-Z]+)\.complete\.(true|false)\.error\.(true|false)`
	requestCounterRegexp  = regexp.MustCompile(requestCounterPattern)
)

func (cw clientWrapper) requestCounterValues(now time.Time, counters map[string]int64) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !requestCounterRegexp.MatchString(k) {
			cw.logger.Debug("counter", k, "doesn't match")
			continue
		}
		params := requestCounterRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"layer":    params[0],
			"name":     params[1],
			"complete": params[2],
			"error":    params[3],
		}
		fields := map[string]interface{}{
			"count": int(count),
		}

		countersPoint, err := client.NewPoint("requests", tags, fields, now)
		if err != nil {
			cw.logger.Error("creating counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	return res
}
