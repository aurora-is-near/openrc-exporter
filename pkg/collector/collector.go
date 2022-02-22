package collector

import (
	"log"
	"strconv"
	"time"

	"github.com/aurora-is-near/openrc-exporter/pkg/openrc"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "openrc"
)

type Collector struct {
	logger *log.Logger

	currentRunlevelMetric       *prometheus.GaugeVec
	serviceDaemonsCrashedMetric *prometheus.Desc
	serviceStatusMetric         *prometheus.Desc
	serviceStateMetric          *prometheus.Desc
	serviceUptimeSecondsMetric  *prometheus.Desc
	serviceStartCountMetric     *prometheus.Desc

	serviceRespawnDelaySecondsMetric  *prometheus.Desc
	serviceRespawnMaxMetric           *prometheus.Desc
	serviceRespawnPeriodSecondsMetric *prometheus.Desc
}

func New(logger *log.Logger) *Collector {
	return &Collector{
		logger: logger,

		currentRunlevelMetric: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "runlevel",
				Name:      "current",
				Help:      "The current runlevel",
			},
			[]string{"runlevel"},
		),
		serviceDaemonsCrashedMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "daemons_crashed"),
			"Whether the daemons started with start-stop-daemon are crashed",
			[]string{"service"}, nil,
		),
		serviceStatusMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "status"),
			"States that the service is in (1 means the service is in the state, 0 otherwise)",
			[]string{"service", "state"}, nil,
		),

		serviceStateMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "state"),
			"State that the service is in (1 means running, 0 means failed, -1 means not started.)",
			[]string{"service"}, nil,
		),

		serviceUptimeSecondsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "uptime_seconds"),
			"Number of seconds that the service is running (0 if not running)",
			[]string{"service"}, nil,
		),
		serviceStartCountMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "start_count"),
			"Number of times a service has been restarted",
			[]string{"service"}, nil,
		),
		serviceRespawnDelaySecondsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "respawn_delay_seconds"),
			"Number of seconds to wait before restarting a process that crashed (supervise-daemon only)",
			[]string{"service"}, nil,
		),
		serviceRespawnMaxMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "respawn_max"),
			"Maximum number of restarts within respawn_period before giving up (supervise-daemon only)",
			[]string{"service"}, nil,
		),
		serviceRespawnPeriodSecondsMetric: prometheus.NewDesc(
			prometheus.BuildFQName(namespace, "service", "respawn_period_seconds"),
			"The period within which respawn counts towards respawn_max (supervise-daemon only)",
			[]string{"service"}, nil,
		),
	}
}

// Describe implements prometheus.Collector's Describe
func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	c.currentRunlevelMetric.Describe(ch)

	ch <- c.serviceDaemonsCrashedMetric
	ch <- c.serviceStatusMetric
	ch <- c.serviceStateMetric
	ch <- c.serviceUptimeSecondsMetric
	ch <- c.serviceStartCountMetric
	ch <- c.serviceRespawnDelaySecondsMetric
	ch <- c.serviceRespawnMaxMetric
	ch <- c.serviceRespawnPeriodSecondsMetric
}

// Collect implements prometheus.Collector's Collect
func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	runlevel := openrc.RunlevelGet()
	currentRunlevelGauge := c.currentRunlevelMetric.WithLabelValues(runlevel)
	currentRunlevelGauge.Set(1)
	ch <- currentRunlevelGauge

	for _, service := range openrc.ServicesInRunlevel(nil) {
		uptimeSeconds := c.serviceUptimeSeconds(service)
		ch <- prometheus.MustNewConstMetric(
			c.serviceUptimeSecondsMetric,
			prometheus.GaugeValue,
			uptimeSeconds,
			service,
		)

		if startCount, err := serviceValueGetFloat64(service, "start_count"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.serviceStartCountMetric,
				prometheus.CounterValue,
				startCount,
				service,
			)
		} else {
			c.logger.Println("start_count:", err)
		}

		if respawnDelay, err := serviceValueGetFloat64(service, "respawn_delay"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.serviceRespawnDelaySecondsMetric,
				prometheus.GaugeValue,
				respawnDelay,
				service,
			)
		} else {
			c.logger.Println("respawn_delay:", err)
		}

		if respawnMax, err := serviceValueGetFloat64(service, "respawn_max"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.serviceRespawnMaxMetric,
				prometheus.GaugeValue,
				respawnMax,
				service,
			)
		} else {
			c.logger.Println("respawn_max:", err)
		}

		if respawnPeriod, err := serviceValueGetFloat64(service, "respawn_period"); err == nil {
			ch <- prometheus.MustNewConstMetric(
				c.serviceRespawnPeriodSecondsMetric,
				prometheus.GaugeValue,
				respawnPeriod,
				service,
			)
		} else {
			c.logger.Println("respawn_period:", err)
		}

		daemonsCrashed := openrc.ServiceDaemonsCrashed(service)
		ch <- prometheus.MustNewConstMetric(
			c.serviceDaemonsCrashedMetric,
			prometheus.GaugeValue,
			boolToValue(daemonsCrashed),
			service,
		)
		state := StateNum(service)
		if state >= 0 {
			ch <- prometheus.MustNewConstMetric(
				c.serviceStateMetric,
				prometheus.GaugeValue,
				float64(state),
				service,
			)
		}

		for state, stateName := range openrc.ServiceStateNames {
			isState := isServiceState(service, state)
			value := boolToValue(isState)

			ch <- prometheus.MustNewConstMetric(
				c.serviceStatusMetric,
				prometheus.GaugeValue,
				value,
				service, stateName,
			)
		}
	}
}

func (c *Collector) serviceUptimeSeconds(service string) float64 {
	startTimeStr := openrc.ServiceValueGet(service, "start_time")
	if startTimeStr != "" {
		uptimeSeconds, err := uptimeSecondsFromStartTime(startTimeStr)
		if err != nil {
			c.logger.Println(err)
			return float64(0)
		}

		return float64(uptimeSeconds)
	} else {
		return float64(0)
	}
}

func serviceValueGetFloat64(service string, option string) (float64, error) {
	valueStr := openrc.ServiceValueGet(service, option)
	if valueStr == "" {
		return float64(0), nil
	}
	return strconv.ParseFloat(valueStr, 64)
}

func uptimeSecondsFromStartTime(startTimeStr string) (float64, error) {
	layout := "2006-01-02 15:04:05"
	location, err := time.LoadLocation("Local")
	if err != nil {
		return 0, err
	}

	startTime, err := time.ParseInLocation(layout, startTimeStr, location)
	if err != nil {
		return 0, err
	}

	uptimeSeconds := float64(time.Since(startTime).Seconds())
	return uptimeSeconds, nil
}

func boolToValue(boolean bool) float64 {
	if boolean {
		return float64(1)
	} else {
		return float64(0)
	}
}

func isServiceState(service string, state int) bool {
	return (openrc.ServiceState(service) & state) > 0
}

func StateNum(service string) int {
	state := openrc.ServiceState(service)
	switch {
	case state&openrc.ServiceCrashed > 0:
		return 0
	case state&openrc.ServiceFailed > 0:
		return 0
	case state&openrc.ServiceStarted > 0:
		return 1
	default:
		return -1
	}
}
