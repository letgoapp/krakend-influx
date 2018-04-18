package counter

import (
	"regexp"
	"time"

	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client"
)

func Points(now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	points := requestPoints(now, counters, logger)
	points = append(points, responsePoints(now, counters, logger)...)
	points = append(points, connectionPoints(now, counters, logger)...)
	return points
}

var (
	connectionCounterPattern = `krakend\.router\.(dis)?connected(-total)?`
	connectionCounterRegexp  = regexp.MustCompile(connectionCounterPattern)

	requestCounterPattern = `krakend\.proxy\.requests\.layer\.([a-zA-Z]+)\.name\.([\/_a-zA-Z]+)\.complete\.(true|false)\.error\.(true|false)`
	requestCounterRegexp  = regexp.MustCompile(requestCounterPattern)

	responseCounterPattern = `krakend\.router\.response\.([\/_a-zA-Z]+)\.status\.([\d]{3})\.count`
	responseCounterRegexp  = regexp.MustCompile(responseCounterPattern)
)

func requestPoints(now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !requestCounterRegexp.MatchString(k) {
			logger.Debug("counter", k, "doesn't match")
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
			logger.Error("creating counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	return res
}

func responsePoints(now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !responseCounterRegexp.MatchString(k) {
			logger.Debug("counter", k, "doesn't match")
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
			logger.Error("creating counters point:", err.Error())
			continue
		}
		res = append(res, countersPoint)
	}
	return res
}

func connectionPoints(now time.Time, counters map[string]int64, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, count := range counters {
		if !connectionCounterRegexp.MatchString(k) {
			logger.Debug("counter", k, "doesn't match")
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
