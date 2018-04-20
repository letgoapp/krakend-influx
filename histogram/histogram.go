package histogram

import (
	"regexp"
	"time"

	metrics "github.com/devopsfaith/krakend-metrics"
	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client/v2"
)

func Points(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	points := latencyPoints(hostname, now, histograms, logger)
	points = append(points, routerPoints(hostname, now, histograms, logger)...)
	if p := debugPoint(hostname, now, histograms, logger); p != nil {
		points = append(points, p)
	}
	if p := runtimePoint(hostname, now, histograms, logger); p != nil {
		points = append(points, p)
	}
	return points
}

var (
	latencyPattern = `krakend\.proxy\.latency\.layer\.([a-zA-Z]+)\.name\.(.*)\.complete\.(true|false)\.error\.(true|false)`
	latencyRegexp  = regexp.MustCompile(latencyPattern)

	routerPattern = `krakend\.router\.response\.(.*)\.(size|time)`
	routerRegexp  = regexp.MustCompile(routerPattern)
)

func latencyPoints(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, histogram := range histograms {
		if !latencyRegexp.MatchString(k) {
			continue
		}
		params := latencyRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"host":     hostname,
			"layer":    params[0],
			"name":     params[1],
			"complete": params[2],
			"error":    params[3],
		}
		fields := map[string]interface{}{
			"max":      int(histogram.Max),
			"min":      int(histogram.Min),
			"mean":     int(histogram.Mean),
			"stddev":   int(histogram.Stddev),
			"variance": int(histogram.Variance),
		}

		histogramPoint, err := client.NewPoint("requests", tags, fields, now)
		if err != nil {
			logger.Error("creating histogram point:", err.Error())
			continue
		}
		res = append(res, histogramPoint)
	}
	return res
}

func routerPoints(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) []*client.Point {
	res := []*client.Point{}
	for k, histogram := range histograms {
		if !routerRegexp.MatchString(k) {
			continue
		}
		params := routerRegexp.FindAllStringSubmatch(k, -1)[0][1:]
		tags := map[string]string{
			"host": hostname,
			"name": params[0],
		}
		fields := map[string]interface{}{
			"max":      int(histogram.Max),
			"min":      int(histogram.Min),
			"mean":     int(histogram.Mean),
			"stddev":   int(histogram.Stddev),
			"variance": int(histogram.Variance),
		}

		histogramPoint, err := client.NewPoint("router.response-"+params[1], tags, fields, now)
		if err != nil {
			logger.Error("creating histogram point:", err.Error())
			continue
		}
		res = append(res, histogramPoint)
	}
	return res
}

func debugPoint(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) *client.Point {
	hd, ok := histograms["krakend.service.debug.GCStats.Pause"]
	if !ok {
		return nil
	}
	tags := map[string]string{
		"host": hostname,
	}
	fields := map[string]interface{}{
		"max":      int(hd.Max),
		"min":      int(hd.Min),
		"mean":     int(hd.Mean),
		"stddev":   int(hd.Stddev),
		"variance": int(hd.Variance),
	}

	histogramPoint, err := client.NewPoint("service.debug.GCStats.Pause", tags, fields, now)
	if err != nil {
		logger.Error("creating histogram point:", err.Error())
		return nil
	}
	return histogramPoint
}

func runtimePoint(hostname string, now time.Time, histograms map[string]metrics.HistogramData, logger logging.Logger) *client.Point {
	hd, ok := histograms["krakend.service.runtime.MemStats.PauseNs"]
	if !ok {
		return nil
	}
	tags := map[string]string{
		"host": hostname,
	}
	fields := map[string]interface{}{
		"max":      int(hd.Max),
		"min":      int(hd.Min),
		"mean":     int(hd.Mean),
		"stddev":   int(hd.Stddev),
		"variance": int(hd.Variance),
	}

	histogramPoint, err := client.NewPoint("service.runtime.MemStats.PauseNs", tags, fields, now)
	if err != nil {
		logger.Error("creating histogram point:", err.Error())
		return nil
	}
	return histogramPoint
}
