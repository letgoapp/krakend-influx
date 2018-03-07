package influxdb

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/pkg/errors"
	"github.com/influxdata/influxdb/client/v2"
	"github.com/devopsfaith/krakend-metrics"
	"time"
	"context"
)

const Namespace = "github_com/letgoapp/krakend-influx"

type influxConfig struct {
	address  string
	username string
	password string
	ttl time.Duration
}

func configGetter(extraConfig config.ExtraConfig) interface{} {
	value, ok := extraConfig[Namespace]

	if !ok {
		return nil
	}

	castedConfig, ok := value.(map[string]interface{})

	if !ok {
		return nil
	}

	cfg := influxConfig{}

	if value, ok := castedConfig["address"]; ok {
		cfg.address = value.(string)
	}

	if value, ok := castedConfig["username"]; ok {
		cfg.username = value.(string)
	}

	if value, ok := castedConfig["password"]; ok {
		cfg.password = value.(string)
	}

	if value, ok := castedConfig["ttl"]; ok {
		s, ok := value.(string)

		if !ok {
		    return nil
		}
		var err error
		cfg.ttl, err = time.ParseDuration(s)

		if err != nil {
			return nil
		}
	}

	return cfg
}

var errNoConfig = errors.New("Unable to load custom config from the extra config")

func New(ctx context.Context, extraConfig config.ExtraConfig, metricsCollector *metrics.Metrics) error {
	cfg, ok := configGetter(extraConfig).(influxConfig)

	if !ok {
	    return errNoConfig
	}

	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.address,
		Username: cfg.username,
		Password: cfg.password,
	})

	if err != nil {
		return err
	}

	go keepUpdated(ctx, influxdbClient, metricsCollector)
	return nil
}

func keepUpdated(ctx context.Context, influxdbClient client.Client, metricsCollector *metrics.Metrics) {
	for {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return
		}

		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  "supu",
			Precision: "s",
		})

		snapshot := metricsCollector.Snapshot()

		fields := make(map[string]interface{}, len(snapshot.Counters))

		for k, v := range snapshot.Counters {
			fields[k] = v
		}

		now := time.Unix(snapshot.Time, 0)

		countersPoint, _ := client.NewPoint("counters", map[string]string{"type": "counters"}, fields, now)

		fields = make(map[string]interface{}, len(snapshot.Gauges))

		for k, v := range snapshot.Gauges {
			fields[k] = v
		}

		gaugesPoint, _ := client.NewPoint("gauges", map[string]string{"type": "gauges"}, fields, now)

		bp.AddPoint(countersPoint)
		bp.AddPoint(gaugesPoint)

		influxdbClient.Write(bp)
	}
}
