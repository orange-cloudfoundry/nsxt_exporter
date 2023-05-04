package metrics

import (
	"fmt"
	"reflect"

	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/exp/constraints"
)

const (
	StatusStable    = "STABLE"
	StatusConnected = "CONNECTED"
	StatusUp        = "UP"
	StatusInSync    = "in_sync"
)

type floatable interface {
	constraints.Float | constraints.Integer
}

func zero[T any](value *T) T {
	var zero T
	if value != nil {
		return *value
	}
	return zero
}

func set[T floatable](metric prometheus.GaugeVec, labels []string, value T) {
	metric.WithLabelValues(labels...).Set(float64(value))
}

func setp[T floatable](metric prometheus.GaugeVec, labels []string, value *T) {
	if value != nil {
		metric.WithLabelValues(labels...).Set(float64(*value))
	}
}

func setb(metric prometheus.GaugeVec, labels []string, value *bool) {
	v := 0
	if value != nil && *value {
		v = 1
	}
	set(metric, labels, v)
}

// nolint: unparam
func setv(metric prometheus.GaugeVec, labels []string, value *string, expect string) {
	v := 0
	if value != nil && *value == expect {
		v = 1
	}
	set(metric, labels, v)
}

// set if labels - set value 1 with given labels if all labels and non-nil
func setl[T any](metric prometheus.GaugeVec, labels []string, values ...*T) {
	oLen := len(labels)
	for _, cVal := range values {
		if cVal != nil {
			var format string
			switch reflect.TypeOf(*cVal).Kind() {
			case reflect.Int64:
				format = "%d"
			default:
				format = "%s"
			}
			labels = slice(labels, fmt.Sprintf(format, *cVal))
		}
	}
	if len(labels) == oLen+len(values) {
		set(metric, labels, 1)
	}
}

func slice(s []string, vals ...string) []string {
	res := []string{}
	res = append(res, s...)
	res = append(res, vals...)
	return res
}

func statusToValue(value string, expect string) float64 {
	if value == expect {
		return 1.0
	}
	return 0.0
}
