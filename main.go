package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin/v2"
	"github.com/orange-cloudfoundry/nsxt_exporter/api"
	"github.com/orange-cloudfoundry/nsxt_exporter/config"
	"github.com/orange-cloudfoundry/nsxt_exporter/metrics"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/sirupsen/logrus"
)

var (
	configFile = kingpin.Flag("config", "Configuration file path").Default("config.yml").File()
)

func main() {
	kingpin.Version(version.Print("nsxt-exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	object := config.NewConfig(*configFile)
	namespace := "nsxt"
	if object.Exporter.Namespace != "" {
		namespace = object.Exporter.Namespace
	}
	manager, err := api.NewNSXApi(object.Nsxt)
	if err != nil {
		logrus.Fatal(err)
	}

	lvl, err := logrus.ParseLevel(object.Log.Level)
	if err != nil {
		logrus.Warnf("invalid log.level value '%s'", object.Log.Level)
		lvl = logrus.InfoLevel
	}
	logrus.SetLevel(lvl)

	recorder := metrics.NewRecorder(manager, namespace)

	go func() {
		for {
			for {
				err := recorder.RecordMetrics()
				if err != nil {
					logrus.Error("Error when creating metrics: " + err.Error())
					break
				}
				time.Sleep(object.Exporter.IntervalDuration)
			}
			time.Sleep(object.Exporter.ErrorIntervalDuration)
		}
	}()
	http.Handle(object.Exporter.Path, promhttp.Handler())
	listen := ":" + strconv.Itoa(object.Exporter.Port)
	logrus.Infof("listening on %s", listen)

	// nolint: gosec
	err = http.ListenAndServe(listen, nil)
	if err != nil {
		logrus.Fatal("Error when serving: " + err.Error())
	}
}
