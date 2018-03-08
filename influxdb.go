package influxdb

import (
	"github.com/devopsfaith/krakend/config"
	"github.com/pkg/errors"
	"github.com/influxdata/influxdb/client/v2"
	metrics "github.com/devopsfaith/krakend-metrics/gin"
	"time"
	"context"
	"fmt"
	"github.com/devopsfaith/krakend/logging"
)

const Namespace = "github_com/letgoapp/krakend-influx"

type influxConfig struct {
	address  string
	username string
	password string
	ttl      time.Duration
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

func New(ctx context.Context, extraConfig config.ExtraConfig, metricsCollector *metrics.Metrics, logger logging.Logger) error {
	logger.Debug("Entering new")
	fmt.println("pepe")
	cfg, ok := configGetter(extraConfig).(influxConfig)

	if !ok {
		logger.Debug("no config")
		return errNoConfig
	}

	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.address,
		Username: cfg.username,
		Password: cfg.password,
	})

	if err != nil {
		logger.Debug("client crashed")
		return err
	}

	t := time.NewTicker(cfg.ttl)

	go keepUpdated(ctx, t.C, influxdbClient, metricsCollector)

	logger.Debug("client up and running")

	return nil
}

func keepUpdated(ctx context.Context, ticker <-chan time.Time, influxdbClient client.Client, metricsCollector *metrics.Metrics) {
	for {
		select {
		case <-ticker:
		case <-ctx.Done():
			return
		}

		fmt.Println("Preparing points")

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

		fmt.Println("WRITE done")
	}
}
