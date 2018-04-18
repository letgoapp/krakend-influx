package influxdb

import (
	"context"
	"time"

	ginmetrics "github.com/devopsfaith/krakend-metrics/gin"
	"github.com/devopsfaith/krakend/config"
	"github.com/devopsfaith/krakend/logging"
	"github.com/influxdata/influxdb/client/v2"

	"github.com/letgoapp/krakend-influx/counter"
	"github.com/letgoapp/krakend-influx/gauge"
)

const Namespace = "github_com/letgoapp/krakend-influx"

type clientWrapper struct {
	influxClient client.Client
	collector    *ginmetrics.Metrics
	logger       logging.Logger
	db           string
}

func New(ctx context.Context, extraConfig config.ExtraConfig, metricsCollector *ginmetrics.Metrics, logger logging.Logger) error {
	logger.Debug("creating a new ifluxdb client")
	cfg, ok := configGetter(extraConfig).(influxConfig)
	if !ok {
		logger.Debug("no config fot the influx client. aborting")
		return errNoConfig
	}

	influxdbClient, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     cfg.address,
		Username: cfg.username,
		Password: cfg.password,
	})
	if err != nil {
		logger.Debug("influx client crashed")
		return err
	}

	t := time.NewTicker(cfg.ttl)

	cw := clientWrapper{
		influxClient: influxdbClient,
		collector:    metricsCollector,
		logger:       logger,
		db:           cfg.database,
	}

	go cw.keepUpdated(ctx, t.C)

	logger.Debug("client up and running")

	return nil
}

func (cw clientWrapper) keepUpdated(ctx context.Context, ticker <-chan time.Time) {
	for {
		select {
		case <-ticker:
		case <-ctx.Done():
			return
		}

		cw.logger.Info("Preparing points")

		snapshot := cw.collector.Snapshot()

		if shouldSendPoints := len(snapshot.Counters) > 0 || len(snapshot.Gauges) > 0; !shouldSendPoints {
			cw.logger.Info("no metrics to send to influx")
			continue
		}

		bp, _ := client.NewBatchPoints(client.BatchPointsConfig{
			Database:  cw.db,
			Precision: "s",
		})
		now := time.Unix(0, snapshot.Time)

		for _, p := range counter.Points(now, snapshot.Counters, cw.logger) {
			bp.AddPoint(p)
		}

		for _, p := range gauge.Points(now, snapshot.Gauges, cw.logger) {
			bp.AddPoint(p)
		}

		// TODO: collect all the other points

		if err := cw.influxClient.Write(bp); err != nil {
			cw.logger.Error("writting to influx:", err.Error())
		}

		cw.logger.Info(len(bp.Points()), "datapoints sent to Influx")
	}
}
