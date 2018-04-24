package histogram

import (
	"testing"

	metrics "github.com/devopsfaith/krakend-metrics"
)

func Test_isEmpty(t *testing.T) {
	for _, h := range []metrics.HistogramData{
		{
			Max: 42,
		},
		{
			Min: 42,
		},
		{
			Mean: 42.0,
		},
		{
			Stddev: 42.0,
		},
		{
			Variance: 42.0,
		},
		{
			Percentiles: []float64{42.0, 0, 10},
		},
	} {
		if isEmpty(h) {
			t.Errorf("the histogram %v is not empty", h)
		}
	}

	if !isEmpty(metrics.HistogramData{}) {
		t.Errorf("unable to detect an empty histogram")
	}
}
